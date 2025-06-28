package manager

// InputEvent represents a broadcasted input event
type InputEvent struct {
	Context InputContext
}

// InputContext defines an interface for any input context that can be processed
type InputContext interface {
	// Process handles the input and returns true if any input was detected
	Process() bool
	// GetType returns the type identifier for this input context
	GetType() string
}
