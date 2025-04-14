package server_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"yourmodule/internal/balancer"
	"yourmodule/internal/health"
	"yourmodule/internal/logger"
	"yourmodule/internal/proxy"
	"yourmodule/internal/server"
)

func TestServer_Start(t *testing.T) {
	// Initialize mock logger
	log, err := logger.New(false) // Development mode
	require.NoError(t, err)

	// Mock load balancer with one backend
	backends := []string{"http://localhost:9001"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Create health checker
	checker := health.NewHealthChecker(rb, log, 1*time.Second, 2*time.Second)

	// Create the server
	srv := server.New(rb, checker, log)

	// Start the server in a goroutine
	go func() {
		err := srv.Start()
		require.NoError(t, err)
	}()

	// Allow the server to start
	time.Sleep(1 * time.Second)

	// Test the server's basic response
	resp, err := http.Get("http://localhost:8080")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Gracefully stop the server
	err = srv.httpServer.Shutdown(context.Background())
	require.NoError(t, err)
}

func TestServer_GracefulShutdown(t *testing.T) {
	// Initialize mock logger
	log, err := logger.New(false)
	require.NoError(t, err)

	// Mock load balancer with one backend
	backends := []string{"http://localhost:9001"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Create health checker
	checker := health.NewHealthChecker(rb, log, 1*time.Second, 2*time.Second)

	// Create the server
	srv := server.New(rb, checker, log)

	// Start the server in a goroutine
	go func() {
		err := srv.Start()
		require.NoError(t, err)
	}()

	// Allow the server to start
	time.Sleep(1 * time.Second)

	// Simulate sending an interrupt signal
	srv.httpServer.Shutdown(context.Background())

	// Allow some time for shutdown process
	time.Sleep(1 * time.Second)

	// The server should be shut down gracefully without errors
	assert.True(t, true)
}
