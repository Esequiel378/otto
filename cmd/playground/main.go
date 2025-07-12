package main

import (
	"log"
	"otto"
	"otto/cmd/playground/cube"
	"otto/cmd/playground/player"
	"otto/manager"
	"otto/system"
	"otto/system/camera"
	"otto/system/input"
	"otto/system/physics"
	"otto/system/renderer"
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

	inputPID := e.Spawn(input.New(), "input")
	rendererPID := e.Spawn(renderer.New(), "renderer")
	physicsPID := e.Spawn(physics.New(), "physics")

	playerPID := e.Spawn(player.NewPlayer(physicsPID, rendererPID, inputPID), "player")

	e.Spawn(cube.NewCube(physicsPID, rendererPID, inputPID), "test_cube")

	window, err := otto.NewSDLBackendWithOpenGL(1200, 900, "Hello from cimgui-go")
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
		cameraResp := e.Request(playerPID, camera.RequestCamera{}, 10*time.Millisecond)
		cameraRes, err := cameraResp.Result()
		if err != nil {
			log.Printf("failed to request camera: %v", err)
			return
		}
		camera, ok := cameraRes.(camera.ResponseCamera)
		if !ok {
			log.Printf("failed to cast camera response: %v", cameraRes)
			return
		}

		// Request entities from renderer
		resp := e.Request(rendererPID, renderer.RequestEntities{}, 10*time.Millisecond)

		res, err := resp.Result()
		if err != nil {
			log.Printf("failed to request entities: %v", err)
			return
		}

		entities, ok := res.(renderer.EntitiesResponse)
		if !ok {
			log.Printf("failed to cast entities response: %v", res)
			return
		}

		// Render entities using OpenGL
		for _, entity := range entities.Entities {
			log.Printf("rendering entity: %v", entity)
			otto.RenderEntity(shaderManager, modelManager, &entity, &camera.Camera)
		}
	})
}
