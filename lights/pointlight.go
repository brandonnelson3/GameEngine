package lights

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// MaximumPointLights is the maximum number of lights that the pointlight system is prepared to handle.
	MaximumPointLights = 1024
)

var (
	// PointLights are the current pointlights in the scene.
	PointLights [MaximumPointLights]PointLight
)

func init() {
	PointLights[0].Color = mgl32.Vec4{1.0, 0.0, 0.0, 1.0}
	PointLights[0].PositionAndRadius = mgl32.Vec4{0, 3, 0, 10}

	PointLights[1].Color = mgl32.Vec4{0.0, 1.0, 0.0, 1.0}
	PointLights[1].PositionAndRadius = mgl32.Vec4{36, 3, 36, 10}

	PointLights[2].Color = mgl32.Vec4{1.0, 1.0, 0.0, 1.0}
	PointLights[2].PositionAndRadius = mgl32.Vec4{36, 3, 0, 10}

	PointLights[3].Color = mgl32.Vec4{0.0, 0.0, 1.0, 1.0}
	PointLights[3].PositionAndRadius = mgl32.Vec4{0, 3, 36, 10}
}

// PointLight represents all of the data about a PointLight.
type PointLight struct {
	Color, PositionAndRadius mgl32.Vec4
}

// VisibleIndex is a wrapper around an index.
type VisibleIndex struct {
	index int32
}
