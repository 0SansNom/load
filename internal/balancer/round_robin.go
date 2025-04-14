package balancer

import (
	"errors"
	"sync"
)

// RoundRobin is a simple round-robin load balancer that selects backends in order.
type RoundRobin struct {
	backends []string
	mu       sync.Mutex
	counter  int
}

// NewRoundRobin creates a new RoundRobin load balancer.
func NewRoundRobin(backends []string) (*RoundRobin, error) {
	if len(backends) == 0 {
		return nil, errors.New("at least one backend is required")
	}
	return &RoundRobin{backends: backends}, nil
}

// Next returns the next backend in a round-robin fashion.
func (r *RoundRobin) Next() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	backend := r.backends[r.counter]
	r.counter = (r.counter + 1) % len(r.backends)
	return backend
}

// SetBackends allows updating the backends list at runtime.
func (r *RoundRobin) SetBackends(backends []string) error {
	if len(backends) == 0 {
		return errors.New("at least one backend is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.backends = backends
	r.counter = 0 // Reset counter when backends are updated
	return nil
}
