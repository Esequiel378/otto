package player

import (
	"otto"
	"otto/cmd/playground/camera"
	"otto/system/input"
	"otto/system/physics"
	"otto/util"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Player struct {
	*otto.Entity
	cameraPID   *actor.PID
	rendererPID *actor.PID
	physicsPID  *actor.PID
	inputPID    *actor.PID
}

var _ actor.Receiver = (*Player)(nil)

func NewPlayer(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Player{
			rendererPID: rendererPID,
			physicsPID:  physicsPID,
			inputPID:    inputPID,
			Entity:      otto.NewEntity(physicsPID, nil, inputPID),
		}
	}
}

func (p *Player) Receive(ctx *actor.Context) {
	defer p.Entity.Receive(ctx)

	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		p.cameraPID = ctx.SpawnChild(camera.New(p.physicsPID, p.rendererPID, p.inputPID), "camera")
		input.RegisterInputs(
			ctx,
			p.inputPID,
			&InputPlayerMovement{PID: ctx.PID()},
		)
	case input.EventInput:
		p.HandleInput(ctx, msg)
	case physics.EventRigidBodyTransform:
		ctx.Send(p.cameraPID, physics.EventPositionUpdate{
			PID:      ctx.PID(),
			Position: msg.Position,
		})
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

		ctx.Send(p.physicsPID, physics.EventRigidBodyUpdate{
			PID:             ctx.PID(),
			Velocity:        velocity,
			AngularVelocity: mgl64.Vec3{}, // No rotation for movement input
		})
	}
}
