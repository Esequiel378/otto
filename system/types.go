package system

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
