package receiver

import (
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl64"
)

func (h *InputPlayerMovement) Handle() {
	if ctx := imgui.CurrentContext(); ctx == nil {
		return
	}

	io := imgui.CurrentIO()

	// Don't process keyboard input if ImGui wants to capture it
	if io.WantCaptureKeyboard() {
		return
	}

	// Process keyboard input
	if imgui.IsKeyDown(imgui.KeyW) {
		h.Velocity = mgl64.Vec3{0, 0, 1}
	}

	if imgui.IsKeyDown(imgui.KeyS) {
		h.Velocity = mgl64.Vec3{0, 0, -1}
	}

	if imgui.IsKeyDown(imgui.KeyA) {
		h.Velocity = mgl64.Vec3{-1, 0, 0}
	}

	if imgui.IsKeyDown(imgui.KeyD) {
		h.Velocity = mgl64.Vec3{1, 0, 0}
	}

	if imgui.IsKeyDown(imgui.KeySpace) {
		h.Velocity = mgl64.Vec3{0, 1, 0}
	}

	if imgui.IsKeyDown(imgui.KeyLeftShift) {
		h.Velocity = mgl64.Vec3{0, -1, 0}
	}
}
