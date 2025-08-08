package player

import (
	"otto"
	"otto/cmd/playground/camera"
	"otto/cmd/playground/floor"
	"otto/system/input"
	"otto/system/physics"
	"otto/system/renderer"
	"otto/util"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Player struct {
	*otto.Entity
	rendererPID *actor.PID
	physicsPID  *actor.PID
	inputPID    *actor.PID

	cameraPID *actor.PID
	floorPID  *actor.PID
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
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		p.cameraPID = ctx.SpawnChild(camera.New(p.physicsPID, p.rendererPID, p.inputPID), "camera")
		p.floorPID = ctx.SpawnChild(floor.New(p.rendererPID), "floor")

		input.RegisterInputs(
			ctx,
			p.inputPID,
			&InputPlayerMovement{PID: ctx.PID()},
		)

		// Send initialization to physics system
		ctx.Send(p.physicsPID, physics.EventRigidBodyRegister{
			PID:             ctx.PID(),
			EntityRigidBody: p.ToRigidBody(),
		})
		// Send initialization to renderer system
		ctx.Send(p.rendererPID, renderer.EventEntityRegister{
			PID:             ctx.PID(),
			EntityRigidBody: p.ToRigidBody(),
		})
	case input.EventInput:
		p.HandleInput(ctx, msg)
	case physics.EventRotationUpdate:
		p.Entity.Rotation = msg.Rotation
	case physics.EventRigidBodyTransform:
		ctx.Send(p.cameraPID, physics.EventPositionUpdate{
			PID:      ctx.PID(),
			Position: msg.Position,
		})
		ctx.Send(p.floorPID, physics.EventPositionUpdate{
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

		// Create horizontal-only versions of front and right vectors (Y = 0)
		frontHorizontal := mgl64.Vec3{front.X(), 0, front.Z()}.Normalize()
		rightHorizontal := mgl64.Vec3{right.X(), 0, right.Z()}.Normalize()

		// Transform horizontal movement (X and Z) using camera-relative vectors
		horizontalVelocity := rightHorizontal.Mul(input.Velocity.X()).
			Add(frontHorizontal.Mul(input.Velocity.Z()))

		// Keep Y movement separate and direct (controlled by SPACE/SHIFT)
		// Invert Y movement to fix the direction
		verticalVelocity := mgl64.Vec3{0, -input.Velocity.Y(), 0}

		// Combine horizontal and vertical movement
		velocity := horizontalVelocity.Add(verticalVelocity)

		ctx.Send(p.physicsPID, physics.EventRigidBodyUpdate{
			PID:             ctx.PID(),
			Velocity:        velocity,
			AngularVelocity: mgl64.Vec3{}, // No rotation for movement input
		})
	}
}
