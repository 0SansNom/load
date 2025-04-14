package proxy_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"yourmodule/internal/balancer"
	"yourmodule/internal/logger"
	"yourmodule/internal/proxy"
)

func TestProxy_ServeHTTP(t *testing.T) {
	// Initialize mock logger
	log, err := logger.New(false) // Development mode
	require.NoError(t, err)

	// Set up the load balancer with two mock backends
	mockBackend1 := setupMockBackend(t, http.StatusOK, "Hello from backend 1")
	defer mockBackend1.Close()

	mockBackend2 := setupMockBackend(t, http.StatusOK, "Hello from backend 2")
	defer mockBackend2.Close()

	backends := []string{mockBackend1.URL, mockBackend2.URL}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Create the proxy instance
	proxy := proxy.New(rb, log)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	require.NoError(t, err)

	// Create a test HTTP response recorder
	rr := &http.ResponseRecorder{}

	// Serve the HTTP request using the proxy
	proxy.ServeHTTP(rr, req)

	// Check that the response status code is 200 (OK)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify the response body comes from the first backend (round-robin selection)
	assert.Equal(t, "Hello from backend 1", rr.Body.String())
}

func TestProxy_RoundRobinSelection(t *testing.T) {
	// Initialize mock logger
	log, err := logger.New(false)
	require.NoError(t, err)

	// Set up two mock backends
	mockBackend1 := setupMockBackend(t, http.StatusOK, "Response from backend 1")
	defer mockBackend1.Close()

	mockBackend2 := setupMockBackend(t, http.StatusOK, "Response from backend 2")
	defer mockBackend2.Close()

	backends := []string{mockBackend1.URL, mockBackend2.URL}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Create the proxy instance
	proxy := proxy.New(rb, log)

	// Perform several requests to test round-robin backend selection
	for i := 0; i < 5; i++ {
		req, err := http.NewRequest("GET", "http://localhost:8080", nil)
		require.NoError(t, err)

		rr := &http.ResponseRecorder{}
		proxy.ServeHTTP(rr, req)

		if i%2 == 0 {
			assert.Equal(t, "Response from backend 1", rr.Body.String())
		} else {
			assert.Equal(t, "Response from backend 2", rr.Body.String())
		}
	}
}

func TestProxy_HealthCheck_Routing(t *testing.T) {
	// Initialize mock logger
	log, err := logger.New(false)
	require.NoError(t, err)

	// Set up one healthy and one unhealthy backend
	mockBackendHealthy := setupMockBackend(t, http.StatusOK, "Healthy backend")
	defer mockBackendHealthy.Close()

	mockBackendUnhealthy := setupMockBackend(t, http.StatusInternalServerError, "Unhealthy backend")
	defer mockBackendUnhealthy.Close()

	// Mock the load balancer to only use healthy backends
	backends := []string{mockBackendHealthy.URL, mockBackendUnhealthy.URL}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Create the proxy instance
	proxy := proxy.New(rb, log)

	// Perform a request to verify healthy backend is selected
	req, err := http.NewRequest("GET", "http://localhost:8080", nil)
	require.NoError(t, err)

	rr := &http.ResponseRecorder{}
	proxy.ServeHTTP(rr, req)

	// Verify that the healthy backend is used
	assert.Equal(t, "Healthy backend", rr.Body.String())
}

func setupMockBackend(t *testing.T, statusCode int, responseBody string) *http.Server {
	t.Helper()

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
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
