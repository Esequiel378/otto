package camera

import (
	"otto/system"
	"otto/system/input"

	"math"

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

func New(inputPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &CameraActor{
			camera: Camera{
				Position: mgl64.Vec3{0, 0, -2}, // Start slightly back from origin
				Rotation: mgl64.Vec2{0, 0},     // No rotation (looking forward)
				Zoom:     1.0,                  // Default zoom
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
		c.Send(ca.inputPID, input.EventRegisterInputs{
			Contexts: []input.Context{&InputCameraControl{}},
		})
	case input.EventInput:
		switch msg.Context.(type) {
		case *InputCameraControl:
			ca.handleCameraInput(msg.Context.(*InputCameraControl))
		}
	case RequestCamera:
		// Respond with current camera state
		c.Respond(ResponseCamera{
			Camera: ca.camera,
		})
	case system.Tick:
		// Broadcast camera update on every tick
		c.Engine().BroadcastEvent(CameraUpdate{
			Camera: ca.camera,
		})
	}
}

func (ca *CameraActor) handleCameraInput(input *InputCameraControl) {
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

// GetFrontVector returns the camera's front (forward) direction vector
func (c *Camera) GetFrontVector() mgl64.Vec3 {
	pitch := c.Rotation[0]
	yaw := c.Rotation[1]

	cosPitch := math.Cos(pitch)
	sinPitch := math.Sin(pitch)
	cosYaw := math.Cos(yaw)
	sinYaw := math.Sin(yaw)

	return mgl64.Vec3{
		cosPitch * sinYaw,
		sinPitch,
		cosPitch * cosYaw,
	}
}

// GetRightVector returns the camera's right direction vector
func (c *Camera) GetRightVector() mgl64.Vec3 {
	yaw := c.Rotation[1]
	cosYaw := math.Cos(yaw)
	sinYaw := math.Sin(yaw)

	return mgl64.Vec3{
		cosYaw,
		0,
		-sinYaw,
	}
}

// GetUpVector returns the camera's up direction vector
func (c *Camera) GetUpVector() mgl64.Vec3 {
	right := c.GetRightVector()
	front := c.GetFrontVector()
	return right.Cross(front)
}
