package receiver

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Player struct {
	Entity
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: Entity{
				position:    mgl64.Vec3{0, 0, 0}, // Start at origin
				velocity:    mgl64.Vec3{0, 0, 0},
				scale:       mgl64.Vec3{1, 1, 1}, // Make it visible size
				rotation:    mgl64.Vec3{0, 0, 0},
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
