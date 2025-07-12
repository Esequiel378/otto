package input

// EventInput represents a broadcasted input event
type EventInput struct {
	Context Context
}

// Context defines an interface for any input context that can be processed
type Context interface {
	// Process handles the input and returns true if any input was detected
	Process() bool
	// GetType returns the type identifier for this input context
	GetType() string
}
