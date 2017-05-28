package lights

import (
	"sync"
	"unsafe"

	"github.com/brandonnelson3/GameEngine/window"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// MaximumPointLights is the maximum number of lights that the pointlight system is prepared to handle.
	MaximumPointLights = 1024
)

var (
	// PointLights are the current pointlights in the scene.
	PointLights    [MaximumPointLights]PointLight
	numPointLights = uint32(0)
	mu             sync.Mutex

	lightBuffer, visibleLightIndicesBuffer uint32
)

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

// InitPointLights sets up buffer space for light culling calculations and storage.
func InitPointLights() {
	AddPointLight(mgl32.Vec3{0, 12, 0}, mgl32.Vec3{1, 0, 0}, 1.0, 10.0)
	AddPointLight(mgl32.Vec3{36, 12, 0}, mgl32.Vec3{0, 1, 0}, 1.0, 10.0)
	AddPointLight(mgl32.Vec3{0, 12, 36}, mgl32.Vec3{0, 0, 1}, 1.0, 10.0)
	AddPointLight(mgl32.Vec3{36, 12, 36}, mgl32.Vec3{1, 1, 0}, 1.0, 10.0)

	// Prepare light buffers
	gl.GenBuffers(1, &lightBuffer)
	gl.GenBuffers(1, &visibleLightIndicesBuffer)

	// Bind light buffer
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, lightBuffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, MaximumPointLights*int(unsafe.Sizeof(&PointLight{})), unsafe.Pointer(&PointLights), gl.DYNAMIC_DRAW)

	// Bind visible light indices buffer
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, visibleLightIndicesBuffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(window.GetTotalNumTiles())*int(unsafe.Sizeof(&VisibleIndex{}))*MaximumPointLights, nil, gl.STATIC_DRAW)

	// Unbind for safety.
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)
}

// GetNumPointLights returns the number of PointLights that are currently in the scene.
func GetNumPointLights() uint32 {
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

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, lightBuffer)
	//pointLights := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_WRITE)
	/*
		for (int i = 0; i < NUM_LIGHTS; i++) {
			PointLight &light = pointLights[i];
			float min = LIGHT_MIN_BOUNDS[1];
			float max = LIGHT_MAX_BOUNDS[1];

			light.position.y = fmod((light.position.y + (-4.5f * deltaTime) - min + max), max) + min;
		}
	*/
	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, 0)

	numPointLights++
}

// GetPointLightBuffer retrieves the private lightBuffer variable.
func GetPointLightBuffer() uint32 {
	return lightBuffer
}

// GetPointLightVisibleLightIndicesBuffer retrieves the private visibleLightIndicesBuffer variable.
func GetPointLightVisibleLightIndicesBuffer() uint32 {
	return visibleLightIndicesBuffer
}
