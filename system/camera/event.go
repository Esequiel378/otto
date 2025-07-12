package camera

type CameraUpdate struct {
	Camera Camera
}

// RequestCamera is sent to get the current camera state
type RequestCamera struct{}

// ResponseCamera contains the current camera state
type ResponseCamera struct {
	Camera Camera
}
