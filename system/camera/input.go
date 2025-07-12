package camera

import (
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl64"
)

// InputCameraControl handles camera movement input
type InputCameraControl struct {
	Rotation mgl64.Vec2 // Pitch, Yaw
	Zoom     float64
}

// GetType returns the type identifier for camera control input
func (c *InputCameraControl) GetType() string {
	return "camera_control"
}

// Process handles camera control input
func (ic *InputCameraControl) Process() bool {
	// Reset state at the beginning
	ic.Rotation = mgl64.Vec2{}
	ic.Zoom = 0.0

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
	ic.Rotation = rotation
	ic.Zoom = zoom

	// Return true if right mouse button is pressed (for state tracking)
	// The actual rotation/zoom values are stored in the context for processing
	return rightMouseDown
}
