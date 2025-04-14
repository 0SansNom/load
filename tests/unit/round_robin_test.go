package balancer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"yourmodule/internal/balancer"
)

func TestNewRoundRobin(t *testing.T) {
	// Test with valid backends
	backends := []string{"http://localhost:9001", "http://localhost:9002"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)
	assert.Equal(t, 2, len(rb.backends))

	// Test with empty backends
	rb, err = balancer.NewRoundRobin([]string{})
	require.Error(t, err)
	assert.Nil(t, rb)
}

func TestRoundRobin_Next(t *testing.T) {
	backends := []string{"http://localhost:9001", "http://localhost:9002"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Test that backends rotate correctly
	assert.Equal(t, "http://localhost:9001", rb.Next())
	assert.Equal(t, "http://localhost:9002", rb.Next())
	assert.Equal(t, "http://localhost:9001", rb.Next())
}

func TestRoundRobin_SetBackends(t *testing.T) {
	backends := []string{"http://localhost:9001", "http://localhost:9002"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Change backends
	newBackends := []string{"http://localhost:9003", "http://localhost:9004"}
	err = rb.SetBackends(newBackends)
	require.NoError(t, err)

	// Verify that the backends were updated
	assert.Equal(t, "http://localhost:9003", rb.Next())
	assert.Equal(t, "http://localhost:9004", rb.Next())
}

func TestRoundRobin_EmptyBackends(t *testing.T) {
	backends := []string{"http://localhost:9001", "http://localhost:9002"}
	rb, err := balancer.NewRoundRobin(backends)
	require.NoError(t, err)

	// Test empty backend list after update
	err = rb.SetBackends([]string{})
	require.Error(t, err)
}
