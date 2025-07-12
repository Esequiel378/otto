package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"otto/system"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/go-gl/mathgl/mgl64"
)

// MockRenderer is a simple mock renderer that just counts messages
type MockRenderer struct {
	messageCount int64
}

func NewMockRenderer() actor.Producer {
	return func() actor.Receiver {
		return &MockRenderer{
			messageCount: 0,
		}
	}
}

func (r *MockRenderer) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Initialized:
		log.Printf("MockRenderer initialized")
	case system.EventEntityRenderUpdate:
		atomic.AddInt64(&r.messageCount, 1)
		if atomic.LoadInt64(&r.messageCount)%1000 == 0 {
			log.Printf("Renderer received %d messages", atomic.LoadInt64(&r.messageCount))
		}
	case system.RequestEntities:
		c.Respond(system.EntitiesResponse{
			Entities: []system.EntityRigidBody{},
		})
	}
}

// MockPhysics is a simple mock physics system
type MockPhysics struct {
	messageCount *int64
}

func NewMockPhysics(messageCounter *int64) actor.Producer {
	return func() actor.Receiver {
		return &MockPhysics{
			messageCount: messageCounter,
		}
	}
}

func (m *MockPhysics) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Initialized:
		log.Printf("MockPhysics initialized")
	case system.EventEntityUpdate:
		atomic.AddInt64(m.messageCount, 1)
		// Simulate some physics processing
		time.Sleep(time.Microsecond * 10)
		if atomic.LoadInt64(m.messageCount)%1000 == 0 {
			log.Printf("Physics received %d messages", atomic.LoadInt64(m.messageCount))
		}
	case system.SetCameraPID:
		// Just acknowledge the camera PID
	}
}

// MockCamera is a simple mock camera
type MockCamera struct {
	camera system.Camera
}

func NewMockCamera() actor.Producer {
	return func() actor.Receiver {
		return &MockCamera{
			camera: system.Camera{
				Position: mgl64.Vec3{0, 0, -2},
				Rotation: mgl64.Vec2{0, 0},
				Zoom:     1.0,
			},
		}
	}
}

func (c *MockCamera) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Initialized:
		log.Printf("MockCamera initialized")
	case system.RequestCamera:
		ctx.Respond(system.ResponseCamera{
			Camera: c.camera,
		})
	}
}

// BenchmarkPlayer is a player that generates random movements for benchmarking
type BenchmarkPlayer struct {
	physicsPID   *actor.PID
	rendererPID  *actor.PID
	playerID     int
	messageCount int64
}

func NewBenchmarkPlayer(physicsPID, rendererPID *actor.PID, playerID int) actor.Producer {
	return func() actor.Receiver {
		return &BenchmarkPlayer{
			physicsPID:   physicsPID,
			rendererPID:  rendererPID,
			playerID:     playerID,
			messageCount: 0,
		}
	}
}

func (p *BenchmarkPlayer) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Initialized:
		log.Printf("BenchmarkPlayer %d initialized", p.playerID)
	case system.Tick:
		// Generate random movement for this player
		velocity := mgl64.Vec3{
			(rand.Float64() - 0.5) * 2, // Random between -1 and 1
			(rand.Float64() - 0.5) * 2,
			(rand.Float64() - 0.5) * 2,
		}

		// Normalize velocity
		if velocity.Len() > 0 {
			velocity = velocity.Normalize()
		}

		// Send movement update to physics
		c.Send(p.physicsPID, system.EventEntityUpdate{
			PID:      c.PID(),
			Velocity: velocity,
		})

		// Send render update to renderer
		c.Send(p.rendererPID, system.EventEntityRenderUpdate{
			PID: c.PID(),
			EntityRigidBody: system.EntityRigidBody{
				Position: mgl64.Vec3{0, 0, 0}, // Mock position
				Velocity: velocity,
				Scale:    mgl64.Vec3{1, 1, 1},
				Rotation: mgl64.Vec3{0, 0, 0},
			},
		})

		atomic.AddInt64(&p.messageCount, 1)
	}
}

// BenchmarkStats tracks performance metrics
type BenchmarkStats struct {
	TotalTicks    int64
	TotalMessages int64
	StartTime     time.Time
	EndTime       time.Time
	TickRate      int
	NumPlayers    int
}

func (s *BenchmarkStats) Print() {
	duration := s.EndTime.Sub(s.StartTime)
	messagesPerSecond := float64(s.TotalMessages) / duration.Seconds()
	ticksPerSecond := float64(s.TotalTicks) / duration.Seconds()

	fmt.Printf("\n=== BENCHMARK RESULTS ===\n")
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Tick Rate: %d Hz\n", s.TickRate)
	fmt.Printf("Number of Players: %d\n", s.NumPlayers)
	fmt.Printf("Total Ticks: %d\n", s.TotalTicks)
	fmt.Printf("Total Messages: %d\n", s.TotalMessages)
	fmt.Printf("Messages per second: %.2f\n", messagesPerSecond)
	fmt.Printf("Ticks per second: %.2f\n", ticksPerSecond)
	fmt.Printf("Messages per tick: %.2f\n", float64(s.TotalMessages)/float64(s.TotalTicks))
	fmt.Printf("========================\n")
}

func init() {
	runtime.LockOSThread()
}

func main() {
	var (
		duration = flag.Duration("duration", 10*time.Second, "Benchmark duration")
		tickRate = flag.Int("tickrate", 64, "Tick rate in Hz (64 or 128)")
		players  = flag.Int("players", 100, "Number of players to simulate")
	)
	flag.Parse()

	// Validate tick rate
	if *tickRate != 64 && *tickRate != 128 {
		log.Fatal("Tick rate must be either 64 or 128")
	}

	log.Printf("Starting benchmark with %d players, %d Hz tick rate, for %v", *players, *tickRate, *duration)

	// Create actor engine
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatalf("failed to create actor engine: %v", err)
	}

	// Initialize stats and message counters
	var totalMessages int64
	stats := &BenchmarkStats{
		TickRate:   *tickRate,
		NumPlayers: *players,
		StartTime:  time.Now(),
	}

	// Create mock systems with shared message counter
	rendererPID := e.Spawn(NewMockRenderer(), "renderer")
	physicsPID := e.Spawn(NewMockPhysics(&totalMessages), "physics")
	cameraPID := e.Spawn(NewMockCamera(), "camera")

	// Set camera PID in physics system
	e.Send(physicsPID, system.SetCameraPID{PID: cameraPID})

	// Create benchmark players
	playerPIDs := make([]*actor.PID, *players)
	for i := 0; i < *players; i++ {
		playerPIDs[i] = e.Spawn(NewBenchmarkPlayer(physicsPID, rendererPID, i), fmt.Sprintf("player_%d", i))
	}

	tickInterval := time.Second / time.Duration(*tickRate)
	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	// Create a channel to signal when the benchmark should stop
	done := make(chan bool)
	go func() {
		time.Sleep(*duration)
		done <- true
	}()

	log.Printf("Starting tick loop...")
	latestTick := time.Now()

	// Main tick loop
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			deltaTime := now.Sub(latestTick).Seconds()
			latestTick = now

			atomic.AddInt64(&stats.TotalTicks, 1)
			e.BroadcastEvent(system.Tick{DeltaTime: deltaTime})

		case <-done:
			stats.EndTime = time.Now()

			// Get physics message count
			physicsMessages := atomic.LoadInt64(&totalMessages)

			// Calculate total messages: each player sends 2 messages per tick (physics + renderer)
			expectedMessagesPerTick := int64(*players) * 2
			stats.TotalMessages = physicsMessages + (stats.TotalTicks * expectedMessagesPerTick)

			// Print results
			stats.Print()
			return
		}
	}
}
