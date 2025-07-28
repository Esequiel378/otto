package camera

import (
	"otto/system/input"

	"math"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// InputCamera handles camera movement input using Euler angles like the reference code
type InputCamera struct {
	PID          *actor.PID
	lastRotation mgl64.Vec2 // Track last mouse position for delta calculation
	yaw          float64    // Euler angle for horizontal rotation
	pitch        float64    // Euler angle for vertical rotation
	Zoom         float64
}

var _ input.Context = (*InputCamera)(nil)

// NewInputCamera creates a new InputCamera with default values
func NewInputCamera(pid *actor.PID) *InputCamera {
	return &InputCamera{
		PID: pid,
	}
}

// GetPID returns the PID of the input context
func (c *InputCamera) GetPID() *actor.PID {
	return c.PID
}

// GetYaw returns the current yaw angle
func (c *InputCamera) GetYaw() float64 {
	return c.yaw
}

// GetPitch returns the current pitch angle
func (c *InputCamera) GetPitch() float64 {
	return c.pitch
}

// GetFrontVector returns the front direction vector calculated from Euler angles
func (c *InputCamera) GetFrontVector() mgl64.Vec3 {
	// Calculate the new Front vector (like reference code)
	front := mgl64.Vec3{
		math.Cos(mgl64.DegToRad(c.yaw)) * math.Cos(mgl64.DegToRad(c.pitch)),
		math.Sin(mgl64.DegToRad(c.pitch)),
		math.Sin(mgl64.DegToRad(c.yaw)) * math.Cos(mgl64.DegToRad(c.pitch)),
	}.Normalize()

	return front
}

// Process handles camera control input using the input state
func (c *InputCamera) Process(deltaTime float64, state *input.InputState, captureKeyboard, captureMouse bool) bool {
	// Reset zoom at the beginning
	c.Zoom = 0.0

	// Check if right mouse button is held down for camera rotation
	// Only process mouse input if UI doesn't want to capture it
	rightMouseDown := !captureMouse && state.IsMouseButtonPressed(input.MouseButtonRight)

	if rightMouseDown {
		// Get mouse delta for rotation
		mouseDelta := state.GetMouseDelta()
		deltaX := mouseDelta.X()
		deltaY := mouseDelta.Y()

		// Calculate offset from last position (similar to reference code)
		offsetX := deltaX - c.lastRotation.X()
		offsetY := c.lastRotation.Y() - deltaY

		// Update last position
		c.lastRotation = mgl64.Vec2{deltaX, deltaY}

		if c.lastRotation.X() != 0 || c.lastRotation.Y() != 0 {
			sensitivity := 0.5
			c.yaw += offsetX * sensitivity * deltaTime
			c.pitch += offsetY * sensitivity * deltaTime

			// Clamp pitch between -89 and 89 degrees (like reference code)
			if c.pitch > 89.0 {
				c.pitch = 89.0
			}
			if c.pitch < -89.0 {
				c.pitch = -89.0
			}
		} else {
			c.yaw = 0
			c.pitch = 0
		}
	}

	// Zoom controls with mouse wheel (always process, not just when right mouse is down)
	// Only process mouse wheel if UI doesn't want to capture mouse
	if !captureMouse {
		mouseWheel := state.GetMouseWheel()
		if mouseWheel != 0 {
			// Apply deltaTime to zoom for consistent speed
			c.Zoom = mouseWheel * 0.1 * deltaTime
		}
	}

	return rightMouseDown || c.Zoom != 0.0
}
