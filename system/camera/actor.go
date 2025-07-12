package camera

import (
	"otto/system"
	"otto/system/physics"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// Camera manages camera state and handles camera input events
type Camera struct {
	physicsPID  *actor.PID
	rendererPID *actor.PID
	camera      system.Camera
}

var _ actor.Receiver = (*Camera)(nil)

// NewCamera creates a new camera confinguration that does not need to be used as an actor
func NewCamera(physicsPID, rendererPID *actor.PID) *Camera {
	return &Camera{
		physicsPID:  physicsPID,
		rendererPID: rendererPID,
		camera: system.Camera{
			Position: mgl64.Vec3{0, 0, -2},
			Rotation: mgl64.Vec2{0, 0},
			Zoom:     1.0,
		},
	}
}

// New creates a new camera configuration that must be used as an actor
func New(physicsPID, rendererPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return NewCamera(physicsPID, rendererPID)
	}
}

func (c *Camera) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case physics.EventRigidBodyTransform:
		c.camera.Position = msg.Position
		c.camera.Rotation = mgl64.Vec2{msg.Rotation[0], msg.Rotation[1]}
		// TODO: Send to renderer when it can actually render the game
	case RequestCamera:
		ctx.Respond(ResponseCamera{
			Camera: c.camera,
		})
	}
}
