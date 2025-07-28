package input

import (
	"github.com/go-gl/mathgl/mgl64"
)

// Key represents a keyboard key
type Key int

// MouseButton represents a mouse button
type MouseButton int

// Common key constants
const (
	KeyW Key = iota
	KeyS
	KeyA
	KeyD
	KeySpace
	KeyLeftShift
	KeyEqual
	KeyMinus
	// Add more keys as needed
)

// Common mouse button constants
const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

// InputState represents the current state of all input devices
type InputState struct {
	// Keyboard state
	keyStates map[Key]bool

	// Mouse state
	mouseButtonStates map[MouseButton]bool
	mousePosition     mgl64.Vec2
	mouseDelta        mgl64.Vec2
	mouseWheel        float64

	// Input capture flags
	wantCaptureKeyboard bool
	wantCaptureMouse    bool
}

// NewInputState creates a new input state
func NewInputState() *InputState {
	return &InputState{
		keyStates:           make(map[Key]bool),
		mouseButtonStates:   make(map[MouseButton]bool),
		mousePosition:       mgl64.Vec2{},
		mouseDelta:          mgl64.Vec2{},
		mouseWheel:          0.0,
		wantCaptureKeyboard: false,
		wantCaptureMouse:    false,
	}
}

// IsKeyPressed returns true if the specified key is currently pressed
func (is *InputState) IsKeyPressed(key Key) bool {
	return is.keyStates[key]
}

// IsKeyReleased returns true if the specified key was just released (not currently pressed)
func (is *InputState) IsKeyReleased(key Key) bool {
	return !is.keyStates[key]
}

// IsMouseButtonPressed returns true if the specified mouse button is currently pressed
func (is *InputState) IsMouseButtonPressed(button MouseButton) bool {
	return is.mouseButtonStates[button]
}

// IsMouseButtonReleased returns true if the specified mouse button was just released
func (is *InputState) IsMouseButtonReleased(button MouseButton) bool {
	return !is.mouseButtonStates[button]
}

// MousePosition returns the current mouse position
func (is *InputState) MousePosition() mgl64.Vec2 {
	return is.mousePosition
}

// MouseDelta returns the mouse movement delta since last frame
func (is *InputState) MouseDelta() mgl64.Vec2 {
	return is.mouseDelta
}

// MouseWheel returns the mouse wheel delta since last frame
func (is *InputState) MouseWheel() float64 {
	return is.mouseWheel
}

// WantCaptureKeyboard returns true if the UI wants to capture keyboard input
func (is *InputState) WantCaptureKeyboard() bool {
	return is.wantCaptureKeyboard
}

// WantCaptureMouse returns true if the UI wants to capture mouse input
func (is *InputState) WantCaptureMouse() bool {
	return is.wantCaptureMouse
}

// InputProvider defines the interface for input providers (ImGui, GLFW, etc.)
type InputProvider interface {
	// GetInputState returns the current input state
	GetInputState() *InputState

	// Update updates the input state with current input data
	Update() error

	// IsValid returns true if the input provider is in a valid state
	IsValid() bool
}
