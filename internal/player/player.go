package player

import (
	"otto/manager"
	"otto/system"

	"github.com/anthdm/hollywood/actor"
)

type Player struct {
	*system.Entity
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: system.NewEntity(physicsPID, rendererPID, inputPID),
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
			c.Send(p.PhysicsPID(), system.EventEntityUpdate{
				PID:      c.PID(),
				Velocity: ctx.Velocity,
			})
		}
	}
}
