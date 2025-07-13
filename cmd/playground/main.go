package main

import (
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

	tickRate := 64 // Increased from 64 for smoother input processing
	tickInterval := time.Second / time.Duration(tickRate)
	latestTick := time.Now()

	// FPS tracking variables
	var lastFPSUpdate time.Time
	var currentFPS float64
	var frameTimes []float64
	maxFrameTimes := 60 // Keep last 60 frame times for averaging

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

		// Render FPS overlay
		imgui.Begin("Performance")
		imgui.Text(fmt.Sprintf("FPS: %.1f", currentFPS))
		imgui.Text(fmt.Sprintf("Frame Time: %.3f ms", deltaTime*1000))
		imgui.Text(fmt.Sprintf("Tick Rate: %d Hz", tickRate))
		imgui.End()

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

		// Render entities using OpenGL
		for _, entity := range response.Entities {
			otto.RenderEntity(shaderManager, modelManager, &entity, &response.Camera)
		}
	})
}
