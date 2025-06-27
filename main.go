package main

import (
	"log"
	"otto/manager"
	"otto/receiver"
	"runtime"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatalf("failed to create actor engine: %v", err)
	}

	rendererPID := e.Spawn(receiver.NewRender(), "renderer")
	physicsPID := e.Spawn(receiver.NewPhysics(), "physics")

	playerPID := e.Spawn(receiver.NewPlayer(physicsPID, rendererPID), "player")

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
			e.BroadcastEvent(receiver.Tick{DeltaTime: deltaTime})
		}
	}()

	window.Run(func(deltaTime float64) {
		// Handle input
		input := receiver.InputPlayerMovement{}
		input.Handle()
		if input.Velocity != (mgl64.Vec3{}) {
			e.Send(playerPID, input)
		}

		// Request entities from renderer
		resp := e.Request(rendererPID, receiver.RequestEntities{}, 10*time.Millisecond)

		res, err := resp.Result()
		if err != nil {
			log.Printf("failed to request entities: %v", err)
			return
		}

		entities, ok := res.(receiver.EntitiesResponse)
		if !ok {
			log.Printf("failed to cast entities response: %v", res)
			return
		}

		// Render entities using OpenGL
		for _, entity := range entities.Entities {
			RenderEntity(shaderManager, modelManager, &entity)
		}
	})
}
