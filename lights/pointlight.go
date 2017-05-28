package lights

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	// MaximumPointLights is the maximum number of lights that the pointlight system is prepared to handle.
	MaximumPointLights = 1024
)

var (
	// PointLights are the current pointlights in the scene.
	PointLights    [MaximumPointLights]PointLight
	numPointLights uint
	mu             sync.Mutex
)

func init() {
	PointLights[0].Color = mgl32.Vec3{1.0, 0.0, 0.0}
	PointLights[0].Intensity = 1.0
	PointLights[0].Position = mgl32.Vec3{0, 3, 0}
	PointLights[0].Radius = 10.0

	PointLights[1].Color = mgl32.Vec3{0.0, 1.0, 0.0}
	PointLights[1].Intensity = 1.0
	PointLights[1].Position = mgl32.Vec3{36, 3, 0}
	PointLights[1].Radius = 10.0

	PointLights[2].Color = mgl32.Vec3{0.0, 0.0, 1.0}
	PointLights[2].Intensity = 1.0
	PointLights[2].Position = mgl32.Vec3{0, 3, 36}
	PointLights[2].Radius = 10.0

	PointLights[3].Color = mgl32.Vec3{1.0, 1.0, 0.0}
	PointLights[3].Intensity = 1.0
	PointLights[3].Position = mgl32.Vec3{36, 3, 36}
	PointLights[3].Radius = 10.0
}

// PointLight represents all of the data about a PointLight.
type PointLight struct {
	Color     mgl32.Vec3
	Intensity float32
	Position  mgl32.Vec3
	Radius    float32
}

// VisibleIndex is a wrapper around an index.
type VisibleIndex struct {
	index int32
}

// GetNumPointLights returns the number of PointLights that are currently in the scene.
func GetNumPointLights() uint {
	return numPointLights
}

// AddPointLight adds a PointLight to the scene with the given attributes.
func AddPointLight(position, color mgl32.Vec3, intensity, radius float32) {
	mu.Lock()
	defer mu.Unlock()

	PointLights[numPointLights].Color = color
	PointLights[numPointLights].Intensity = intensity
	PointLights[numPointLights].Position = position
	PointLights[numPointLights].Radius = radius

	numPointLights++
}
