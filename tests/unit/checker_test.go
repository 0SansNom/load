package health_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"yourmodule/internal/balancer"
	"yourmodule/internal/health"
	"yourmodule/internal/logger"
)

func TestHealthChecker(t *testing.T) {
	// Create a test logger
	log, err := logger.New(false) // Development mode
	require.NoError(t, err)

	// Mock a load balancer with a few backends
	backends := []string{"http://localhost:9001", "http://localhost:9002"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Create a health checker
	checker := health.NewHealthChecker(rb, log, 1*time.Second, 2*time.Second)

	// Start the health checker in a separate goroutine
	go checker.Start()

	// Set up mock HTTP servers for the backends
	mockBackend1 := setupMockBackend(t, http.StatusOK)
	defer mockBackend1.Close()

	mockBackend2 := setupMockBackend(t, http.StatusInternalServerError)
	defer mockBackend2.Close()

	// Update backends to point to mock servers
	rb.SetBackends([]string{mockBackend1.URL, mockBackend2.URL})

	// Allow some time for health checks to complete
	time.Sleep(3 * time.Second)

	// Check that the health status is updated correctly
	assert.True(t, checker.IsHealthy(mockBackend1.URL))
	assert.False(t, checker.IsHealthy(mockBackend2.URL))
}

// setupMockBackend creates a mock HTTP server that returns the given status code.
func setupMockBackend(t *testing.T, statusCode int) *http.Server {
	t.Helper()

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	})

	server := &http.Server{
		Addr:    "localhost:0", // Random port
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("could not start mock server: %v", err)
		}
	}()

	// Return the server's URL
	return server
}
