package timer

import "github.com/go-gl/glfw/v3.1/glfw"

var (
	previousTime    float64
	previousElapsed float64
)

func init() {
	previousTime = 0
	previousElapsed = 0
}

// Update is expected to be called at the same point in every frame to work properly.
func Update() {
	time := glfw.GetTime()
	previousElapsed = time - previousTime
	previousTime = time
}

// GetPreviousFrameLength returns the time in seconds as a float64 of the previous frame.
func GetPreviousFrameLength() float64 {
	return previousElapsed
}
