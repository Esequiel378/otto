package otto

import (
	"otto/system/input"
	"otto/system/physics"
	"otto/system/renderer"

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
		// Send initialization to physics system
		c.Send(e.PhysicsPID(), physics.EventRigidBodyRegister{
			PID: c.PID(),
			EntityRigidBody: physics.EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
		// Also send initialization to renderer
		c.Send(e.RendererPID(), renderer.EventEntityRegister{
			PID: c.PID(),
			EntityRigidBody: physics.EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
	case physics.EventRigidBodyTransform:
		e.Transform(c, msg)
		c.Send(e.RendererPID(), renderer.EventEntityRenderUpdate{
			PID: c.PID(),
			EntityRigidBody: physics.EntityRigidBody{
				Position: e.position,
				Velocity: e.velocity,
				Scale:    e.scale,
				Rotation: e.rotation,
			},
		})
	}
}

func (e *Entity) Transform(c *actor.Context, msg physics.EventRigidBodyTransform) {
	e.position = msg.Position
}

func (e *Entity) RegisterInputs(c *actor.Context, contexts ...input.Context) {
	c.Send(e.InputPID(), input.EventRegisterInputs{
		Contexts: contexts,
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
