package fragmentshader

import (
	"fmt"
	"strings"

	"github.com/brandonnelson3/GameEngine/buffers"
	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/brandonnelson3/GameEngine/uniforms"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	originalFragmentSourceFile = `shader.frag`
	fragSrc                    = `
#version 450

// TODO: Probably can pull this out into a common place.
struct PointLight {
	vec3 color;
	float intensity;
	vec3 position;
	float radius;
};

struct VisibleIndex {
	int index;
};

struct DirectionalLight {
	vec3 color;
	float brightness;
	vec3 direction;
};

// Shader storage buffer objects
layout(std430, binding = 0) readonly buffer LightBuffer {
	PointLight data[];
} lightBuffer;

layout(std430, binding = 1) readonly buffer VisibleLightIndicesBuffer {
	VisibleIndex data[];
} visibleLightIndicesBuffer;

layout(std430, binding = 2) readonly buffer DirectionalLightBuffer {
	DirectionalLight data;
} directionalLightBuffer;

uniform int renderMode;
uniform uint numTilesX;
uniform sampler2D diffuse;

in VERTEX_OUT
{
    vec4 gl_Position;
	vec3 worldPosition;
	vec3 normal;
	vec2 uv;
} fragment_in;

out vec4 outputColor;

void main() {
	ivec2 location = ivec2(gl_FragCoord.xy);
	// TODO: Put this 16 somewhere constant.
	ivec2 tileID = location / ivec2(16, 16);
	uint index = tileID.y * numTilesX + tileID.x;

	// TODO 1024 should be somewhere constant.
	uint offset = index * 1024;
	
	if (renderMode == 0) {
		vec3 pointLightColor = vec3(0, 0, 0);

		uint i=0;
		for (i; i < 1024 && visibleLightIndicesBuffer.data[offset + i].index != -1; i++) {
			uint lightIndex = visibleLightIndicesBuffer.data[offset + i].index;
			PointLight light = lightBuffer.data[lightIndex];
			vec3 lightVector = light.position - fragment_in.worldPosition;
			float dist = length(lightVector);
			float NdL = max(0.0f, dot(fragment_in.normal, lightVector*(1.0f/dist)));
			float attenuation = 1.0f - clamp(dist * (1.0/(light.radius)), 0.0, 1.0);
			vec3 diffuse = NdL * light.color * light.intensity;
			pointLightColor += attenuation * diffuse;
		}

		DirectionalLight directionalLight = directionalLightBuffer.data;
		float NdL = max(0.0f, dot(fragment_in.normal, -1*directionalLight.direction));
		vec3 directionalLightColor = NdL * directionalLight.color * directionalLight.brightness;

		outputColor = texture(diffuse, fragment_in.uv) * vec4(pointLightColor+directionalLightColor, 1.0);
	} else if (renderMode == 1) {
		uint i=0;
		for (i; i < 1024 && visibleLightIndicesBuffer.data[offset + i].index != -1; i++) {}
		outputColor = vec4(vec3(float(i)/256)+vec3(0.1), 1.0);
	} else if (renderMode == 2) {
		outputColor = vec4(abs(fragment_in.normal), 1.0);
	} else if (renderMode == 3) {
		outputColor = vec4(fragment_in.uv, 0, 1.0);
	} else if (renderMode == 4) {
		outputColor = texture(diffuse, fragment_in.uv);
	}	
}
` + "\x00"
)

// FragmentShader represents a FragmentShader
type FragmentShader struct {
	uint32

	RenderMode *uniforms.Int
	NumTilesX  *uniforms.UInt
	Diffuse    *uniforms.Sampler2D

	LightBuffer, VisibleLightIndicesBuffer, DirectionalLightBuffer *buffers.Binding
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

	renderModeLoc := gl.GetUniformLocation(program, gl.Str("renderMode\x00"))
	numTilesXLoc := gl.GetUniformLocation(program, gl.Str("numTilesX\x00"))
	diffuseLoc := gl.GetUniformLocation(program, gl.Str("diffuse\x00"))

	gl.DeleteShader(shader)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	fs := &FragmentShader{
		uint32:                    program,
		RenderMode:                uniforms.NewInt(program, renderModeLoc),
		NumTilesX:                 uniforms.NewUInt(program, numTilesXLoc),
		Diffuse:                   uniforms.NewSampler2D(program, diffuseLoc),
		LightBuffer:               buffers.NewBinding(0),
		VisibleLightIndicesBuffer: buffers.NewBinding(1),
		DirectionalLightBuffer:    buffers.NewBinding(2),
	}

	messagebus.RegisterType("key", func(m *messagebus.Message) {
		pressedKeys := m.Data1.([]glfw.Key)
		for _, key := range pressedKeys {
			if key >= glfw.KeyF1 && key <= glfw.KeyF25 {
				fs.RenderMode.Set(int32(key - glfw.KeyF1))
			}
		}
	})

	return fs, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *FragmentShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.FRAGMENT_SHADER_BIT, s.uint32)
}
