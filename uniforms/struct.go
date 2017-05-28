package uniforms

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// Struct is a wrapper around a program/uniform for binding. This struct only works for a pure float struct.
type Struct struct {
	program uint32
	uniform int32
}

// NewStruct instantiates a Struct for the provided program and uniform location.
func NewStruct(p uint32, u int32) *Struct {
	return &Struct{p, u}
}

// Set sets this Int to the provided data, and updates the uniform data.
func (s *Struct) Set(count int32, ptr *float32) {
	f := [7]float32{1, 1, 1, 1, 0, -1, 0}

	gl.ProgramUniform1fv(s.program, s.uniform, count, &f[0])
}
