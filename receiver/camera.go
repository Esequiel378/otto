package receiver

import (
	"otto/manager"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// Camera represents a 3D camera with position, rotation, and zoom
type Camera struct {
	Position mgl64.Vec3
	Rotation mgl64.Vec2 // Pitch, Yaw
	Zoom     float64
}

// CameraActor manages camera state and handles camera input events
type CameraActor struct {
	camera   Camera
	inputPID *actor.PID
}

var _ actor.Receiver = (*CameraActor)(nil)

func NewCamera(inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &CameraActor{
			camera: Camera{
				Position: mgl64.Vec3{0, 0, 5}, // Start at origin looking down -Z
				Rotation: mgl64.Vec2{0, 0},    // No rotation
				Zoom:     1.0,                 // Default zoom
			},
			inputPID: inputPID,
		}
	}
}

func (ca *CameraActor) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		// Subscribe to input events
		c.Engine().Subscribe(c.PID())
		// Register our input context with the input actor
		c.Send(ca.inputPID, RegisterInputContext{
			Context: &manager.InputCameraControl{},
		})
	case manager.InputEvent:
		switch msg.Context.(type) {
		case *manager.InputCameraControl:
			ca.handleCameraInput(msg.Context.(*manager.InputCameraControl))
		}
	case RequestCamera:
		// Respond with current camera state
		c.Respond(ResponseCamera{
			Camera: ca.camera,
		})
	}
}

func (ca *CameraActor) handleCameraInput(input *manager.InputCameraControl) {
	// Apply rotation
	if input.Rotation != (mgl64.Vec2{}) {
		rotationSpeed := 0.1
		ca.camera.Rotation = ca.camera.Rotation.Add(input.Rotation.Mul(rotationSpeed))

		// Clamp pitch to prevent gimbal lock
		if ca.camera.Rotation[0] > 1.5 {
			ca.camera.Rotation[0] = 1.5
		}
		if ca.camera.Rotation[0] < -1.5 {
			ca.camera.Rotation[0] = -1.5
		}
	}

	// Apply zoom
	if input.Zoom != 0 {
		zoomSpeed := 0.1
		ca.camera.Zoom += input.Zoom * zoomSpeed

		// Clamp zoom
		if ca.camera.Zoom < 0.1 {
			ca.camera.Zoom = 0.1
		}
		if ca.camera.Zoom > 10.0 {
			ca.camera.Zoom = 10.0
		}
	}
}

// RequestCamera is sent to get the current camera state
type RequestCamera struct{}

// ResponseCamera contains the current camera state
type ResponseCamera struct {
	Camera Camera
}
