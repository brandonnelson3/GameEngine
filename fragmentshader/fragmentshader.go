package fragmentshader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"

	"github.com/brandonnelson3/GameEngine/math"
)

const (
	originalFragmentSourceFile = `shader.frag`
	fragSrc                    = `
#version 450
uniform vec4 color;
out vec4 outputColor;
void main() {
    outputColor = color;
}` + "\x00"
)

// FragmentShader represents a FragmentShader
type FragmentShader struct {
	uint32

	Color *math.Vector4
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

	colorLoc := gl.GetUniformLocation(program, gl.Str("color\x00"))

	gl.DeleteShader(shader)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	return &FragmentShader{
		uint32: program,
		Color:  math.NewVector4(program, colorLoc),
	}, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *FragmentShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.FRAGMENT_SHADER_BIT, s.uint32)
}
