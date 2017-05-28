package framerate

import (
	"fmt"
	"sync"
	"time"

	"github.com/brandonnelson3/GameEngine/messagebus"
)

const (
	numAveragedFrameLengths = 25
	framerateCap            = 105
)

var (
	times  [numAveragedFrameLengths]float64
	frames int
	i      = 0
	mu     sync.Mutex

	thisFrameStart float64
)

func init() {
	go log()
}

func log() {
	for range time.Tick(time.Millisecond * 500) {
		if frames < numAveragedFrameLengths {
			continue
		}

		averageFrameTime, averageFramesPerSecond := calculateFrameDetails()
		messagebus.SendSync(&messagebus.Message{System: "FrameRate", Type: "log", Data: fmt.Sprintf("Length: %.3f ms - Avg FPS: %.1f - Limiting framerate to %d", averageFrameTime*1000, averageFramesPerSecond, framerateCap)})
	}
}

func calculateFrameDetails() (float64, float64) {
	totalTime := float64(0)
	mu.Lock()
	for _, t := range times {
		totalTime += t
	}
	mu.Unlock()
	averageFrameTime := totalTime / numAveragedFrameLengths
	averageFramesPerSecond := 1 / averageFrameTime

	return averageFrameTime, averageFramesPerSecond
}

// BeginningOfFrame is intended to be called at the beginning of the frame.
func BeginningOfFrame(now float64) {
	thisFrameStart = now
}

// EndOfFrame is intended to be called at the end of the frame with the current time. This function will block to maintain the intended maximum frame rate.
func EndOfFrame(now float64) {
	d := now - thisFrameStart
	mu.Lock()
	times[i] = d
	frames++
	i++
	if i >= numAveragedFrameLengths {
		i = 0
	}
	mu.Unlock()
	// Sleep for as long as we need to...
	time.Sleep(time.Duration((float64(time.Second) / framerateCap) - (float64(time.Second) * d)))
}
