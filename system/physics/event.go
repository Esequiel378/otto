package physics

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

type EventRigidBodyRegister struct {
	PID             *actor.PID
	EntityRigidBody EntityRigidBody
}
type EventRigidBodyUpdate struct {
	PID             *actor.PID
	Velocity        mgl64.Vec3
	AngularVelocity mgl64.Vec3
}
type EventRigidBodyTransform struct {
	PID      *actor.PID
	Position mgl64.Vec3
	Rotation mgl64.Vec3
}
type EventPositionUpdate struct {
	PID      *actor.PID
	Position mgl64.Vec3
}
type EventRotationUpdate struct {
	PID      *actor.PID
	Rotation mgl64.Vec3
}
