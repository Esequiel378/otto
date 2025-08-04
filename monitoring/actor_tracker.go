package monitoring

import (
	"os"
	"sync"

	"github.com/anthdm/hollywood/actor"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Actor message counter
	actorMessageCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "otto_actor_messages_total",
		Help: "Total number of messages received by actor",
	}, []string{"actor_name", "message_type"})

	// Actor message rate
	actorMessageRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "otto_actor_messages_per_second",
		Help: "Messages per second for each actor",
	}, []string{"actor_name"})

	// Track actor PIDs and names for cleanup
	actorPIDs  = make(map[string]*actor.PID)
	actorNames = make(map[string]string) // PID string -> actor name
	actorMutex sync.RWMutex
)

// ActorTracker tracks messages per actor
type ActorTracker struct {
	enabled bool
}

// NewActorTracker creates a new actor tracker
func NewActorTracker() *ActorTracker {
	enabled := os.Getenv("OTTO_METRICS_ENABLED") == "true"

	if enabled {
		// Register metrics
		registry.MustRegister(actorMessageCounter)
		registry.MustRegister(actorMessageRate)
	}

	return &ActorTracker{
		enabled: enabled,
	}
}

// WithActorTracking creates middleware to track actor messages
func (at *ActorTracker) WithActorTracking(actorName string) func(actor.ReceiveFunc) actor.ReceiveFunc {
	return func(next actor.ReceiveFunc) actor.ReceiveFunc {
		return func(c *actor.Context) {
			if at.enabled {
				// Track actor creation
				switch c.Message().(type) {
				case actor.Initialized:
					at.trackActor(c, actorName)
				case actor.Stopped:
					at.untrackActor(c)
				}

				// Count messages
				at.countMessage(c)
			}
			next(c)
		}
	}
}

// trackActor adds an actor to tracking
func (at *ActorTracker) trackActor(c *actor.Context, actorName string) {
	actorMutex.Lock()
	defer actorMutex.Unlock()

	pid := c.PID()
	pidStr := pid.String()

	actorPIDs[pidStr] = pid
	actorNames[pidStr] = actorName

	// Initialize message rate tracking
	actorMessageRate.WithLabelValues(actorName).Set(0)
}

// untrackActor removes an actor from tracking
func (at *ActorTracker) untrackActor(c *actor.Context) {
	actorMutex.Lock()
	defer actorMutex.Unlock()

	pid := c.PID()
	pidStr := pid.String()
	delete(actorPIDs, pidStr)
	delete(actorNames, pidStr)
}

// countMessage increments the message counter for an actor
func (at *ActorTracker) countMessage(c *actor.Context) {
	pid := c.PID()
	pidStr := pid.String()

	// Get actor name from stored names, or use default
	actorName := "unknown"
	actorMutex.RLock()
	if name, exists := actorNames[pidStr]; exists {
		actorName = name
	}
	actorMutex.RUnlock()

	messageType := "unknown"
	switch c.Message().(type) {
	case actor.Initialized:
		messageType = "initialized"
	case actor.Started:
		messageType = "started"
	case actor.Stopped:
		messageType = "stopped"
	default:
		messageType = "custom"
	}

	actorMessageCounter.WithLabelValues(actorName, messageType).Inc()
}

// GetActorCount returns the current number of tracked actors
func (at *ActorTracker) GetActorCount() int {
	actorMutex.RLock()
	defer actorMutex.RUnlock()
	return len(actorPIDs)
}

// GetActorMessageCount returns message count for a specific actor
func (at *ActorTracker) GetActorMessageCount(actorName, actorID string) float64 {
	// This is a simplified implementation
	// In a real implementation, you'd want to track this properly
	return 0
}

// UpdateActorMessageRates updates the message rate metrics
func (at *ActorTracker) UpdateActorMessageRates() {
	if !at.enabled {
		return
	}

	actorMutex.RLock()
	defer actorMutex.RUnlock()

	// This is a simplified rate calculation
	// In a real implementation, you'd want to track time windows
	for pidStr := range actorPIDs {
		actorName := "unknown"                                     // You'd need to store actor names
		actorMessageRate.WithLabelValues(actorName, pidStr).Set(0) // Placeholder
	}
}
