package player

import (
	"otto"
	"otto/system/input"
	"otto/system/physics"

	"github.com/anthdm/hollywood/actor"
)

type Player struct {
	*otto.Entity
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: otto.NewEntity(physicsPID, rendererPID, inputPID),
		}
	}
}

func (p *Player) Receive(c *actor.Context) {
	defer p.Entity.Receive(c)

	switch msg := c.Message().(type) {
	case actor.Initialized:
		p.RegisterInputs(c, &InputPlayerMovement{}, &InputPlayerCamera{})
	case input.EventInput:
		switch ctx := msg.Context.(type) {
		case *InputPlayerMovement:
			c.Send(p.PhysicsPID(), physics.EventRigidBodyUpdate{
				PID:      c.PID(),
				Velocity: ctx.Velocity,
			})
		}
	}
}
