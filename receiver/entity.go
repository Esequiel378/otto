package receiver

import (
	"otto/manager"

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
	inputPID    *actor.PID
}

var _ actor.Receiver = (*Entity)(nil)

func NewEntity(physicsPID, rendererPID, inputPID *actor.PID) *Entity {
	return &Entity{
		scale:       mgl64.Vec3{1, 1, 1},
		physicsPID:  physicsPID,
		rendererPID: rendererPID,
		inputPID:    inputPID,
	}
}

// Receive implements actor.Receiver.
func (e *Entity) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		c.Engine().Subscribe(c.PID())
		// Send initialization to physics system
		c.Send(e.PhysicsPID(), EventEntityInitialized{
			PID: c.PID(),
			EntityRigidBody: EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
		// Also send initialization to renderer
		c.Send(e.RendererPID(), EventEntityInitialized{
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
		c.Send(e.RendererPID(), EventEntityRenderUpdate{
			PID: c.PID(),
			EntityRigidBody: EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
	}
}

func (e *Entity) Transform(c *actor.Context, msg EventEntityTransform) {
	e.position = msg.Position
}

func (e *Entity) RegisterInputContext(c *actor.Context, inputContext manager.InputContext) {
	c.Send(e.InputPID(), RegisterInputContext{
		Context: inputContext,
	})
}

func (e *Entity) PhysicsPID() *actor.PID {
	return e.physicsPID
}

func (e *Entity) RendererPID() *actor.PID {
	return e.rendererPID
}

func (e *Entity) InputPID() *actor.PID {
	return e.inputPID
}
