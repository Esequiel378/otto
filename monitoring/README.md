# Otto Game Engine Monitoring System

This monitoring system provides comprehensive metrics for your Otto game engine using Prometheus and Grafana.

## Features

- **Game-specific metrics**: FPS, entity count, render calls, input events, physics calculations
- **System metrics**: Memory usage, Go runtime metrics, process metrics
- **Toggleable**: Enable/disable with environment variable
- **Docker-based**: Easy setup with Docker Compose
- **Real-time dashboards**: Beautiful Grafana dashboards

## Quick Start

### 1. Start the monitoring stack

```bash
make monitor-start
```

This will start:
- **Prometheus** on http://localhost:9090
- **Grafana** on http://localhost:3000 (admin/admin)

### 2. Run your game engine with metrics enabled

```bash
make run
```

Or manually:

```bash
OTTO_METRICS_ENABLED=true go run ./cmd/playground/main.go
```

### 3. View your metrics

- **Metrics endpoint**: http://localhost:8080/metrics
- **Grafana dashboard**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090

## Available Metrics

### Game-specific Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `otto_fps` | Gauge | Current frames per second |
| `otto_entity_count` | Gauge | Number of entities in the game world |
| `otto_render_calls_total` | Counter | Total number of render calls |
| `otto_input_events_total` | Counter | Total number of input events processed |
| `otto_physics_calculations_total` | Counter | Total number of physics calculations |
| `otto_frame_time_seconds` | Histogram | Frame rendering time distribution |
| `otto_memory_usage_bytes` | GaugeVec | Memory usage by type (heap_alloc, heap_sys, heap_idle, heap_inuse) |

### System Metrics (Automatic)

- Go runtime metrics (goroutines, GC stats, etc.)
- Process metrics (CPU, memory, file descriptors)
- Standard Prometheus metrics

## Integration Guide

### 1. Add monitoring to your main.go

```go
import "otto/monitoring"

func main() {
    // Initialize monitoring
    metricsManager := monitoring.NewMetricsManager()
    if err := metricsManager.Start(); err != nil {
        log.Printf("Warning: failed to start metrics server: %v", err)
    }
    defer metricsManager.Stop()
    
    // ... rest of your game engine code
}
```

### 2. Update metrics in your game loop

```go
// Update FPS
metricsManager.UpdateFPS(currentFPS)

// Update entity count
metricsManager.UpdateEntityCount(len(entities))

// Increment counters
metricsManager.IncrementRenderCalls()
metricsManager.IncrementInputEvents()
metricsManager.IncrementPhysicsCalculations()

// Record frame time
frameDuration := time.Since(frameStart)
metricsManager.RecordFrameTime(frameDuration)

// Update memory stats
var memStats runtime.MemStats
runtime.ReadMemStats(&memStats)
metricsManager.UpdateMemoryUsage(
    memStats.HeapAlloc,
    memStats.HeapSys,
    memStats.HeapIdle,
    memStats.HeapInuse,
)
```

### 3. Check if metrics are enabled

```go
if metricsManager.IsEnabled() {
    // Update metrics
    metricsManager.UpdateFPS(fps)
}
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make monitor-start` | Start Prometheus and Grafana |
| `make monitor-stop` | Stop the monitoring stack |
| `make monitor-logs` | Show monitoring stack logs |
| `make monitor-status` | Show monitoring stack status |
| `make run` | Run game engine with metrics enabled |

## Configuration

### Environment Variables

- `OTTO_METRICS_ENABLED=true` - Enable metrics collection
- `OTTO_METRICS_ENABLED=false` or unset - Disable metrics collection

### Prometheus Configuration

The Prometheus configuration is in `monitoring/prometheus.yml`. It scrapes metrics from:
- Your game engine: `host.docker.internal:8080/metrics`
- Prometheus itself: `localhost:9090`

### Grafana Configuration

- **Datasource**: Automatically configured to connect to Prometheus
- **Dashboard**: Pre-configured dashboard with all game metrics
- **Credentials**: admin/admin

## Dashboard Features

The Grafana dashboard includes:

1. **FPS Chart** - Real-time frames per second
2. **Entity Count** - Number of entities in the game world
3. **Memory Usage** - Heap allocation and system memory
4. **Performance Metrics** - Render calls, input events, physics calculations per second
5. **Frame Time Distribution** - 50th and 95th percentile frame times

## Troubleshooting

### Metrics not showing up

1. Check if metrics are enabled: `OTTO_METRICS_ENABLED=true`
2. Verify metrics endpoint: http://localhost:8080/metrics
3. Check Prometheus targets: http://localhost:9090/targets

### Docker issues

1. Make sure Docker is running
2. Check container status: `make monitor-status`
3. View logs: `make monitor-logs`

### Port conflicts

If ports 3000, 8080, or 9090 are in use:
- Change ports in `docker-compose.yml`
- Update Prometheus configuration in `monitoring/prometheus.yml`

## Development

### Adding new metrics

1. Add new metric in `monitoring/metrics.go`
2. Register it in `NewMetricsManager()`
3. Update it in your game engine code
4. Add to Grafana dashboard if needed

### Custom dashboards

1. Create new dashboard JSON in `monitoring/grafana/dashboards/`
2. Restart monitoring stack: `make monitor-stop && make monitor-start`

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Game Engine   │    │   Prometheus    │    │     Grafana     │
│                 │    │                 │    │                 │
│  :8080/metrics │───▶│  :9090          │───▶│  :3000          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

- Game engine exposes metrics on `/metrics` endpoint
- Prometheus scrapes metrics every 5 seconds
- Grafana queries Prometheus for visualization 