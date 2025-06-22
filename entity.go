package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Entity struct {
	position mgl64.Vec3
	velocity mgl64.Vec3
	scale    mgl64.Vec3
	rotation mgl64.Vec3

	physicsPID  *actor.PID
	rendererPID *actor.PID
}

var _ actor.Receiver = (*Entity)(nil)

func NewEntity(physicsPID, rendererPID *actor.PID, render func()) actor.Producer {
	return func() actor.Receiver {
		return &Entity{
			physicsPID:  physicsPID,
			rendererPID: rendererPID,
		}
	}
}

// Receive implements actor.Receiver.
func (e *Entity) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		c.Engine().BroadcastEvent(EventEntityInitialized{
			PID: c.PID(),
			EntityRigidBody: EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
	case EventEntityTransform:
		e.Transform(c, msg)
		c.Send(e.rendererPID, EventEntityRenderUpdate{
			PID: c.PID(),
			EntityRigidBody: EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
	case Tick:
		c.Send(e.physicsPID, EventEntityUpdate{
			PID:      c.PID(),
			Velocity: e.velocity,
		})
	}
}

func (e *Entity) Transform(c *actor.Context, msg EventEntityTransform) {
	e.position = msg.Position
}
