package input

import "github.com/anthdm/hollywood/actor"

// EventInput represents a broadcasted input event
type EventInput struct {
	Context Context
}

// Context defines an interface for any input context that can be processed
type Context interface {
	// Process handles the input using the provided input state and returns true if any input was detected
	// The captureKeyboard and captureMouse parameters indicate if the UI wants to capture those input types
	Process(state *InputState, captureKeyboard, captureMouse bool) bool
	// GetPID returns the PID of the input context
	GetPID() *actor.PID
}
