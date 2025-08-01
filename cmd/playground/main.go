package main

import (
	"context"
	"fmt"
	"log"
	"otto"
	"otto/cmd/playground/cube"
	"otto/cmd/playground/player"
	"otto/manager"
	"otto/system"
	"otto/system/input"
	"otto/system/physics"
	"otto/system/renderer"
	"runtime"
	"time"

	"github.com/AllenDang/cimgui-go/imgui"
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

	inputPID := e.Spawn(input.New(), "input")
	rendererPID := e.Spawn(renderer.New(), "renderer")
	physicsPID := e.Spawn(physics.New(), "physics")

	e.Spawn(player.NewPlayer(physicsPID, rendererPID, inputPID), "player")

	// Spawn 100 cubes in a 10x10 grid with a gap of 1 cube between them
	// Each cube is 1 unit, so we place them 2 units apart (1 unit for cube + 1 unit gap)
	// Using batch rendering for better performance
	for i := 0; i < 50; i++ {
		for j := 0; j < 50; j++ {
			e.Spawn(
				cube.NewCubeWithPosition(
					physicsPID,
					rendererPID,
					inputPID,
					mgl64.Vec3{float64(i * 2), 0.5, float64(j * 2)}, // 2 units apart for 1 unit gap, Y=1 to be above grid
				),
				fmt.Sprintf("cube_%d_%d", i, j),
			)
		}
	}

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverTickRate := 256
	clientTickRate := 1_000

	go func(ctx context.Context) {
		tickInterval := time.Second / time.Duration(serverTickRate)
		latestTick := time.Now()

		ticker := time.NewTicker(tickInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				deltaTime := now.Sub(latestTick).Seconds()
				latestTick = now
				e.BroadcastEvent(system.ServerTick{DeltaTime: deltaTime})
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second / time.Duration(clientTickRate))
		defer ticker.Stop()
		var latestTick time.Time

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				now := time.Now()
				deltaTime := now.Sub(latestTick).Seconds()
				latestTick = now
				// TODO: Maybe we should broadcast this instead of sending it to the input PID?
				e.Send(inputPID, system.ClientTick{DeltaTime: deltaTime})
			}
		}
	}(ctx)

	// FPS tracking variables
	var lastFPSUpdate time.Time
	var currentFPS float64
	var frameTimes []float64
	maxFrameTimes := 60 // Keep last 60 frame times for averaging

	window.Run(func(deltaTime float64) {
		// Track frame time for FPS calculation
		if len(frameTimes) >= maxFrameTimes {
			frameTimes = frameTimes[1:]
		}
		frameTimes = append(frameTimes, deltaTime)

		// Update FPS every second
		now := time.Now()
		if now.Sub(lastFPSUpdate) >= time.Second {
			if len(frameTimes) > 0 {
				totalTime := 0.0
				for _, ft := range frameTimes {
					totalTime += ft
				}
				averageFrameTime := totalTime / float64(len(frameTimes))
				currentFPS = 1.0 / averageFrameTime
			}
			lastFPSUpdate = now
		}

		// Request entities from renderer
		resp := e.Request(rendererPID, renderer.RequestEntities{}, 10*time.Millisecond)

		res, err := resp.Result()
		if err != nil {
			log.Printf("failed to request entities: %v", err)
			return
		}

		response, ok := res.(renderer.EntitiesResponse)
		if !ok {
			log.Printf("failed to cast entities response: %v", res)
			return
		}

		// Render FPS overlay
		imgui.Begin("Performance")
		imgui.Text(fmt.Sprintf("FPS: %.1f", currentFPS))
		imgui.Text(fmt.Sprintf("Frame Time: %.3f ms", deltaTime*1000))
		imgui.Text(fmt.Sprintf("Tick Rate: %d Hz", serverTickRate))
		imgui.Text(fmt.Sprintf("Entities: %d", len(response.Entities)))
		imgui.End()

		// Render entities using OpenGL batch rendering for better performance
		var floor physics.EntityRigidBody
		entities := make([]*physics.EntityRigidBody, 0, len(response.Entities))
		for i := range response.Entities {
			if response.Entities[i].ModelName == "plane" {
				floor = response.Entities[i]
				continue
			}
			entities = append(entities, &response.Entities[i])
		}

		otto.RenderGridFloor(shaderManager, modelManager, &floor, &response.Camera)
		otto.RenderEntityBatch(shaderManager, modelManager, entities, &response.Camera)
	})
}
