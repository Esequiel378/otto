package monitoring

import (
	"net/http"
	"testing"
	"time"
)

func TestMetricsManager(t *testing.T) {
	// Test with metrics disabled
	manager := NewMetricsManager()
	if manager.IsEnabled() {
		t.Error("Metrics should be disabled by default")
	}

	// Test metrics methods don't panic when disabled
	manager.UpdateFPS(60.0)
	manager.UpdateEntityCount(100)
	manager.IncrementRenderCalls()
	manager.IncrementInputEvents()
	manager.IncrementPhysicsCalculations()
	manager.RecordFrameTime(time.Millisecond * 16)
	manager.UpdateMemoryUsage(1024, 2048, 512, 1536)

	// Test stop doesn't panic
	if err := manager.Stop(); err != nil {
		t.Errorf("Stop should not return error when disabled: %v", err)
	}
}

func TestMetricsManagerEnabled(t *testing.T) {
	// This test requires the environment variable to be set manually
	// Run with: OTTO_METRICS_ENABLED=true go test ./monitoring -run TestMetricsManagerEnabled

	manager := NewMetricsManager()

	// Only run the full test if metrics are enabled
	if manager.IsEnabled() {
		// Test start
		if err := manager.Start(); err != nil {
			t.Errorf("Start failed: %v", err)
		}

		// Give the server a moment to start
		time.Sleep(100 * time.Millisecond)

		// Test metrics endpoint is accessible
		resp, err := http.Get("http://localhost:8080/metrics")
		if err != nil {
			t.Errorf("Failed to access metrics endpoint: %v", err)
		} else {
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
		}

		// Test stop
		if err := manager.Stop(); err != nil {
			t.Errorf("Stop failed: %v", err)
		}
	} else {
		t.Skip("Skipping test - metrics not enabled. Run with OTTO_METRICS_ENABLED=true")
	}
}
