package system

import "github.com/go-gl/mathgl/mgl64"

type TickInput struct{}
type Tick struct {
	DeltaTime float64
}

// Camera represents a 3D camera with position, rotation, and zoom
type Camera struct {
	Position mgl64.Vec3
	Rotation mgl64.Vec2 // Pitch, Yaw
	Zoom     float64
}
