package main

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Matrix4 is a wrapper around a mgl32.Mat4, and a program/uniform for binding.
type Matrix4 struct {
	mgl32.Mat4

	program uint32
	uniform int32
}

// NewMatrix4 instantiates an identity matrix for the provided program and uniform location.
func NewMatrix4(p uint32, u int32) *Matrix4 {
	return &Matrix4{mgl32.Ident4(), p, u}
}

// Set Sets this Matrix4 to the provided data, and updates the uniform data.
func (m *Matrix4) Set(nm mgl32.Mat4) {
	m.Mat4 = nm
	gl.ProgramUniformMatrix4fv(m.program, m.uniform, 1, false, &m.Mat4[0])
}
