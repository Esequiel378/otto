package player

import (
	"otto"
	"otto/system/camera"
	"otto/system/input"
	"otto/system/physics"
	"otto/util"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Player struct {
	*otto.Entity
	camera *camera.Camera
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			Entity: otto.NewEntity(physicsPID, nil, inputPID),
		}
	}
}

func (p *Player) Receive(ctx *actor.Context) {
	defer p.Entity.Receive(ctx)
	defer p.camera.Receive(ctx)

	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		p.camera = camera.NewCamera(nil, p.RendererPID())
		input.RegisterInputs(
			ctx,
			p.InputPID(),
			&InputPlayerMovement{PID: ctx.PID()},
			&InputPlayerCamera{PID: ctx.PID()},
		)
	case input.EventInput:
		p.HandleInput(ctx, msg)
	}
}

func (p *Player) HandleInput(ctx *actor.Context, event input.EventInput) {
	switch input := event.Context.(type) {
	case *InputPlayerMovement:
		// Use camera vectors to transform velocity into world space
		front := util.Vec3FrontVector(p.Entity.Rotation)
		right := util.Vec3RightVector(p.Entity.Rotation)
		up := util.Vec3UpVector(p.Entity.Rotation)

		// Transform velocity from camera-relative to world coordinates
		velocity := right.Mul(input.Velocity.X()).
			Add(up.Mul(input.Velocity.Y())).
			Add(front.Mul(input.Velocity.Z()))

		ctx.Send(p.PhysicsPID(), physics.EventRigidBodyUpdate{
			PID:             ctx.PID(),
			Velocity:        velocity,
			AngularVelocity: mgl64.Vec3{}, // No rotation for movement input
		})
	case *InputPlayerCamera:
		ctx.Send(p.PhysicsPID(), physics.EventRigidBodyUpdate{
			PID: ctx.PID(),
			// Convert 2D rotation (pitch, yaw) to 3D angular velocity
			AngularVelocity: mgl64.Vec3{
				input.Rotation[0],
				input.Rotation[1],
				0,
			},
		})
	}
}
