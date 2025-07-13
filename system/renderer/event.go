package renderer

import (
	"otto/system"
	"otto/system/physics"

	"github.com/anthdm/hollywood/actor"
)

type EventEntityRegister struct {
	PID             *actor.PID
	EntityRigidBody physics.EntityRigidBody
}

type EventEntityRenderUpdate struct {
	PID             *actor.PID
	EntityRigidBody physics.EntityRigidBody
}

type RequestEntities struct{}
type EntitiesResponse struct {
	Entities []physics.EntityRigidBody
	Camera   system.Camera
}

type EventUpdateCamera struct {
	Camera system.Camera
}
