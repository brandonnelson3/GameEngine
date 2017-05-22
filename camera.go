package main

import (
	"math"

	"github.com/brandonnelson3/GameEngine/input"
	"github.com/brandonnelson3/GameEngine/messagebus"
	"github.com/brandonnelson3/GameEngine/window"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	pi2 = math.Pi / 2.0
)

// FirstPersonCamera is a camera which behaves like a FirstPersonShooter Camera would. WASD control the movement and the mouse controls the direction.
type FirstPersonCamera struct {
	position        mgl32.Vec3
	direction       mgl32.Vec3
	horizontalAngle float32
	verticalAngle   float32
	sensitivity     float32
	speed           float32
}

// NewFirstPersonCamera instantiates a new FirstPersonCamera.
func NewFirstPersonCamera() *FirstPersonCamera {
	c := &FirstPersonCamera{position: mgl32.Vec3{-23, 30, -23}, horizontalAngle: 5.56, verticalAngle: -0.56, sensitivity: 0.001, speed: 20}
	messagebus.RegisterType("key", c.handleMovement)
	messagebus.RegisterType("mouse", c.handleMouse)
	return c
}

// Update is called every frame to execute this frame's movement.
func (c *FirstPersonCamera) Update(d float64) {
	if c.direction.X() != 0 || c.direction.Y() != 0 || c.direction.Z() != 0 {
		delta := c.direction.Normalize().Mul(float32(d) * c.speed)
		c.position = c.position.Add(delta)
		c.direction = mgl32.Vec3{0, 0, 0}
	}
}

// GetForward returns the forward unit vector for this camera.
func (c *FirstPersonCamera) GetForward() mgl32.Vec3 {
	return mgl32.Rotate3DY(c.horizontalAngle).Mul3x1(mgl32.Rotate3DZ(c.verticalAngle).Mul3x1((mgl32.Vec3{1, 0, 0})))
}

// GetRight returns the right unit vector for this camera.
func (c *FirstPersonCamera) GetRight() mgl32.Vec3 {
	return mgl32.Rotate3DY(c.horizontalAngle).Mul3x1(mgl32.Vec3{0, 0, 1})
}

// GetView returns the current view matrix for this camera.
func (c *FirstPersonCamera) GetView() mgl32.Mat4 {
	return mgl32.LookAtV(c.position, c.position.Add(c.GetForward()), mgl32.Vec3{0, 1, 0})
}

func (c *FirstPersonCamera) handleMovement(m *messagebus.Message) {
	direction := mgl32.Vec3{0, 0, 0}
	pressedKeys := m.Data.([]glfw.Key)
	for _, key := range pressedKeys {
		switch key {
		case glfw.KeyW:
			direction = direction.Add(c.GetForward())
		case glfw.KeyS:
			direction = direction.Sub(c.GetForward())
		case glfw.KeyD:
			direction = direction.Add(c.GetRight())
		case glfw.KeyA:
			direction = direction.Sub(c.GetRight())
		}
	}
	c.direction = direction
}

func (c *FirstPersonCamera) handleMouse(m *messagebus.Message) {
	mouseInput := m.Data.(input.MouseInput)
	c.verticalAngle -= c.sensitivity * float32(mouseInput.Y-float64(window.Height)/2)
	if c.verticalAngle < -pi2 {
		c.verticalAngle = float32(-pi2 + 0.0001)
	}
	c.horizontalAngle -= c.sensitivity * float32(mouseInput.X-float64(window.Width)/2)
	for c.horizontalAngle < 0 {
		c.horizontalAngle += float32(2 * math.Pi)
	}
}
