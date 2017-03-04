package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

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

	vertexShader, err := NewVertexShader()
	if err != nil {
		panic(err)
	}

	fragmentShader, err := NewFragmentShader()
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

	vertexShader.Projection.Set(mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100.0))
	vertexShader.Camera.Set(mgl32.LookAtV(mgl32.Vec3{-25, 13, -25}, mgl32.Vec3{15, 0, 15}, mgl32.Vec3{0, 1, 0}))
	vertexShader.Model.Set(mgl32.Ident4())

	fragmentShader.Color.Set(mgl32.Vec4{0, 1, 0, 1})

	fragmentShader.BindFragmentOutputDataLocation()

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
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		modelRotation := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		gl.BindVertexArray(vao)

		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				modelTranslation := mgl32.Translate3D(float32(5*x), 0.0, float32(5*y))
				vertexShader.Model.Set(modelTranslation.Mul4(modelRotation))
				fragmentShader.Color.Set(mgl32.Vec4{float32(x) / 10, float32(y) / 10, float32(x*y) / 100, 1})
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

var cubeVertices = []Vertex{
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
