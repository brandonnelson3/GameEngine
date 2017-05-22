package fragmentshader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

const (
	originalFragmentSourceFile = `shader.frag`
	fragSrc                    = `
#version 450

// TODO: Probably can pull this out into a common place.
struct PointLight {
	vec4 color;
	vec4 position;
	vec4 paddingAndRadius;
};

struct VisibleIndex {
	int index;
};

// Shader storage buffer objects
layout(std430, binding = 0) readonly buffer LightBuffer {
	PointLight data[];
} lightBuffer;

layout(std430, binding = 1) readonly buffer VisibleLightIndicesBuffer {
	VisibleIndex data[];
} visibleLightIndicesBuffer;

out vec4 outputColor;


void main() {
	ivec2 location = ivec2(gl_FragCoord.xy);
	// TODO: Put this 16 somewhere constant.
	ivec2 tileID = location / ivec2(16, 16);
	// TODO: Pass in numberOfTilesX as a uniform.
	uint numberOfTilesX = 50;
	uint index = tileID.y * numberOfTilesX + tileID.x;

	// TODO 1024 should be somewhere constant.
	uint offset = index * 1024;
	uint i = 0;
	for (i; i < 1024 && visibleLightIndicesBuffer.data[offset + i].index != -1; i++) {
		//uint lightIndex = visibleLightIndicesBuffer.data[offset + i].index;
		//PointLight light = lightBuffer.data[0];
	}
	if (i == 1) {
		outputColor = vec4(1.0,0,0,1.0);
	} else {
		outputColor = vec4(0, 1.0, 0, 1.0);
	}

	//outputColor = outputColor + vec4(0.2, 0.2, 0.2, 1.0);
}` + "\x00"
)

// FragmentShader represents a FragmentShader
type FragmentShader struct {
	uint32
}

// NewFragmentShader instantiates and initializes a FragmentShader object.
func NewFragmentShader() (*FragmentShader, error) {
	program := gl.CreateProgram()
	shader := gl.CreateShader(gl.FRAGMENT_SHADER)

	csources, free := gl.Strs(fragSrc)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to compile %v: %v", originalFragmentSourceFile, log)
	}

	gl.AttachShader(program, shader)
	gl.ProgramParameteri(program, gl.PROGRAM_SEPARABLE, gl.TRUE)
	gl.LinkProgram(program)

	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("failed to link %v: %v", originalFragmentSourceFile, log)
	}

	gl.DeleteShader(shader)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	return &FragmentShader{
		uint32: program,
	}, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *FragmentShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.FRAGMENT_SHADER_BIT, s.uint32)
}
