package player

import (
	"otto/manager"
	"otto/receiver"

	"github.com/anthdm/hollywood/actor"
)

type Player struct {
	*receiver.Entity
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: receiver.NewEntity(physicsPID, rendererPID, inputPID),
		}
	}
}

func (p *Player) Receive(c *actor.Context) {
	defer p.Entity.Receive(c)

	switch msg := c.Message().(type) {
	case actor.Initialized:
		p.RegisterInputContext(c, &InputPlayerMovement{})
	case manager.InputEvent:
		switch ctx := msg.Context.(type) {
		case *InputPlayerMovement:
			c.Send(p.PhysicsPID(), receiver.EventEntityUpdate{
				PID:      c.PID(),
				Velocity: ctx.Velocity,
			})
		}
	}
}
