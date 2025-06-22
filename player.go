package main

import (
	"github.com/anthdm/hollywood/actor"
)

type Player struct {
	Entity
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: Entity{
				physicsPID:  physicsPID,
				rendererPID: rendererPID,
			},
		}
	}
}
