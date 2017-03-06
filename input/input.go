package input

import (
	"fmt"

	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const keyRange = 349

var (
	down          [keyRange]bool
	downThisFrame [keyRange]bool
)

// Update calls all of the currently pressed keys.
func Update() {
	// glfw.KeySpace is the lowest key.
	for i := glfw.KeySpace; i < keyRange; i++ {
		if down[i] {
			messagebus.SendSync(&messagebus.Message{Type: "key", Data: i})
		}
	}
}

// KeyCallBack is the function bound to handle key events from OpenGL.
func KeyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		down[key] = true
		downThisFrame[key] = true
	}
	if action == glfw.Release {
		down[key] = false
		downThisFrame[key] = false
	}
}

// MouseButtonCallback is the function bound to handle mouse button events from OpenGL.
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	fmt.Printf("Got mouse press: %v\n", button)
}

// CursorPosCallback is the function bound to handle mouse movement events from OpenGL.
func CursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	fmt.Printf("Got mouse movement: [%f, %f]\n", xpos, ypos)
}
