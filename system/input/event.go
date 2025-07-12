package input

import "github.com/anthdm/hollywood/actor"

type EventRegisterInputs struct {
	Contexts []Context
}

// EventUnregisterInputs is sent by entities to unregister their input contexts
type EventUnregisterInputs struct {
	ContextTypes []string
}

// RegisterInputs is used to register input contexts with the input actor
func RegisterInputs(c *actor.Context, contexts ...Context) {
	c.Send(c.PID(), EventRegisterInputs{
		Contexts: contexts,
	})
}
