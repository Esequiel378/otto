package player

import (
	"otto/system/input"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type InputPlayerMovement struct {
	PID      *actor.PID
	Velocity mgl64.Vec3
}

var _ input.Context = (*InputPlayerMovement)(nil)

// GetPID returns the PID of the input context
func (h *InputPlayerMovement) GetPID() *actor.PID {
	return h.PID
}

// Process handles player movement input using the input state
func (h *InputPlayerMovement) Process(deltaTime float64, state *input.InputState, captureKeyboard, captureMouse bool) bool {
	h.Velocity = mgl64.Vec3{} // Reset velocity

	// Skip keyboard input if UI wants to capture it
	if captureKeyboard {
		return false
	}

	// Process keyboard input - Allow multiple keys to be pressed simultaneously
	if state.IsKeyPressed(input.KeyW) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 0, 1}) // Forward (Z+)
	}

	if state.IsKeyPressed(input.KeyS) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 0, -1}) // Backward (Z-)
	}

	if state.IsKeyPressed(input.KeyA) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{1, 0, 0}) // Left (X+)
	}

	if state.IsKeyPressed(input.KeyD) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{-1, 0, 0}) // Right (X-)
	}

	if state.IsKeyPressed(input.KeySpace) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, -1, 0}) // Up (Y-)
	}

	if state.IsKeyPressed(input.KeyLeftShift) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 1, 0}) // Down (Y+)
	}

	// Normalize the velocity to prevent faster diagonal movement
	if h.Velocity.Len() > 0 {
		h.Velocity = h.Velocity.Normalize()
	}

	return h.Velocity != (mgl64.Vec3{})
}
