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
		Entity:      otto.NewEntity(physicsPID, nil, inputPID),
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
		return NewCamera(nil, rendererPID, inputPID)
	}
}

func (c *Camera) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		// Initialize InputCamera with default Euler angles (like reference code)
		inputCamera := NewInputCamera(ctx.PID())
		input.RegisterInputs(
			ctx,
			c.inputPID,
			inputCamera,
		)
		ctx.Send(c.rendererPID, renderer.EventUpdateCamera{
			Camera: c.camera,
		})
	case input.EventInput:
		c.HandleInput(ctx, msg)
	case physics.EventPositionUpdate:
		c.camera.Position = msg.Position
		ctx.Send(c.rendererPID, renderer.EventUpdateCamera{
			Camera: c.camera,
		})
	}
}

func (c *Camera) HandleInput(ctx *actor.Context, event input.EventInput) {
	switch input := event.Context.(type) {
	case *InputCamera:
		// Update camera rotation with Euler angles
		c.camera.Rotation = c.camera.Rotation.Add(input.Rotation())

		// TODO: This should be sent to the physics system too

		// Send updated camera to renderer
		ctx.Send(c.rendererPID, renderer.EventUpdateCamera{
			Camera: c.camera,
		})

		if ctx.Parent() != nil {
			ctx.Send(ctx.Parent(), physics.EventRotationUpdate{
				PID:      ctx.PID(),
				Rotation: c.camera.Rotation.Vec3(0),
			})
		}
	}
}
