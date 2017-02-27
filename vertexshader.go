package main

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	originalVertexSourceFile = `shader.vert`
	vertSrc                  = `
#version 450
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
out gl_PerVertex
{
    vec4 gl_Position;
};
void main() {
    gl_Position = projection * camera * model * vec4(vert, 1);
}` + "\x00"
)

// Vertex is a Vertex.
type Vertex struct {
	vert mgl32.Vec3
}

// VertexShader is a VertexShader.
type VertexShader struct {
	program uint32

	// Uniforms.
	projectionLoc, cameraLoc, modelLoc int32
	projectionPtr, cameraPtr, modelPtr *mgl32.Mat4
}

// NewVertexShader instantiates and initializes a shader object.
func NewVertexShader() (*VertexShader, error) {
	program := gl.CreateProgram()
	shader := gl.CreateShader(gl.VERTEX_SHADER)

	csources, free := gl.Strs(vertSrc)
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

		return nil, fmt.Errorf("failed to compile %v: %v", originalVertexSourceFile, log)
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

		return nil, fmt.Errorf("failed to link %v: %v", originalVertexSourceFile, log)
	}

	projectionLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraLoc := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.DeleteShader(shader)

	return &VertexShader{
		program:       program,
		projectionLoc: projectionLoc,
		cameraLoc:     cameraLoc,
		modelLoc:      modelLoc,
	}, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *VertexShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.VERTEX_SHADER_BIT, s.program)
}

// SetProjectionMatrix sets the pointer to the ProjectionMatrix desired.
func (s *VertexShader) SetProjectionMatrix(m *mgl32.Mat4) {
	s.projectionPtr = m
}

// SetCameraMatrix sets the pointer to the CameraMatrix desired.
func (s *VertexShader) SetCameraMatrix(m *mgl32.Mat4) {
	s.cameraPtr = m
}

// SetModelMatrix sets the pointer to the ModelMatrix desired.
func (s *VertexShader) SetModelMatrix(m *mgl32.Mat4) {
	s.modelPtr = m
}

// UpdateProjectionMatrix updates the copy of the ProjectionMatrix on the GPU.
func (s *VertexShader) UpdateProjectionMatrix() {
	gl.ProgramUniformMatrix4fv(s.program, s.projectionLoc, 1, false, &s.projectionPtr[0])
}

// UpdateCameraMatrix updates the copy of the CameraMatrix on the GPU.
func (s *VertexShader) UpdateCameraMatrix() {
	gl.ProgramUniformMatrix4fv(s.program, s.cameraLoc, 1, false, &s.cameraPtr[0])
}

// UpdateModelMatrix updates the copy of the ModelMatrix on the GPU.
func (s *VertexShader) UpdateModelMatrix() {
	gl.ProgramUniformMatrix4fv(s.program, s.modelLoc, 1, false, &s.modelPtr[0])
}

// UpdateAll updates the copy of all uniforms on GPU.
func (s *VertexShader) UpdateAll() {
	gl.ProgramUniformMatrix4fv(s.program, s.projectionLoc, 1, false, &s.projectionPtr[0])
	gl.ProgramUniformMatrix4fv(s.program, s.cameraLoc, 1, false, &s.cameraPtr[0])
	gl.ProgramUniformMatrix4fv(s.program, s.modelLoc, 1, false, &s.modelPtr[0])
}

// BindVertexAttributes binds the attributes per vertex.
func (s *VertexShader) BindVertexAttributes() {
	vertAttrib := uint32(gl.GetAttribLocation(s.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
}
