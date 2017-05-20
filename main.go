package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brandonnelson3/GameEngine/fragmentshader"
	"github.com/brandonnelson3/GameEngine/framerate"
	"github.com/brandonnelson3/GameEngine/input"
	"github.com/brandonnelson3/GameEngine/timer"
	"github.com/brandonnelson3/GameEngine/vertexshader"
	"github.com/brandonnelson3/GameEngine/window"
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	w := window.Create()

	w.SetKeyCallback(input.KeyCallBack)
	w.SetMouseButtonCallback(input.MouseButtonCallback)
	w.SetCursorPosCallback(input.CursorPosCallback)
	w.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	vertexShader, err := vertexshader.NewVertexShader()
	if err != nil {
		panic(err)
	}

	fragmentShader, err := fragmentshader.NewFragmentShader()
	if err != nil {
		panic(err)
	}

	var pipeline uint32
	gl.CreateProgramPipelines(1, &pipeline)

	vertexShader.AddToPipeline(pipeline)
	fragmentShader.AddToPipeline(pipeline)

	gl.ValidateProgramPipeline(pipeline)

	gl.UseProgram(0)
	gl.BindProgramPipeline(pipeline)

	vertexShader.Projection.Set(window.GetProjection())

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*3*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertexShader.BindVertexAttributes()

	angle := 0.0

	camera := NewFirstPersonCamera()

	for !w.ShouldClose() {
		timer.BeginningOfFrame()
		framerate.BeginningOfFrame(timer.GetTime())
		input.Update()
		camera.Update(timer.GetPreviousFrameLength())

		vertexShader.Camera.Set(camera.GetView())

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		elapsed := timer.GetPreviousFrameLength()

		angle += elapsed
		modelRotation := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		gl.BindVertexArray(vao)

		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				modelTranslation := mgl32.Translate3D(float32(4*x), 0.0, float32(4*y))
				vertexShader.Model.Set(modelTranslation.Mul4(modelRotation))
				fragmentShader.Color.Set(mgl32.Vec4{float32(x) / 10, float32(y) / 10, float32(x*y) / 100, 1})
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}

		// Maintenance
		w.SwapBuffers()
		glfw.PollEvents()

		window.RecenterCursor()
		framerate.EndOfFrame(timer.GetTime())
	}
}

var cubeVertices = []vertexshader.Vertex{
	//  X, Y, Z
	// Bottom
	{mgl32.Vec3{-1.0, -1.0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}},

	// Top
	{mgl32.Vec3{-1.0, 1.0, -1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}},

	// Front
	{mgl32.Vec3{-1.0, -1.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}},

	// Back
	{mgl32.Vec3{-1.0, -1.0, -1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}},

	// Left
	{mgl32.Vec3{-1.0, -1.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}},
	{mgl32.Vec3{-1.0, -1.0, -1.0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}},

	// Right
	{mgl32.Vec3{1.0, -1.0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}},
}
