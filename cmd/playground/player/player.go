package player

import (
	"otto"
	"otto/system/input"
	"otto/system/physics"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Player struct {
	*otto.Entity
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: otto.NewEntity(physicsPID, rendererPID),
		}
	}
}

func (p *Player) Receive(c *actor.Context) {
	defer p.Entity.Receive(c)

	switch msg := c.Message().(type) {
	case actor.Initialized:
		input.RegisterInputs(c, &InputPlayerMovement{}, &InputPlayerCamera{})
	case input.EventInput:
		switch ctx := msg.Context.(type) {
		case *InputPlayerMovement:
			c.Send(p.PhysicsPID(), physics.EventRigidBodyUpdate{
				PID:             c.PID(),
				Velocity:        ctx.Velocity,
				AngularVelocity: mgl64.Vec3{}, // No rotation for movement input
			})
		case *InputPlayerCamera:
			// Convert 2D rotation (pitch, yaw) to 3D angular velocity
			angularVelocity := mgl64.Vec3{
				ctx.Rotation[0], // Pitch -> X rotation
				ctx.Rotation[1], // Yaw -> Y rotation
				0,               // No roll
			}
			c.Send(p.PhysicsPID(), physics.EventRigidBodyUpdate{
				PID:             c.PID(),
				Velocity:        mgl64.Vec3{}, // No movement for camera input
				AngularVelocity: angularVelocity,
			})
		}
	}
}
