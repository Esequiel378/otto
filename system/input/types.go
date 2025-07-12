package input

import "github.com/anthdm/hollywood/actor"

// EventInput represents a broadcasted input event
type EventInput struct {
	Context Context
}

// Context defines an interface for any input context that can be processed
type Context interface {
	// Process handles the input and returns true if any input was detected
	Process() bool
	// GetPID returns the PID of the input context
	GetPID() *actor.PID
}
