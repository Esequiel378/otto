package input

import (
	"log"
	"sync"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/mathgl/mgl64"
)

// ImGuiProvider implements InputProvider using ImGui as the backend
type ImGuiProvider struct {
	inputState *InputState
	mu         sync.RWMutex
	lastValid  bool
}

// NewImGuiProvider creates a new ImGui input provider
func NewImGuiProvider() *ImGuiProvider {
	return &ImGuiProvider{
		inputState: NewInputState(),
	}
}

// GetInputState returns the current input state
func (p *ImGuiProvider) GetInputState() *InputState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.inputState
}

// Update updates the input state with current ImGui input data
func (p *ImGuiProvider) Update() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if ImGui context is valid with proper error handling
	if !p.isValidContext() {
		// If context was previously valid but now isn't, log it
		if p.lastValid {
			log.Printf("ImGui context became invalid, skipping input update")
		}
		p.lastValid = false
		return nil
	}

	p.lastValid = true

	// Safely get IO with error handling
	io := p.getIO()
	if io == nil {
		return nil
	}

	// Update keyboard state with error handling
	p.updateKeyboardState(io)

	// Update mouse state with error handling
	p.updateMouseState(io)

	return nil
}

// isValidContext checks if ImGui context is valid with proper error handling
func (p *ImGuiProvider) isValidContext() bool {
	defer func() {
		if r := recover(); r != nil {
			// ImGui context is invalid or destroyed
		}
	}()

	// Try to get current context
	ctx := imgui.CurrentContext()
	return ctx != nil
}

// getIO safely gets the ImGui IO with error handling
func (p *ImGuiProvider) getIO() *imgui.IO {
	defer func() {
		if r := recover(); r != nil {
			// ImGui IO is not available
		}
	}()

	return imgui.CurrentIO()
}

// updateKeyboardState updates keyboard state with error handling
func (p *ImGuiProvider) updateKeyboardState(io *imgui.IO) {
	defer func() {
		if r := recover(); r != nil {
			// Keyboard state update failed
		}
	}()

	p.inputState.keyStates[KeyW] = imgui.IsKeyDown(imgui.KeyW)
	p.inputState.keyStates[KeyS] = imgui.IsKeyDown(imgui.KeyS)
	p.inputState.keyStates[KeyA] = imgui.IsKeyDown(imgui.KeyA)
	p.inputState.keyStates[KeyD] = imgui.IsKeyDown(imgui.KeyD)
	p.inputState.keyStates[KeySpace] = imgui.IsKeyDown(imgui.KeySpace)
	p.inputState.keyStates[KeyLeftShift] = imgui.IsKeyDown(imgui.KeyLeftShift)
	p.inputState.keyStates[KeyEqual] = imgui.IsKeyDown(imgui.KeyEqual)
	p.inputState.keyStates[KeyMinus] = imgui.IsKeyDown(imgui.KeyMinus)
}

// updateMouseState updates mouse state with error handling
func (p *ImGuiProvider) updateMouseState(io *imgui.IO) {
	defer func() {
		if r := recover(); r != nil {
			// Mouse state update failed
		}
	}()

	// Update mouse button state
	p.inputState.mouseButtonStates[MouseButtonLeft] = imgui.IsMouseDown(0)
	p.inputState.mouseButtonStates[MouseButtonRight] = imgui.IsMouseDown(1)
	p.inputState.mouseButtonStates[MouseButtonMiddle] = imgui.IsMouseDown(2)

	// Update mouse position and delta
	mousePos := io.MousePos()
	p.inputState.mousePosition = mgl64.Vec2{float64(mousePos.X), float64(mousePos.Y)}

	mouseDelta := io.MouseDelta()
	p.inputState.mouseDelta = mgl64.Vec2{float64(mouseDelta.X), float64(mouseDelta.Y)}

	// Update mouse wheel
	p.inputState.mouseWheel = float64(io.MouseWheel())

	// Update input capture flags
	p.inputState.wantCaptureKeyboard = io.WantCaptureKeyboard()
	p.inputState.wantCaptureMouse = io.WantCaptureMouse()
}

// IsValid returns true if the ImGui context is valid
func (p *ImGuiProvider) IsValid() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastValid
}
