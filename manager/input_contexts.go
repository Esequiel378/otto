package manager

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
	// Check for camera control keys
	rotation := mgl64.Vec2{}
	zoom := 0.0

	// Test with different keys to see if the issue is with specific key constants
	// Try using arrow keys instead of I/J/K/L
	upPressed := imgui.IsKeyDown(imgui.KeyUpArrow)
	downPressed := imgui.IsKeyDown(imgui.KeyDownArrow)
	leftPressed := imgui.IsKeyDown(imgui.KeyLeftArrow)
	rightPressed := imgui.IsKeyDown(imgui.KeyRightArrow)

	// Also test the original keys
	iPressed := imgui.IsKeyDown(imgui.KeyI)
	kPressed := imgui.IsKeyDown(imgui.KeyK)
	jPressed := imgui.IsKeyDown(imgui.KeyJ)
	lPressed := imgui.IsKeyDown(imgui.KeyL)
	plusPressed := imgui.IsKeyDown(imgui.KeyEqual)
	minusPressed := imgui.IsKeyDown(imgui.KeyMinus)

	// Rotation controls (try arrow keys first, then I/J/K/L)
	if upPressed || iPressed {
		rotation[0] -= 1.0 // Pitch down
	}
	if downPressed || kPressed {
		rotation[0] += 1.0 // Pitch up
	}
	if leftPressed || jPressed {
		rotation[1] -= 1.0 // Yaw left
	}
	if rightPressed || lPressed {
		rotation[1] += 1.0 // Yaw right
	}

	// Zoom controls (+/- keys)
	if plusPressed || imgui.IsKeyDown(imgui.KeyKeypadAdd) {
		zoom += 1.0 // Zoom in
	}
	if minusPressed || imgui.IsKeyDown(imgui.KeyKeypadSubtract) {
		zoom -= 1.0 // Zoom out
	}

	// Update state
	ic.Rotation = rotation
	ic.Zoom = zoom

	// Return true if any input was detected
	return rotation != (mgl64.Vec2{}) || zoom != 0.0
}

// InputUIInteraction handles UI-related input
type InputUIInteraction struct {
	MousePosition mgl64.Vec2
	MouseClicked  bool
	KeyPressed    string
}

// GetType returns the type identifier for UI interaction input
func (u *InputUIInteraction) GetType() string {
	return "ui_interaction"
}

// Process handles UI interaction input
func (u *InputUIInteraction) Process() bool {
	u.MouseClicked = false
	u.KeyPressed = ""

	// Get mouse position
	io := imgui.CurrentIO()
	u.MousePosition = mgl64.Vec2{float64(io.MousePos().X), float64(io.MousePos().Y)}

	// Check for mouse clicks (simplified - using key instead)
	if imgui.IsKeyDown(imgui.KeyEscape) {
		u.KeyPressed = "escape"
	}
	if imgui.IsKeyDown(imgui.KeyEnter) {
		u.KeyPressed = "enter"
	}
	if imgui.IsKeyDown(imgui.KeyTab) {
		u.KeyPressed = "tab"
	}

	return u.KeyPressed != ""
}

// InputGameActions handles game-specific action input
type InputGameActions struct {
	Action string
	Active bool
}

// GetType returns the type identifier for game actions input
func (g *InputGameActions) GetType() string {
	return "game_actions"
}

// Process handles game action input
func (g *InputGameActions) Process() bool {
	g.Action = ""
	g.Active = false

	// Check for action keys
	if imgui.IsKeyDown(imgui.KeyF) {
		g.Action = "interact"
		g.Active = true
	}
	if imgui.IsKeyDown(imgui.KeyR) {
		g.Action = "reload"
		g.Active = true
	}
	if imgui.IsKeyDown(imgui.KeyE) {
		g.Action = "use"
		g.Active = true
	}
	if imgui.IsKeyDown(imgui.KeyQ) {
		g.Action = "drop"
		g.Active = true
	}

	return g.Active
}
