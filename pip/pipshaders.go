package pip

import (
	"fmt"
	"strings"

	"github.com/brandonnelson3/GameEngine/uniforms"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	originalSourceFile = `pipshader.`
	fragSrc            = `
#version 450

uniform sampler2D textureSampler;
uniform mat4 projection;

in VERTEX_OUT
{
    vec4 gl_Position;
	vec2 uv;
} fragment_in;

out vec4 outputColor;

void main() {
	float depth = texture(textureSampler, fragment_in.uv).r;
	// Linearize the depth value from depth buffer (must do this because we created it using projection)
	depth = 1 - 1/log((0.5 * projection[3][2]) / (depth + 0.5 * projection[2][2] - 0.5));

	outputColor = vec4(vec3(depth), 1.0);
}` + "\x00"
	vertSrc = `
#version 450

in vec2 pos;
in vec2 uv;

uniform mat4 projection;

out gl_PerVertex
{
    vec4 gl_Position;
	vec2 uv;
} vertex_out;

void main() {
    gl_Position = projection * vec4(pos, 0, 1);
	vertex_out.uv = uv;
}` + "\x00"
)

// FragmentShader represents a FragmentShader
type FragmentShader struct {
	uint32

	DepthMap   *uniforms.Sampler2D
	Projection *uniforms.Matrix4
}

// Vertex is a Vertex.
type Vertex struct {
	Pos, UV mgl32.Vec2
}

// VertexShader is a VertexShader.
type VertexShader struct {
	uint32

	Projection *uniforms.Matrix4
}

// NewFragmentShader instantiates and initializes a PipFragmentShader object.
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

		return nil, fmt.Errorf("failed to compile %v: %v", originalSourceFile+"frag", log)
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

		return nil, fmt.Errorf("failed to link %v: %v", originalSourceFile+"frag", log)
	}

	depthMapLoc := gl.GetUniformLocation(program, gl.Str("textureSampler\x00"))
	projectionLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	gl.DeleteShader(shader)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	return &FragmentShader{
		uint32:     program,
		DepthMap:   uniforms.NewSampler2D(program, depthMapLoc),
		Projection: uniforms.NewMatrix4(program, projectionLoc),
	}, nil
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

		return nil, fmt.Errorf("failed to compile %v: %v", originalSourceFile+"vert", log)
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

		return nil, fmt.Errorf("failed to link %v: %v", originalSourceFile+"vert", log)
	}

	projectionLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	gl.DeleteShader(shader)

	return &VertexShader{
		uint32:     program,
		Projection: uniforms.NewMatrix4(program, projectionLoc),
	}, nil
}

// BindVertexAttributes binds the attributes per vertex.
func (s *VertexShader) BindVertexAttributes() {
	posAttrib := uint32(gl.GetAttribLocation(s.uint32, gl.Str("pos\x00")))
	gl.EnableVertexAttribArray(posAttrib)
	gl.VertexAttribPointer(posAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	uvAttrib := uint32(gl.GetAttribLocation(s.uint32, gl.Str("uv\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointer(uvAttrib, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(8))
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *FragmentShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.FRAGMENT_SHADER_BIT, s.uint32)
}

// AddToPipeline adds this shader to the provided pipeline.
func (s *VertexShader) AddToPipeline(pipeline uint32) {
	gl.UseProgramStages(pipeline, gl.VERTEX_SHADER_BIT, s.uint32)
}
