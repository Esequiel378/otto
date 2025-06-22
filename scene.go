package main

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
)

type Scene struct {
}

var _ actor.Receiver = (*Scene)(nil)

func NewScene() actor.Producer {
	return func() actor.Receiver {
		return &Scene{}
	}
}

// Receive implements actor.Receiver.
func (s *Scene) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Started:
		fmt.Println("Started")
	case Tick:
		fmt.Println("Update")
	}
}
