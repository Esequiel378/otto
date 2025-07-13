package util

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

// Vec2FrontVector returns the front (forward) direction vector for a 2D vector
func Vec2FrontVector(vec mgl64.Vec2) mgl64.Vec2 {
	front := Vec3FrontVector(mgl64.Vec3{vec[0], vec[1], 0})
	return mgl64.Vec2{front.X(), front.Z()}
}

// Vec2RightVector returns the right direction vector for a 2D vector
func Vec2RightVector(vec mgl64.Vec2) mgl64.Vec2 {
	right := Vec3RightVector(mgl64.Vec3{vec[0], vec[1], 0})
	return mgl64.Vec2{right.X(), right.Z()}
}

// Vec2UpVector returns the up direction vector for a 2D vector
func Vec2UpVector(vec mgl64.Vec2) mgl64.Vec2 {
	up := Vec3UpVector(mgl64.Vec3{vec[0], vec[1], 0})
	return mgl64.Vec2{up.X(), up.Z()}
}

// Vec3FrontVector returns the front (forward) direction vector
func Vec3FrontVector(vec mgl64.Vec3) mgl64.Vec3 {
	pitch := vec[0]
	yaw := vec[1]

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

// Vec3RightVector returns the right direction vector
func Vec3RightVector(vec mgl64.Vec3) mgl64.Vec3 {
	yaw := vec[1]
	cosYaw := math.Cos(yaw)
	sinYaw := math.Sin(yaw)

	return mgl64.Vec3{
		cosYaw,
		0,
		-sinYaw,
	}
}

// Vec3UpVector returns the up direction vector
func Vec3UpVector(vec mgl64.Vec3) mgl64.Vec3 {
	right := Vec3RightVector(vec)
	front := Vec3FrontVector(vec)
	return right.Cross(front)
}

// Convert mgl64.Vec3 to mgl32.Vec3
func Vec64ToVec32(v mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v.X()), float32(v.Y()), float32(v.Z())}
}
