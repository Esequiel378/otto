package camera

import (
	"otto/system/input"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// InputCamera handles camera movement input using Euler angles like the reference code
type InputCamera struct {
	PID          *actor.PID
	lastRotation mgl64.Vec2 // Track last mouse position for delta calculation
	rotation     mgl64.Vec2 // Pitch, Yaw
	fov          float64
	sensitivity  float64
}

var _ input.Context = (*InputCamera)(nil)

// NewInputCamera creates a new InputCamera with default values
func NewInputCamera(pid *actor.PID) *InputCamera {
	return &InputCamera{
		PID:         pid,
		sensitivity: 0.5,
	}
}

// GetPID returns the PID of the input context
func (c *InputCamera) GetPID() *actor.PID {
	return c.PID
}

func (c *InputCamera) Rotation() mgl64.Vec2 {
	return c.rotation
}

// Process handles camera control input using the input state
func (c *InputCamera) Process(deltaTime float64, state *input.InputState, captureKeyboard, captureMouse bool) bool {
	// Only process mouse input if UI doesn't want to capture it
	if captureMouse {
		return false
	}

	c.handleCameraRotation(deltaTime, state)
	c.handleCameraZoom(deltaTime, state)

	return c.rotation != (mgl64.Vec2{0, 0}) || c.fov != 0.0
}

func (c *InputCamera) handleCameraRotation(deltaTime float64, state *input.InputState) {
	c.rotation = mgl64.Vec2{0, 0}

	// Check if right mouse button is held down for camera rotation
	if !state.IsMouseButtonPressed(input.MouseButtonRight) {
		return
	}

	// Get mouse delta for rotation
	mouseDelta := state.GetMouseDelta()
	deltaX := mouseDelta.X()
	deltaY := mouseDelta.Y()

	offset := mgl64.Vec2{deltaX, deltaY}.Sub(c.lastRotation)

	// Update last position
	c.lastRotation = mgl64.Vec2{deltaX, deltaY}

	if c.lastRotation.X() == 0 && c.lastRotation.Y() == 0 {
		return
	}

	c.rotation = c.rotation.Add(offset.Mul(c.sensitivity * deltaTime))

	// Clamp pitch between -89 and 89 degrees
	if c.rotation.Y() > 89.0 {
		c.rotation[1] = 89.0
	}
	if c.rotation.Y() < -89.0 {
		c.rotation[1] = -89.0
	}
}

func (c *InputCamera) handleCameraZoom(deltaTime float64, state *input.InputState) {
	mouseWheel := state.GetMouseWheel()
	if mouseWheel != 0 {
		c.fov = mouseWheel * c.sensitivity * deltaTime
	}
}
