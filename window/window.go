package window

import (
	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	// Width is the width of the window.
	Width = 800

	// Height is the height of the window.
	Height = 600

	// Title is the title of the window.
	Title = "Game Engine Demo"

	// Near is the near plane of this window.
	Near = float32(0.1)

	// Far is the far plane of this window.
	Far = float32(1000.0)

	// Fov is the field of view.
	Fov = float32(45.0)

	window *glfw.Window
)

// Create creates a new window.
func Create() *glfw.Window {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	w, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}
	w.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	window = w
	messagebus.RegisterType("key", handleEscape)
	return w
}

// GetProjection returns the projection matrix.
func GetProjection() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(Fov), float32(Width)/float32(Height), Near, Far)
}

// RecenterCursor recenters the mouse in this window.
func RecenterCursor() {
	window.SetCursorPos(float64(Width)/2, float64(Height)/2)
}

func handleEscape(m *messagebus.Message) {
	pressedKeys := m.Data.([]glfw.Key)

	for _, key := range pressedKeys {
		if key == glfw.KeyEscape {
			window.SetShouldClose(true)
		}
	}
}
