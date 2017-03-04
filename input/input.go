package input

import (
	"fmt"

	"github.com/go-gl/glfw/v3.1/glfw"
)

// KeyCallBack is the function bound to handle key events from OpenGL.
func KeyCallBack(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("Got key press: %v\n", key)
}

// MouseButtonCallback is the function bound to handle mouse button events from OpenGL.
func MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	fmt.Printf("Got mouse press: %v\n", button)
}

// CursorPosCallback is the function bound to handle mouse movement events from OpenGL.
func CursorPosCallback(w *glfw.Window, xpos, ypos float64) {
	fmt.Printf("Got mouse movement: [%f, %f]\n", xpos, ypos)
}
