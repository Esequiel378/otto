package main

import (
	"log"
	"otto/internal/cube"
	"otto/internal/player"
	"otto/manager"
	"otto/system"
	"runtime"
	"time"

	"github.com/anthdm/hollywood/actor"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatalf("failed to create actor engine: %v", err)
	}

	inputPID := e.Spawn(system.NewInputActor(), "input")
	cameraPID := e.Spawn(system.NewCamera(inputPID), "camera")
	rendererPID := e.Spawn(system.NewRender(), "renderer")
	physicsPID := e.Spawn(system.NewPhysics(cameraPID), "physics")

	e.Spawn(player.NewPlayer(physicsPID, rendererPID, inputPID), "player")

	e.Spawn(cube.NewCube(physicsPID, rendererPID, inputPID), "test_cube")

	window, err := NewSDLBackendWithOpenGL(1200, 900, "Hello from cimgui-go")
	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}

	// Initialize shader manager after OpenGL context is created
	shaderManager := manager.NewShaderManager()
	if err := shaderManager.Init("./assets/shaders"); err != nil {
		log.Fatalf("failed to initialize shader manager: %v", err)
	}
	defer shaderManager.Cleanup()

	// Initialize model manager
	modelManager := manager.NewModelManager()
	if err := modelManager.Init("./assets/models", "./assets/textures"); err != nil {
		log.Printf("Warning: failed to initialize model manager: %v", err)
		// Continue without models for now
	}
	defer modelManager.Cleanup()

	tickRate := 64
	tickInterval := time.Second / time.Duration(tickRate)
	latestTick := time.Now()

	// TODO: Add cancelation context
	go func() {
		ticker := time.NewTicker(tickInterval)
		defer ticker.Stop()

		// The broadcast overhead makes this a bit less accurate than the tick rate, but it's good enough for now.
		for range ticker.C {
			now := time.Now()
			deltaTime := now.Sub(latestTick).Seconds()
			latestTick = now
			e.BroadcastEvent(system.Tick{DeltaTime: deltaTime})
		}
	}()

	window.Run(func(deltaTime float64) {
		// Request camera data
		cameraResp := e.Request(cameraPID, system.RequestCamera{}, 10*time.Millisecond)
		cameraRes, err := cameraResp.Result()
		if err != nil {
			log.Printf("failed to request camera: %v", err)
			return
		}
		camera, ok := cameraRes.(system.ResponseCamera)
		if !ok {
			log.Printf("failed to cast camera response: %v", cameraRes)
			return
		}

		// Request entities from renderer
		resp := e.Request(rendererPID, system.RequestEntities{}, 10*time.Millisecond)

		res, err := resp.Result()
		if err != nil {
			log.Printf("failed to request entities: %v", err)
			return
		}

		entities, ok := res.(system.EntitiesResponse)
		if !ok {
			log.Printf("failed to cast entities response: %v", res)
			return
		}

		// Render entities using OpenGL
		for _, entity := range entities.Entities {
			RenderEntity(shaderManager, modelManager, &entity, &camera.Camera)
		}
	})
}
