package input

import (
	"log"

	"github.com/anthdm/hollywood/actor"
)

type EventRegisterInputs struct {
	Contexts []Context
}

// RegisterInputs is used to register input contexts with the input actor
func RegisterInputs(ctx *actor.Context, inputPID *actor.PID, contexts ...Context) {
	log.Printf("registering %s input contexts", contexts)
	ctx.Send(inputPID, EventRegisterInputs{
		Contexts: contexts,
	})
}
