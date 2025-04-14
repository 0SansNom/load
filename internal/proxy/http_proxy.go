package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"yourmodule/internal/balancer"
	"yourmodule/internal/logger"
)

// Proxy is a wrapper around the HTTP reverse proxy with load balancing.
type Proxy struct {
	balancer   *balancer.RoundRobin
	httpClient *http.Client
	log        *logger.Logger
}

// New creates a new Proxy instance with the provided RoundRobin balancer.
func New(balancer *balancer.RoundRobin, log *logger.Logger) *Proxy {
	return &Proxy{
		balancer:   balancer,
		httpClient: &http.Client{},
		log:        log,
	}
}

// ServeHTTP is the main entry point for handling HTTP requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := p.balancer.Next()

	// Log the request routing
	p.log.Info("Routing request", zap.String("url", backend), zap.String("method", r.Method))

	// Parse the backend URL
	parsedURL, err := url.Parse(backend)
	if err != nil {
		http.Error(w, "Failed to parse backend URL", http.StatusInternalServerError)
		return
	}

	// Create a reverse proxy for the selected backend
	proxy := httputil.NewSingleHostReverseProxy(parsedURL)

	// Forward the request to the selected backend
	proxy.ServeHTTP(w, r)
}
