package main

import (
	"fmt"
	"log"
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

	rendererPID := e.Spawn(NewRender(), "renderer")
	physicsPID := e.Spawn(NewPhysics(), "physics")

	e.Spawn(NewPlayer(physicsPID, rendererPID), "player")

	window, err := NewSDLBackendWithOpenGL(1200, 900, "Hello from cimgui-go")
	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}

	go func() {
		for {
			e.BroadcastEvent(Tick{deltaTime: 1})
			time.Sleep(1 * time.Second)
		}
	}()

	window.Run(func(deltaTime float64) {
		resp := e.Request(rendererPID, RequestEntities{}, 10*time.Millisecond)

		res, err := resp.Result()
		if err != nil {
			log.Fatalf("failed to request entities: %v", err)
		}

		entities, ok := res.(EntitiesResponse)
		if !ok {
			log.Fatalf("failed to cast entities response: %v", res)
		}

		for pid, entity := range entities.Entities {
			fmt.Printf("Entity %v: %v\n", pid, entity)
		}
	})
}
