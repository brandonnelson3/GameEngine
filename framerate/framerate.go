package framerate

import (
	"fmt"
	"sync"
	"time"
)

const (
	numAveragedFrameLengths = 25
)

var (
	times  [numAveragedFrameLengths]float64
	frames int
	i      = 0
	mu     sync.Mutex
)

func init() {
	go log()
}

// Update is intended to be called at the same point in every frame.
func Update(p float64) {
	mu.Lock()
	times[i] = p
	frames++
	i++
	if i >= numAveragedFrameLengths {
		i = 0
	}
	mu.Unlock()
}

func log() {
	for range time.Tick(time.Millisecond * 500) {
		if frames < numAveragedFrameLengths {
			continue
		}
		totalTime := float64(0)
		mu.Lock()
		for _, t := range times {
			totalTime += t
		}
		mu.Unlock()
		averageFrameTime := totalTime / numAveragedFrameLengths
		averageFramesPerSecond := 1 / averageFrameTime

		fmt.Printf("FPS - Length: %.3f ms - Avg: %.1f", averageFrameTime*1000, averageFramesPerSecond)
	}
}
