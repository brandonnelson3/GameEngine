package uniforms

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

// Sampler2D is a wrapper around a int32 which is the sampler texture id, and a program/uniform for binding.
type Sampler2D struct {
	uint32

	program uint32
	uniform int32
}

// NewSampler2D instantiates a sampler2d for the provided program, and uniform location.
func NewSampler2D(p uint32, u int32) *Sampler2D {
	return &Sampler2D{0, p, u}
}

// Set Sets this Sampler2D to the provided id, and updates the uniform data.
func (m *Sampler2D) Set(samplerID uint32) {
	m.uint32 = samplerID
	gl.ActiveTexture(gl.TEXTURE4)
	gl.Uniform1i(m.uniform, 4)
	gl.BindTexture(gl.TEXTURE_2D, m.uint32)
}
