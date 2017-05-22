package depthfragmentshader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
)

const (
	originalFragmentSourceFile = `shader.frag`
	fragSrc                    = `
#version 450
void main() {
	// We are not drawing anything to the screen, so nothing to be done here
}` + "\x00"
)

// DepthFragmentShader represents a FragmentShader
type DepthFragmentShader struct {
	uint32
}

// NewDepthFragmentShader instantiates and initializes a DepthFragmentShader object.
func NewDepthFragmentShader() (*DepthFragmentShader, error) {
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

	return &DepthFragmentShader{
		uint32: program,
	}, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *DepthFragmentShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.FRAGMENT_SHADER_BIT, s.uint32)
}
