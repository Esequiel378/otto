package floor

import (
	"otto"
	"otto/system/physics"
	"otto/system/renderer"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Floor struct {
	Entity *otto.Entity

	rendererPID *actor.PID
}

var _ actor.Receiver = (*Floor)(nil)

func New(renderPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		entity := otto.NewEntity(nil, nil, nil)
		entity.ModelName = "plane"
		entity.Scale = mgl64.Vec3{100, 1, 100}

		return &Floor{
			Entity:      entity,
			rendererPID: renderPID,
		}
	}
}

func (f *Floor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		ctx.Send(f.rendererPID, renderer.EventEntityRegister{
			PID:             ctx.PID(),
			EntityRigidBody: f.Entity.ToRigidBody(),
		})
	case physics.EventPositionUpdate:
		// Set Y to 0 to make the floor appear at the bottom of the world
		msg.Position[1] = 0
		f.Entity.Position = msg.Position
		ctx.Send(f.rendererPID, renderer.EventEntityRenderUpdate{
			PID:             ctx.PID(),
			EntityRigidBody: f.Entity.ToRigidBody(),
		})
	}
}
