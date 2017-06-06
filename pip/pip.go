package pip

import (
	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/brandonnelson3/GameEngine/window"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	pipeline, planeVao uint32

	vertexShader   *VertexShader
	fragmentShader *FragmentShader

	Enabled  = true
	DepthMap *uint32
)

func Initialize(m *uint32) {
	DepthMap = m

	var err error
	vertexShader, err = NewVertexShader()
	if err != nil {
		panic(err)
	}
	fragmentShader, err = NewFragmentShader()
	if err != nil {
		panic(err)
	}

	gl.CreateProgramPipelines(1, &pipeline)
	vertexShader.AddToPipeline(pipeline)
	fragmentShader.AddToPipeline(pipeline)
	gl.ValidateProgramPipeline(pipeline)
	gl.UseProgram(0)
	gl.BindProgramPipeline(pipeline)

	sizex := float32(480.0)
	sizey := float32(360.0)

	padding := uint32(50)

	topLeft := mgl32.Vec2{float32(window.Width-padding) - sizex, float32(window.Height-padding) - sizey}
	topRight := mgl32.Vec2{float32(window.Width - padding), float32(window.Height-padding) - sizey}
	botLeft := mgl32.Vec2{float32(window.Width-padding) - sizex, float32(window.Height - padding)}
	botRight := mgl32.Vec2{float32(window.Width - padding), float32(window.Height - padding)}

	planeVertices := []Vertex{
		{topLeft, mgl32.Vec2{0, 1}},
		{topRight, mgl32.Vec2{1, 1}},
		{botRight, mgl32.Vec2{1, 0}},
		{topLeft, mgl32.Vec2{0, 1}},
		{botRight, mgl32.Vec2{1, 0}},
		{botLeft, mgl32.Vec2{0, 0}},
	}

	gl.GenVertexArrays(1, &planeVao)
	gl.BindVertexArray(planeVao)

	var planeVbo uint32
	gl.GenBuffers(1, &planeVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4*4, gl.Ptr(planeVertices), gl.STATIC_DRAW)

	vertexShader.BindVertexAttributes()

	messagebus.RegisterType("key", func(m *messagebus.Message) {
		pressedKeys := m.Data1.([]glfw.Key)
		for _, key := range pressedKeys {
			switch key {
			case glfw.KeyPageUp:
				Enabled = true
			case glfw.KeyPageDown:
				Enabled = false
			}
		}
	})
}

func Render(p mgl32.Mat4) {
	gl.BindProgramPipeline(pipeline)
	vertexShader.Projection.Set(mgl32.Ortho(0.0, float32(window.Width), float32(window.Height), 0.0, -1.0, 1.0))
	// This is intentionally different since it needs to be the projection matrix that the depthMap was rendered with.
	fragmentShader.Projection.Set(p)
	fragmentShader.DepthMap.Set(gl.TEXTURE4, 4, *DepthMap)
	gl.BindVertexArray(planeVao)
	gl.DrawArrays(gl.TRIANGLES, 0, 2*3)
	gl.ActiveTexture(gl.TEXTURE4)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}
