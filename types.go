package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type Tick struct {
	deltaTime float64
}

type Collider struct {
	entity   *actor.PID
	position mgl64.Vec3
}

type RequestEntities struct{}
type EntitiesResponse struct {
	Entities []EntityRigidBody
}

type EntityRigidBody struct {
	Position mgl64.Vec3
	Velocity mgl64.Vec3
	Scale    mgl64.Vec3
	Rotation mgl64.Vec3
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
