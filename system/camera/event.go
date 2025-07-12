package camera

import "otto/system"

type CameraUpdate struct {
	Camera system.Camera
}

// RequestCamera is sent to get the current camera state
type RequestCamera struct{}

// ResponseCamera contains the current camera state
type ResponseCamera struct {
	Camera system.Camera
}
