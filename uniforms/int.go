package uniforms

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// Int is a wrapper around a int32, and a program/uniform for binding.
type Int struct {
	int32

	program uint32
	uniform int32
}

// NewInt instantiates a default int for the provided program and uniform location.
func NewInt(p uint32, u int32) *Int {
	return &Int{0, p, u}
}

// Set Sets this Vector4 to the provided data, and updates the uniform data.
func (m *Int) Set(i int32) {
	m.int32 = i
	gl.ProgramUniform1i(m.program, m.uniform, m.int32)
}
