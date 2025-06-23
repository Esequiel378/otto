package receiver

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

func (p *Player) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case InputPlayerMovement:
		c.Send(p.physicsPID, EventEntityUpdate{
			PID:      c.PID(),
			Velocity: msg.Velocity,
		})
	}

	p.Entity.Receive(c)
}
