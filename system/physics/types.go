package physics

import "github.com/go-gl/mathgl/mgl64"

type EntityRigidBody struct {
	Position        mgl64.Vec3
	Velocity        mgl64.Vec3
	Scale           mgl64.Vec3
	Rotation        mgl64.Vec3
	AngularVelocity mgl64.Vec3
	ModelName       string
}
