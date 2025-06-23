package main

import (
	"log"
	"otto/receiver"
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

	rendererPID := e.Spawn(receiver.NewRender(), "renderer")
	physicsPID := e.Spawn(receiver.NewPhysics(), "physics")

	e.Spawn(receiver.NewPlayer(physicsPID, rendererPID), "player")

	window, err := NewSDLBackendWithOpenGL(1200, 900, "Hello from cimgui-go")
	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}

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
		resp := e.Request(rendererPID, receiver.RequestEntities{}, 10*time.Millisecond)

		res, err := resp.Result()
		if err != nil {
			log.Fatalf("failed to request entities: %v", err)
		}

		entities, ok := res.(receiver.EntitiesResponse)
		if !ok {
			log.Fatalf("failed to cast entities response: %v", res)
		}

		for pid, entity := range entities.Entities {
			_, _ = pid, entity
		}
	})
}
