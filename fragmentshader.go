package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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
	program, shader uint32

	// Uniforms.
	colorLoc int32
	colorPtr *mgl32.Vec4
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

	return &FragmentShader{
		program:  program,
		shader:   shader,
		colorLoc: colorLoc,
	}, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *FragmentShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.FRAGMENT_SHADER_BIT, s.program)
}

// SetColorVector sets the pointer to the ColorVector desired.
func (s *FragmentShader) SetColorVector(v *mgl32.Vec4) {
	s.colorPtr = v
}

// UpdateColorVector updates the copy of the ColorVector on the GPU.
func (s *FragmentShader) UpdateColorVector() {
	gl.ProgramUniform4fv(s.program, s.colorLoc, 1, &s.colorPtr[0])
}

// UpdateAll updates the copy of all uniforms on GPU.
func (s *FragmentShader) UpdateAll() {
	gl.ProgramUniform4fv(s.program, s.colorLoc, 1, &s.colorPtr[0])
}

// BindFragmentOutputDataLocation binds attribute which contains the output color.
func (s *FragmentShader) BindFragmentOutputDataLocation() {
	gl.BindFragDataLocation(s.program, 0, gl.Str("outputColor\x00"))
}
