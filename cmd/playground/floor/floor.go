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
	physicsPID  *actor.PID
}

var _ actor.Receiver = (*Floor)(nil)

func New(renderPID, physicsPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		entity := otto.NewEntity(nil, nil, nil)
		entity.ModelName = "plane"
		entity.EntityType = "floor"
		// Make the floor world-wide with a much larger scale
		entity.Scale = mgl64.Vec3{1000, 1, 1000} // 1000x1000 world units
		// Set initial position at y=0
		entity.Position = mgl64.Vec3{0, 0, 0}

		return &Floor{
			Entity:      entity,
			rendererPID: renderPID,
			physicsPID:  physicsPID,
		}
	}
}

func (f *Floor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		// Register with renderer
		ctx.Send(f.rendererPID, renderer.EventEntityRegister{
			PID:             ctx.PID(),
			EntityRigidBody: f.Entity.ToRigidBody(),
		})
		// Register with physics system
		ctx.Send(f.physicsPID, physics.EventRigidBodyRegister{
			PID:             ctx.PID(),
			EntityRigidBody: f.Entity.ToRigidBody(),
		})
	case physics.EventPositionUpdate:
		// Keep the floor fixed at y=0 regardless of any position updates
		floorPosition := mgl64.Vec3{msg.Position.X(), 0, msg.Position.Z()}
		f.Entity.Position = floorPosition
		ctx.Send(f.rendererPID, renderer.EventEntityRenderUpdate{
			PID:             ctx.PID(),
			EntityRigidBody: f.Entity.ToRigidBody(),
		})
	}
}
