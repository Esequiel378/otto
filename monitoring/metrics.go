package monitoring

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Game-specific metrics
	fpsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "otto_fps",
		Help: "Current frames per second",
	})

	entityCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "otto_entity_count",
		Help: "Number of entities in the game world",
	})

	renderCallsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "otto_render_calls_total",
		Help: "Total number of render calls",
	})

	inputEventsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "otto_input_events_total",
		Help: "Total number of input events processed",
	})

	physicsCalculationsCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "otto_physics_calculations_total",
		Help: "Total number of physics calculations",
	})

	frameTimeHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "otto_frame_time_seconds",
		Help:    "Frame rendering time in seconds",
		Buckets: prometheus.DefBuckets,
	})

	memoryUsageGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "otto_memory_usage_bytes",
		Help: "Memory usage in bytes",
	}, []string{"type"})

	// Metrics registry
	registry = prometheus.NewRegistry()
)

// MetricsManager handles all monitoring functionality
type MetricsManager struct {
	enabled bool
	server  *http.Server
}

// NewMetricsManager creates a new metrics manager
func NewMetricsManager() *MetricsManager {
	enabled := os.Getenv("OTTO_METRICS_ENABLED") == "true"

	if enabled {
		// Register all metrics
		registry.MustRegister(fpsGauge)
		registry.MustRegister(entityCountGauge)
		registry.MustRegister(renderCallsCounter)
		registry.MustRegister(inputEventsCounter)
		registry.MustRegister(physicsCalculationsCounter)
		registry.MustRegister(frameTimeHistogram)
		registry.MustRegister(memoryUsageGauge)

		// Register default Go metrics
		registry.MustRegister(prometheus.NewGoCollector())
		registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	}

	return &MetricsManager{
		enabled: enabled,
	}
}

// Start starts the metrics HTTP server
func (m *MetricsManager) Start() error {
	if !m.enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	m.server = &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return nil
}

// Stop stops the metrics HTTP server
func (m *MetricsManager) Stop() error {
	if !m.enabled || m.server == nil {
		return nil
	}
	return m.server.Close()
}

// IsEnabled returns whether metrics are enabled
func (m *MetricsManager) IsEnabled() bool {
	return m.enabled
}

// UpdateFPS updates the FPS metric
func (m *MetricsManager) UpdateFPS(fps float64) {
	if m.enabled {
		fpsGauge.Set(fps)
	}
}

// UpdateEntityCount updates the entity count metric
func (m *MetricsManager) UpdateEntityCount(count int) {
	if m.enabled {
		entityCountGauge.Set(float64(count))
	}
}

// IncrementRenderCalls increments the render calls counter
func (m *MetricsManager) IncrementRenderCalls() {
	if m.enabled {
		renderCallsCounter.Inc()
	}
}

// IncrementInputEvents increments the input events counter
func (m *MetricsManager) IncrementInputEvents() {
	if m.enabled {
		inputEventsCounter.Inc()
	}
}

// IncrementPhysicsCalculations increments the physics calculations counter
func (m *MetricsManager) IncrementPhysicsCalculations() {
	if m.enabled {
		physicsCalculationsCounter.Inc()
	}
}

// RecordFrameTime records the frame rendering time
func (m *MetricsManager) RecordFrameTime(duration time.Duration) {
	if m.enabled {
		frameTimeHistogram.Observe(duration.Seconds())
	}
}

// UpdateMemoryUsage updates memory usage metrics
func (m *MetricsManager) UpdateMemoryUsage(heapAlloc, heapSys, heapIdle, heapInuse uint64) {
	if m.enabled {
		memoryUsageGauge.WithLabelValues("heap_alloc").Set(float64(heapAlloc))
		memoryUsageGauge.WithLabelValues("heap_sys").Set(float64(heapSys))
		memoryUsageGauge.WithLabelValues("heap_idle").Set(float64(heapIdle))
		memoryUsageGauge.WithLabelValues("heap_inuse").Set(float64(heapInuse))
	}
}
