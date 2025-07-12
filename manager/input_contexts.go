package manager

import (
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl64"
)

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
