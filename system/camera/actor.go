package camera

import (
	"otto"
	"otto/system"
	"otto/system/input"
	"otto/system/physics"
	"otto/system/renderer"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// Camera manages camera state and handles camera input events
type Camera struct {
	*otto.Entity
	physicsPID  *actor.PID
	rendererPID *actor.PID
	inputPID    *actor.PID
	camera      system.Camera
}

var _ actor.Receiver = (*Camera)(nil)

// NewCamera creates a new camera confinguration that does not need to be used as an actor
func NewCamera(physicsPID, rendererPID, inputPID *actor.PID) *Camera {
	return &Camera{
		physicsPID:  physicsPID,
		rendererPID: rendererPID,
		inputPID:    inputPID,
		Entity:      otto.NewEntity(nil, rendererPID, inputPID),
		camera: system.Camera{
			Position: mgl64.Vec3{0, 0, -2},
			Rotation: mgl64.Vec2{0, 0},
			Zoom:     1.0,
		},
	}
}

// New creates a new camera configuration that must be used as an actor
func New(physicsPID, rendererPID, inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return NewCamera(physicsPID, rendererPID, inputPID)
	}
}

func (c *Camera) Receive(ctx *actor.Context) {
	defer c.Entity.Receive(ctx)

	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		input.RegisterInputs(
			ctx,
			c.inputPID,
			&InputCamera{PID: ctx.PID()},
		)
		ctx.Send(c.rendererPID, renderer.EventUpdateCamera{
			Camera: c.camera,
		})
	case input.EventInput:
		c.HandleInput(ctx, msg)
	case physics.EventRigidBodyTransform:
		c.camera.Position = msg.Position
		c.camera.Rotation = mgl64.Vec2{msg.Rotation[0], msg.Rotation[1]}
		ctx.Send(c.rendererPID, renderer.EventUpdateCamera{
			Camera: c.camera,
		})
	}
}

func (c *Camera) HandleInput(ctx *actor.Context, event input.EventInput) {
	switch input := event.Context.(type) {
	case *InputCamera:
		// Also send to physics for entity rotation (if needed)
		ctx.Send(c.physicsPID, physics.EventRigidBodyUpdate{
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
