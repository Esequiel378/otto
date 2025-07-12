package player

import (
	"otto/system/input"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl64"
)

type InputPlayerMovement struct {
	Velocity mgl64.Vec3
}

var _ input.Context = (*InputPlayerMovement)(nil)

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

	return h.Velocity != (mgl64.Vec3{})
}

// InputPlayerCamera handles camera movement input
type InputPlayerCamera struct {
	Rotation mgl64.Vec2 // Pitch, Yaw
	Zoom     float64
}

// GetType returns the type identifier for camera control input
func (c *InputPlayerCamera) GetType() string {
	return "player_camera"
}

// Process handles camera control input
func (c *InputPlayerCamera) Process() bool {
	// Reset state at the beginning
	c.Rotation = mgl64.Vec2{}
	c.Zoom = 0.0

	// Check for camera control keys
	rotation := mgl64.Vec2{}
	zoom := 0.0

	// Get mouse input
	io := imgui.CurrentIO()
	if io == nil {
		return false
	}

	// Check if right mouse button is held down for camera rotation
	rightMouseDown := imgui.IsMouseDown(1) // Right mouse button

	if rightMouseDown {
		// Get mouse delta for rotation
		mouseDelta := io.MouseDelta()
		deltaX := float64(mouseDelta.X)
		deltaY := float64(mouseDelta.Y)

		// Apply mouse sensitivity
		sensitivity := 0.1
		rotation[0] = deltaY * sensitivity // Pitch (Y axis)
		rotation[1] = deltaX * sensitivity // Yaw (X axis)
	}

	// Zoom controls with mouse wheel (always process, not just when right mouse is down)
	mouseWheel := io.MouseWheel()
	if mouseWheel != 0 {
		zoom = float64(mouseWheel) * 0.1
	}

	// Also support keyboard zoom controls as fallback
	plusPressed := imgui.IsKeyDown(imgui.KeyEqual)
	minusPressed := imgui.IsKeyDown(imgui.KeyMinus)

	if plusPressed {
		zoom += 1.0 // Zoom in
	}
	if minusPressed {
		zoom -= 1.0 // Zoom out
	}

	// Update state
	c.Rotation = rotation
	c.Zoom = zoom

	return rightMouseDown || c.Rotation != (mgl64.Vec2{}) || c.Zoom != 0.0
}
