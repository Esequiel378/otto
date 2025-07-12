package input

type EventRegisterInputs struct {
	Contexts []Context
}

// EventUnregisterInputs is sent by entities to unregister their input contexts
type EventUnregisterInputs struct {
	ContextTypes []string
}
