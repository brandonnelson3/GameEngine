package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"unsafe"

	"github.com/brandonnelson3/GameEngine/depthfragmentshader"
	"github.com/brandonnelson3/GameEngine/depthvertexshader"
	"github.com/brandonnelson3/GameEngine/fragmentshader"
	"github.com/brandonnelson3/GameEngine/framerate"
	"github.com/brandonnelson3/GameEngine/input"
	"github.com/brandonnelson3/GameEngine/lightcullingshader"
	"github.com/brandonnelson3/GameEngine/lights"
	"github.com/brandonnelson3/GameEngine/timer"
	"github.com/brandonnelson3/GameEngine/uniforms"
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

	// Initially place cursor dead center.
	window.RecenterCursor()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	var lightBuffer uint32
	var visibleLightIndicesBuffer uint32

	// Prepare light buffers
	gl.GenBuffers(1, &lightBuffer)
	gl.GenBuffers(1, &visibleLightIndicesBuffer)

	// Bind light buffer
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, lightBuffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, lights.MaximumPointLights*int(unsafe.Sizeof(&lights.PointLight{})), unsafe.Pointer(&lights.PointLights), gl.DYNAMIC_DRAW)

	// Bind visible light indices buffer
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, visibleLightIndicesBuffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, window.GetTotalNumTiles()*int(unsafe.Sizeof(&lights.VisibleIndex{}))*lights.MaximumPointLights, nil, gl.STATIC_DRAW)

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	// Build Depth Pipeline
	depthVertexShader, err := depthvertexshader.NewDepthVertexShader()
	if err != nil {
		panic(err)
	}
	depthFragmentShader, err := depthfragmentshader.NewDepthFragmentShader()
	if err != nil {
		panic(err)
	}
	var depthPipeline uint32
	gl.CreateProgramPipelines(1, &depthPipeline)
	depthVertexShader.AddToPipeline(depthPipeline)
	depthFragmentShader.AddToPipeline(depthPipeline)
	gl.ValidateProgramPipeline(depthPipeline)
	depthVertexShader.BindVertexAttributes()

	lightCullingShader, err := lightcullingshader.NewLightCullingShader()
	if err != nil {
		panic(err)
	}

	// Build Normal Pipeline
	vertexShader, err := vertexshader.NewVertexShader()
	if err != nil {
		panic(err)
	}
	fragmentShader, err := fragmentshader.NewFragmentShader()
	if err != nil {
		panic(err)
	}
	var normalPipeline uint32
	gl.CreateProgramPipelines(1, &normalPipeline)
	vertexShader.AddToPipeline(normalPipeline)
	fragmentShader.AddToPipeline(normalPipeline)
	gl.ValidateProgramPipeline(normalPipeline)
	gl.UseProgram(0)
	gl.BindProgramPipeline(normalPipeline)

	// Build Depth FrameBuffer
	var depthMapFBO uint32
	gl.GenFramebuffers(1, &depthMapFBO)
	var depthMap uint32
	gl.GenTextures(1, &depthMap)
	gl.BindTexture(gl.TEXTURE_2D, depthMap)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.DEPTH_COMPONENT, int32(window.Width), int32(window.Height), 0, gl.DEPTH_COMPONENT, gl.FLOAT, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_BORDER)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_BORDER)
	borderColor := []float32{1.0, 1.0, 1.0, 1.0}
	gl.TexParameterfv(gl.TEXTURE_2D, gl.TEXTURE_BORDER_COLOR, &borderColor[0])
	gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.TEXTURE_2D, depthMap, 0)
	gl.DrawBuffer(gl.NONE)
	gl.ReadBuffer(gl.NONE)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*6*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertexShader.BindVertexAttributes()

	angle := 0.0
	camera := NewFirstPersonCamera()
	go camera.Log()

	for !w.ShouldClose() {
		timer.BeginningOfFrame()
		framerate.BeginningOfFrame(timer.GetTime())
		input.Update()
		camera.Update(timer.GetPreviousFrameLength())

		elapsed := timer.GetPreviousFrameLength()

		angle += elapsed
		modelRotation := mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		gl.BindVertexArray(vao)

		// Step 1: Depth Pass
		gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
		gl.BindProgramPipeline(depthPipeline)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		depthVertexShader.View.Set(camera.GetView())
		depthVertexShader.Projection.Set(window.GetProjection())
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				modelTranslation := mgl32.Translate3D(float32(4*x), 0.0, float32(4*y))
				depthVertexShader.Model.Set(modelTranslation.Mul4(modelRotation))
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}

		// Step 2: Light Culling
		lightCullingShader.Use()
		lightCullingShader.View.Set(camera.GetView())
		lightCullingShader.Projection.Set(window.GetProjection())
		lightCullingShader.DepthMap.Set(depthMap)
		lightCullingShader.ScreenSize.Set(uniforms.IVec2{window.Width, window.Height})
		lightCullingShader.LightCount.Set(2)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, lightBuffer)
		gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, visibleLightIndicesBuffer)
		gl.DispatchCompute(50, 40, 1)
		// TODO: Fix this...
		gl.ActiveTexture(gl.TEXTURE4)
		gl.BindTexture(gl.TEXTURE_2D, 0)
		gl.UseProgram(0)

		// Step 3: Normal pass
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.BindProgramPipeline(normalPipeline)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		vertexShader.View.Set(camera.GetView())
		vertexShader.Projection.Set(window.GetProjection())
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				modelTranslation := mgl32.Translate3D(float32(4*x), 0.0, float32(4*y))
				vertexShader.Rotation.Set(modelRotation)
				vertexShader.Model.Set(modelTranslation.Mul4(modelRotation))
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
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec3{0, -1.0, 0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, -1.0, 0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{0, -1.0, 0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, -1.0, 0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{0, -1.0, 0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{0, -1.0, 0}},

	// Top
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{0, 1.0, 0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 1.0, 0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{0, 1.0, 0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{0, 1.0, 0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 1.0, 0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{0, 1.0, 0}},

	// Front
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}},

	// Back
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}},

	// Left
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{-1.0, 0, 0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{-1.0, 0, 0}},
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec3{-1.0, 0, 0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{-1.0, 0, 0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{-1.0, 0, 0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{-1.0, 0, 0}},

	// Right
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{1.0, 0, 0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{1.0, 0, 0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{1.0, 0, 0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{1.0, 0, 0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{1.0, 0, 0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 0, 0}},
}
