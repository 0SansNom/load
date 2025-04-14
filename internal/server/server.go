package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"yourmodule/internal/balancer"
	"yourmodule/internal/health"
	"yourmodule/internal/logger"
	"yourmodule/internal/proxy"
)

// Server struct holds the dependencies to run the server.
type Server struct {
	httpServer *http.Server
}

// New creates a new Server instance with all dependencies initialized.
func New(balancer *balancer.RoundRobin, healthChecker *health.HealthChecker, log *logger.Logger) *Server {
	// Create a proxy that uses the round-robin load balancer
	proxy := proxy.New(balancer, log)

	// Create a new health checker
	healthChecker.Start()

	// Create the HTTP server with the proxy handler
	mux := http.NewServeMux()
	mux.Handle("/", proxy) // Handle all requests via proxy

	httpServer := &http.Server{
		Addr:    ":8080", // Set server address
		Handler: mux,
	}

	return &Server{
		httpServer: httpServer,
	}
}

// Start begins the HTTP server and gracefully shuts it down on interrupt signals.
func (s *Server) Start() error {
	// Set up a signal channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	// Run server in a separate goroutine
	go func() {
