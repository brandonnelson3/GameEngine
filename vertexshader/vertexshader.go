package vertexshader

import (
	"fmt"
	"strings"

	"github.com/brandonnelson3/GameEngine/uniforms"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	originalVertexSourceFile = `shader.vert`
	vertSrc                  = `
#version 450

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

in vec3 vert;
in vec3 norm;

out gl_PerVertex
{
    vec4 gl_Position;
	vec3 worldPosition;
	vec3 normal;
} vertex_out;

void main() {
    gl_Position = projection * view * model * vec4(vert, 1);
	vertex_out.worldPosition = vec3(model * vec4(vert, 1));
	vertex_out.normal = vec3(vec4(norm, 1));
}` + "\x00"
)

// Vertex is a Vertex.
type Vertex struct {
	Vert, Norm mgl32.Vec3
}

// VertexShader is a VertexShader.
type VertexShader struct {
	uint32

	Projection, View, Model, Rotation *uniforms.Matrix4
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
	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))

	gl.DeleteShader(shader)

	return &VertexShader{
		uint32:     program,
		Projection: uniforms.NewMatrix4(program, projectionLoc),
		View:       uniforms.NewMatrix4(program, viewLoc),
		Model:      uniforms.NewMatrix4(program, modelLoc),
	}, nil
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *VertexShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.VERTEX_SHADER_BIT, s.uint32)
}

// BindVertexAttributes binds the attributes per vertex.
func (s *VertexShader) BindVertexAttributes() {
	vertAttrib := uint32(gl.GetAttribLocation(s.uint32, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	normAttrib := uint32(gl.GetAttribLocation(s.uint32, gl.Str("norm\x00")))
	gl.EnableVertexAttribArray(normAttrib)
	gl.VertexAttribPointer(normAttrib, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(12))
}
