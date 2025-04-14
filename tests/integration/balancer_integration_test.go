package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/0SansNom/load/internal/balancer"

)

func TestLoadBalancerRoundRobin(t *testing.T) {
	// Create mock backends using httptest
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Backend 1"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Backend 2"))
	}))
	defer server2.Close()

	// Configure the load balancer with the mock backends
	backends := []string{server1.URL, server2.URL}
	lb, err := balancer.NewRoundRobin(backends)
	assert.NoError(t, err)

	// Test the load balancing by sending two requests
	resp1, err := lb.Next()
	assert.NoError(t, err)
	assert.Equal(t, "Backend 1", resp1)

	resp2, err := lb.Next()
	assert.NoError(t, err)
	assert.Equal(t, "Backend 2", resp2)

	// Ensure it wraps around correctly
	resp3, err := lb.Next()
	assert.NoError(t, err)
	assert.Equal(t, "Backend 1", resp3)
}
