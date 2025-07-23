package otto

import (
	"fmt"
	"time"

	"github.com/AllenDang/cimgui-go/backend"
	"github.com/AllenDang/cimgui-go/backend/sdlbackend"
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Window interface {
	Run(func(deltaTime float64))
	Width() int
	Height() int
}

type SDLWindow struct {
	backend.Backend[sdlbackend.SDLWindowFlags]

	width  int
	height int

	lastTime time.Time
}

var _ Window = (*SDLWindow)(nil)

// NewSDLBackendWithOpenGL creates a new sdl backend with opengl
func NewSDLBackendWithOpenGL(width, height int, title string) (*SDLWindow, error) {
	currBackend, err := backend.CreateBackend(sdlbackend.NewSDLBackend())
	if err != nil {
		return nil, fmt.Errorf("failed to create backend: %w", err)
	}

	currBackend.SetWindowFlags(sdlbackend.SDLWindowFlagsResizable, 1)
	currBackend.CreateWindow(title, width, height)
	currBackend.SetSwapInterval(0)
	currBackend.SetBgColor(imgui.NewVec4(0.2, 0.3, 0.3, 1.0))

	flags := imgui.CurrentIO().ConfigFlags()
	flags |= imgui.ConfigFlagsViewportsEnable

	imgui.CurrentIO().SetConfigFlags(flags)
	imgui.CurrentIO().SetIniFilename("/.imgui.ini")

	// Configure ImGui to capture mouse movement for delta calculation
	// Note: Mouse delta should be captured automatically by the backend

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// Enable depth clamping to prevent z-fighting near the near plane
	gl.Enable(gl.DEPTH_CLAMP)

	// Enable face culling to only render front faces
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	// Set clear color and depth values
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.ClearDepth(1.0)

	window := &SDLWindow{
		Backend:  currBackend,
		lastTime: time.Now(),
		width:    width,
		height:   height,
	}

	return window, nil
}

func (w *SDLWindow) Run(f func(deltaTime float64)) {
	w.Backend.Run(func() {
		// Clear both color and depth buffers completely
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Ensure depth buffer is properly initialized
		gl.ClearDepth(1.0)

		currentTime := time.Now()
		deltaTime := currentTime.Sub(w.lastTime).Seconds()
		w.lastTime = currentTime

		f(deltaTime)
	})
}

func (w *SDLWindow) Width() int {
	return w.width
}

func (w *SDLWindow) Height() int {
	return w.height
}
