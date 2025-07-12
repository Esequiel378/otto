package otto

import (
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
}

var _ actor.Receiver = (*Entity)(nil)

func NewEntity(physicsPID, rendererPID *actor.PID) *Entity {
	return &Entity{
		scale:       mgl64.Vec3{1, 1, 1},
		physicsPID:  physicsPID,
		rendererPID: rendererPID,
	}
}

// Receive implements actor.Receiver.
func (e *Entity) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		if e.physicsPID != nil {
			// Send initialization to physics system
			c.Send(e.physicsPID, physics.EventRigidBodyRegister{
				PID: c.PID(),
				EntityRigidBody: physics.EntityRigidBody{
					Position:        e.position,
					Velocity:        e.velocity,
					Scale:           e.scale,
					Rotation:        e.rotation,
					AngularVelocity: mgl64.Vec3{},
				},
			})
		}
		if e.rendererPID != nil {
			// Also send initialization to renderer
			c.Send(e.rendererPID, renderer.EventEntityRegister{
				PID: c.PID(),
				EntityRigidBody: physics.EntityRigidBody{
					Position:        e.position,
					Velocity:        e.velocity,
					Scale:           e.scale,
					Rotation:        e.rotation,
					AngularVelocity: mgl64.Vec3{},
				},
			})
		}
	case physics.EventRigidBodyTransform:
		e.Transform(c, msg)
		if e.rendererPID != nil {
			c.Send(e.rendererPID, renderer.EventEntityRenderUpdate{
				PID: c.PID(),
				EntityRigidBody: physics.EntityRigidBody{
					Position:        e.position,
					Velocity:        e.velocity,
					Scale:           e.scale,
					Rotation:        e.rotation,
					AngularVelocity: mgl64.Vec3{},
				},
			})
		}
	}
}

func (e *Entity) Transform(c *actor.Context, msg physics.EventRigidBodyTransform) {
	e.position = msg.Position
	e.rotation = msg.Rotation
}
