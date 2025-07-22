package camera

import (
	"otto/system/input"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// InputCamera handles camera movement input
type InputCamera struct {
	PID      *actor.PID
	Rotation mgl64.Vec2 // Pitch, Yaw
	Zoom     float64
}

var _ input.Context = (*InputCamera)(nil)

// GetPID returns the PID of the input context
func (c *InputCamera) GetPID() *actor.PID {
	return c.PID
}

// Process handles camera control input using the input state
func (c *InputCamera) Process(state *input.InputState, captureKeyboard, captureMouse bool) bool {
	// Reset state at the beginning
	c.Rotation = mgl64.Vec2{}
	c.Zoom = 0.0

	// Check for camera control keys
	rotation := mgl64.Vec2{}
	zoom := 0.0

	// Check if right mouse button is held down for camera rotation
	// Only process mouse input if UI doesn't want to capture it
	rightMouseDown := !captureMouse && state.IsMouseButtonPressed(input.MouseButtonRight)

	if rightMouseDown {
		// Get mouse delta for rotation
		mouseDelta := state.GetMouseDelta()
		deltaX := mouseDelta.X()
		deltaY := mouseDelta.Y()

		// Apply mouse sensitivity - direct input, no smoothing
		sensitivity := 0.2
		rotation[0] = deltaY * sensitivity // Pitch (Y axis)
		rotation[1] = deltaX * sensitivity // Yaw (X axis)
	}

	// Zoom controls with mouse wheel (always process, not just when right mouse is down)
	// Only process mouse wheel if UI doesn't want to capture mouse
	if !captureMouse {
		mouseWheel := state.GetMouseWheel()
		if mouseWheel != 0 {
			zoom = mouseWheel * 0.1
		}
	}

	// Update state
	c.Rotation = rotation
	c.Zoom = zoom

	return rightMouseDown || c.Rotation != (mgl64.Vec2{}) || c.Zoom != 0.0
}
