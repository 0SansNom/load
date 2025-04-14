package health

import (
	"net/http"
	"sync"
	"time"

	"yourmodule/internal/balancer"
	"yourmodule/internal/logger"
)

// HealthChecker periodically checks the health of each backend.
type HealthChecker struct {
	balancer          *balancer.RoundRobin
	httpClient        *http.Client
	checkInterval     time.Duration
	checkTimeout      time.Duration
	backendHealth     map[string]bool
	mu                sync.Mutex
	log               *logger.Logger
}

// NewHealthChecker creates a new HealthChecker.
func NewHealthChecker(balancer *balancer.RoundRobin, log *logger.Logger, checkInterval, checkTimeout time.Duration) *HealthChecker {
	return &HealthChecker{
		balancer:      balancer,
		httpClient:    &http.Client{Timeout: checkTimeout},
		checkInterval: checkInterval,
		checkTimeout:  checkTimeout,
		backendHealth: make(map[string]bool),
		log:           log,
	}
}

// Start begins periodic health checks for each backend.
func (h *HealthChecker) Start() {
	ticker := time.NewTicker(h.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.checkBackends()
		}
	}
}

// checkBackends checks the health of all backends and updates their status.
func (h *HealthChecker) checkBackends() {
	h.mu.Lock()
	defer h.mu.Unlock()

	backends := h.balancer.GetBackends()
	for _, backend := range backends {
		go h.checkBackend(backend)
	}
}

// checkBackend performs an HTTP GET request to check if the backend is healthy.
func (h *HealthChecker) checkBackend(backend string) {
	resp, err := h.httpClient.Get(backend)
	if err != nil || resp.StatusCode != http.StatusOK {
		h.backendHealth[backend] = false
		h.log.Error("Backend is unhealthy", zap.String("url", backend), zap.Error(err))
		return
	}

	h.backendHealth[backend] = true
	h.log.Info("Backend is healthy", zap.String("url", backend))
}

// IsHealthy returns whether a given backend is healthy or not.
func (h *HealthChecker) IsHealthy(backend string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.backendHealth[backend]
}
