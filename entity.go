package otto

import (
	"otto/system/physics"
	"otto/system/renderer"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Entity struct {
	Position   mgl64.Vec3
	Velocity   mgl64.Vec3
	Scale      mgl64.Vec3
	Rotation   mgl64.Vec3
	ModelName  string
	EntityType string // "player", "cube", "floor", etc.

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
				PID:             ctx.PID(),
				EntityRigidBody: e.ToRigidBody(),
			})
		}
		if e.rendererPID != nil {
			// Also send initialization to renderer
			ctx.Send(e.rendererPID, renderer.EventEntityRegister{
				PID:             ctx.PID(),
				EntityRigidBody: e.ToRigidBody(),
			})
		}
	case physics.EventRigidBodyTransform:
		e.Transform(ctx, msg)
		if e.rendererPID != nil {
			ctx.Send(e.rendererPID, renderer.EventEntityRenderUpdate{
				PID: ctx.PID(), EntityRigidBody: e.ToRigidBody(),
			})
		}
	}
}

func (e *Entity) ToRigidBody() physics.EntityRigidBody {
	return physics.EntityRigidBody{
		Position:   e.Position,
		Velocity:   e.Velocity,
		Scale:      e.Scale,
		Rotation:   e.Rotation,
		ModelName:  e.ModelName,
		EntityType: e.EntityType,
	}
}

func (e *Entity) Transform(ctx *actor.Context, msg physics.EventRigidBodyTransform) {
	e.Position = msg.Position
	e.Rotation = msg.Rotation
}
