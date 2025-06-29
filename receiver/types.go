package receiver

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Tick struct {
	DeltaTime float64
}

type RequestEntities struct{}
type EntitiesResponse struct {
	Entities []EntityRigidBody
}

type EntityRigidBody struct {
	Position  mgl64.Vec3
	Velocity  mgl64.Vec3
	Scale     mgl64.Vec3
	Rotation  mgl64.Vec3
	ModelName string
}
type EventEntityInitialized struct {
	PID             *actor.PID
	EntityRigidBody EntityRigidBody
}
type EventEntityUpdate struct {
	PID      *actor.PID
	Velocity mgl64.Vec3
}
type EventEntityRenderUpdate struct {
	PID             *actor.PID
	EntityRigidBody EntityRigidBody
}
type EventEntityTransform struct {
	PID      *actor.PID
	Position mgl64.Vec3
}

type SetCameraPID struct {
	PID *actor.PID
}

type CameraUpdate struct {
	Camera Camera
}
