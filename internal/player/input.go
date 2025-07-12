package player

import (
	"otto/manager"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl64"
)

type InputPlayerMovement struct {
	Velocity mgl64.Vec3
}

var _ manager.InputContext = (*InputPlayerMovement)(nil)

// GetType returns the type identifier for player movement input
func (h *InputPlayerMovement) GetType() string {
	return "player_movement"
}

// Process handles player movement input
func (h *InputPlayerMovement) Process() bool {
	h.Velocity = mgl64.Vec3{} // Reset velocity

	// Process keyboard input - Allow multiple keys to be pressed simultaneously
	if imgui.IsKeyDown(imgui.KeyW) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 0, 1}) // Forward (Z+)
	}

	if imgui.IsKeyDown(imgui.KeyS) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 0, -1}) // Backward (Z-)
	}

	if imgui.IsKeyDown(imgui.KeyA) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{1, 0, 0}) // Left (X+)
	}

	if imgui.IsKeyDown(imgui.KeyD) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{-1, 0, 0}) // Right (X-)
	}

	if imgui.IsKeyDown(imgui.KeySpace) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, -1, 0}) // Up (Y-)
	}

	if imgui.IsKeyDown(imgui.KeyLeftShift) {
		h.Velocity = h.Velocity.Add(mgl64.Vec3{0, 1, 0}) // Down (Y+)
	}

	// Normalize the velocity to prevent faster diagonal movement
	if h.Velocity.Len() > 0 {
		h.Velocity = h.Velocity.Normalize()
	}

	// Return true if any movement key is pressed (for state tracking)
	// The actual velocity is stored in the context for processing
	return h.Velocity != (mgl64.Vec3{})
}

// Legacy Handle method for backward compatibility
func (h *InputPlayerMovement) Handle() {
	h.Process()
}
