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

// BeginningOfFrame is expected to be called at the same point in every frame to work properly.
func BeginningOfFrame() {
	time := glfw.GetTime()
	previousElapsed = time - previousTime
	previousTime = time
}

// GetPreviousFrameLength returns the time in seconds as a float64 of the previous frame.
func GetPreviousFrameLength() float64 {
	return previousElapsed
}

// GetTime returns the current time.Now().
func GetTime() float64 {
	return glfw.GetTime()
}
