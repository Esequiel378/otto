package otto

import (
	"otto/system/physics"
	"otto/system/renderer"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Entity struct {
	Position mgl64.Vec3
	Velocity mgl64.Vec3
	Scale    mgl64.Vec3
	Rotation mgl64.Vec3

	physicsPID  *actor.PID
	rendererPID *actor.PID
	inputPID    *actor.PID
}

var _ actor.Receiver = (*Entity)(nil)

func NewEntity(physicsPID, rendererPID, inputPID *actor.PID) *Entity {
	return &Entity{
		Scale:       mgl64.Vec3{1, 1, 1},
		physicsPID:  physicsPID,
		rendererPID: rendererPID,
		inputPID:    inputPID,
	}
}

// Receive implements actor.Receiver.
func (e *Entity) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		if e.physicsPID != nil {
			// Send initialization to physics system
			ctx.Send(e.physicsPID, physics.EventRigidBodyRegister{
				PID: ctx.PID(),
				EntityRigidBody: physics.EntityRigidBody{
					Position:        e.Position,
					Velocity:        e.Velocity,
					Scale:           e.Scale,
					Rotation:        e.Rotation,
					AngularVelocity: mgl64.Vec3{},
				},
			})
		}
		if e.rendererPID != nil {
			// Also send initialization to renderer
			ctx.Send(e.rendererPID, renderer.EventEntityRegister{
				PID: ctx.PID(),
				EntityRigidBody: physics.EntityRigidBody{
					Position:        e.Position,
					Velocity:        e.Velocity,
					Scale:           e.Scale,
					Rotation:        e.Rotation,
					AngularVelocity: mgl64.Vec3{},
				},
			})
		}
	case physics.EventRigidBodyTransform:
		e.Transform(ctx, msg)
		if e.rendererPID != nil {
			ctx.Send(e.rendererPID, renderer.EventEntityRenderUpdate{
				PID: ctx.PID(),
				EntityRigidBody: physics.EntityRigidBody{
					Position:        e.Position,
					Velocity:        e.Velocity,
					Scale:           e.Scale,
					Rotation:        e.Rotation,
					AngularVelocity: mgl64.Vec3{},
				},
			})
		}
	}
}

func (e *Entity) Transform(ctx *actor.Context, msg physics.EventRigidBodyTransform) {
	e.Position = msg.Position
	e.Rotation = msg.Rotation
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
