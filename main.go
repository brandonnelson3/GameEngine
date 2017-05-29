package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/brandonnelson3/GameEngine/depthfragmentshader"
	"github.com/brandonnelson3/GameEngine/depthvertexshader"
	"github.com/brandonnelson3/GameEngine/fragmentshader"
	"github.com/brandonnelson3/GameEngine/framerate"
	"github.com/brandonnelson3/GameEngine/input"
	"github.com/brandonnelson3/GameEngine/lightcullingshader"
	"github.com/brandonnelson3/GameEngine/lights"
	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/brandonnelson3/GameEngine/pip"
	"github.com/brandonnelson3/GameEngine/textures"
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
	gl.Enable(gl.MULTISAMPLE)
	gl.DepthMask(true)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	pip.Initialize()

	lights.InitPointLights()
	lights.InitDirectionalLights()

	diffuseTexture, err := textures.NewFromPng("crate1_diffuse.png")
	if err != nil {
		panic(err)
	}

	sandTexture, err := textures.NewFromPng("sand.png")
	if err != nil {
		panic(err)
	}

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

	pip.DepthMap = &depthMap

	// Configure the vertex data
	var cubeVao uint32
	gl.GenVertexArrays(1, &cubeVao)
	gl.BindVertexArray(cubeVao)

	var cubeVbo uint32
	gl.GenBuffers(1, &cubeVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*8*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertexShader.BindVertexAttributes()

	var planeVao uint32
	gl.GenVertexArrays(1, &planeVao)
	gl.BindVertexArray(planeVao)

	var planeVbo uint32
	gl.GenBuffers(1, &planeVbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, planeVbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*8*4, gl.Ptr(planeVertices), gl.STATIC_DRAW)

	vertexShader.BindVertexAttributes()

	camera := NewFirstPersonCamera()

	messagebus.RegisterType("key", func(m *messagebus.Message) {
		pressedKeys := m.Data2.([]glfw.Key)
		for _, key := range pressedKeys {
			switch key {
			case glfw.KeyL:
				lights.AddPointLight(camera.GetPosition(), mgl32.Vec3{1, 1, 1}, 1, 10)
			}
		}
	})

	for !w.ShouldClose() {
		timer.BeginningOfFrame()
		framerate.BeginningOfFrame(timer.GetTime())
		input.Update()
		camera.Update(timer.GetPreviousFrameLength())

		// Step 1: Depth Pass
		gl.BindFramebuffer(gl.FRAMEBUFFER, depthMapFBO)
		gl.BindProgramPipeline(depthPipeline)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
		depthVertexShader.View.Set(camera.GetView())
		depthVertexShader.Projection.Set(window.GetProjection())
		gl.BindVertexArray(cubeVao)
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				modelTranslation := mgl32.Translate3D(float32(4*x), 5.0, float32(4*y))
				depthVertexShader.Model.Set(modelTranslation)
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}
		depthVertexShader.Model.Set(mgl32.Ident4())
		gl.BindVertexArray(planeVao)
		gl.DrawArrays(gl.TRIANGLES, 0, 2*3)

		// Step 2: Light Culling
		lightCullingShader.Use()
		lightCullingShader.View.Set(camera.GetView())
		lightCullingShader.Projection.Set(window.GetProjection())
		lightCullingShader.DepthMap.Set(gl.TEXTURE4, 4, depthMap)
		lightCullingShader.ScreenSize.Set(uniforms.UIVec2{window.Width, window.Height})
		lightCullingShader.LightCount.Set(lights.GetNumPointLights())
		lightCullingShader.LightBuffer.Set(lights.GetPointLightBuffer())
		lightCullingShader.VisibleLightIndicesBuffer.Set(lights.GetPointLightVisibleLightIndicesBuffer())
		gl.DispatchCompute(window.GetNumTilesX(), window.GetNumTilesY(), 1)

		gl.UseProgram(0)

		// Step 3: Normal pass
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		gl.BindProgramPipeline(normalPipeline)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		vertexShader.View.Set(camera.GetView())
		vertexShader.Projection.Set(window.GetProjection())
		fragmentShader.NumTilesX.Set(window.GetNumTilesX())
		fragmentShader.LightBuffer.Set(lights.GetPointLightBuffer())
		fragmentShader.VisibleLightIndicesBuffer.Set(lights.GetPointLightVisibleLightIndicesBuffer())
		fragmentShader.DirectionalLightBuffer.Set(lights.GetDirectionalLightBuffer())
		fragmentShader.Diffuse.Set(gl.TEXTURE0, 0, diffuseTexture)
		gl.BindVertexArray(cubeVao)
		for x := 0; x < 10; x++ {
			for y := 0; y < 10; y++ {
				modelTranslation := mgl32.Translate3D(float32(4*x), 5.0, float32(4*y))
				vertexShader.Model.Set(modelTranslation)
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, 0)

		vertexShader.Model.Set(mgl32.Ident4())
		fragmentShader.Diffuse.Set(gl.TEXTURE0, 0, sandTexture)
		gl.BindVertexArray(planeVao)
		gl.DrawArrays(gl.TRIANGLES, 0, 2*3)

		// PIP
		if pip.Enabled {
			gl.Disable(gl.DEPTH_TEST)
			pip.Render(window.GetProjection())
			gl.Enable(gl.DEPTH_TEST)
		}

		// Maintenance
		w.SwapBuffers()
		glfw.PollEvents()

		window.RecenterCursor()
		framerate.EndOfFrame(timer.GetTime())
	}
}

var planeVertices = []vertexshader.Vertex{
	{mgl32.Vec3{-1000.0, 0, -1000.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{1000.0, 0, -1000.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{0, 50}},
	{mgl32.Vec3{-1000.0, 0, 1000.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{50, 0}},
	{mgl32.Vec3{1000.0, 0, -1000.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{0, 50}},
	{mgl32.Vec3{1000.0, 0, 1000.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{50, 50}},
	{mgl32.Vec3{-1000.0, 0, 1000.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{50, 0}},
}

var cubeVertices = []vertexshader.Vertex{
	//  X, Y, Z
	// Bottom
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec3{0, -1.0, 0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, -1.0, 0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{0, -1.0, 0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, -1.0, 0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{0, -1.0, 0}, mgl32.Vec2{1, 1}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{0, -1.0, 0}, mgl32.Vec2{0, 1}},

	// Top
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{0, 1.0, 0}, mgl32.Vec2{1, 1}},

	// Front
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}, mgl32.Vec2{1, 1}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{0, 0, 1.0}, mgl32.Vec2{1, 1}},

	// Back
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{0, 0, -1.0}, mgl32.Vec2{1, 1}},

	// Left
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{-1.0, 0, 0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{-1.0, 0, 0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{-1.0, -1.0, -1.0}, mgl32.Vec3{-1.0, 0, 0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{-1.0, -1.0, 1.0}, mgl32.Vec3{-1.0, 0, 0}, mgl32.Vec2{0, 1}},
	{mgl32.Vec3{-1.0, 1.0, 1.0}, mgl32.Vec3{-1.0, 0, 0}, mgl32.Vec2{1, 1}},
	{mgl32.Vec3{-1.0, 1.0, -1.0}, mgl32.Vec3{-1.0, 0, 0}, mgl32.Vec2{1, 0}},

	// Right
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{1.0, 0, 0}, mgl32.Vec2{1, 1}},
	{mgl32.Vec3{1.0, -1.0, -1.0}, mgl32.Vec3{1.0, 0, 0}, mgl32.Vec2{1, 0}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{1.0, 0, 0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{1.0, -1.0, 1.0}, mgl32.Vec3{1.0, 0, 0}, mgl32.Vec2{1, 1}},
	{mgl32.Vec3{1.0, 1.0, -1.0}, mgl32.Vec3{1.0, 0, 0}, mgl32.Vec2{0, 0}},
	{mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 0, 0}, mgl32.Vec2{0, 1}},
}
