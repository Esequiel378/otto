package physics

import (
	"otto/system/camera"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type EventRigidBodyRegister struct {
	PID             *actor.PID
	EntityRigidBody EntityRigidBody
}
type EventRigidBodyUpdate struct {
	PID      *actor.PID
	Velocity mgl64.Vec3
}
type EventRigidBodyTransform struct {
	PID      *actor.PID
	Position mgl64.Vec3
}

type EventCameraUpdate struct {
	Camera camera.Camera
}
