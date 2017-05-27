package window

import (
	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	TileSize = 16
)

var (
	// Width is the width of the window.
	Width = uint32(1920)

	// Height is the height of the window.
	Height = uint32(1080)

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
	w, err := glfw.CreateWindow(int(Width), int(Height), Title, nil, nil)
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

// GetNumTilesX returns back the number of tiles in each the X dimension that are needed for the current window size.
func GetNumTilesX() uint32 {
	return (Width + TileSize - 1) / TileSize
}

// GetNumTilesY returns back the number of tiles in each the Y dimension that are needed for the current window size.
func GetNumTilesY() uint32 {
	return (Height + TileSize - 1) / TileSize
}

// GetTotalNumTiles returns back the total number of tiles required to cover the entire screen.
func GetTotalNumTiles() uint32 {
	return GetNumTilesX() * GetNumTilesY()
}
