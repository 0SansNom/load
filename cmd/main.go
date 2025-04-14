package main

import (
	"log"
	"time"

	"github.com/0sansnom/load/internal/balancer"
	"github.com/0sansnom/load/internal/health"
	"github.com/0sansnom/load/internal/logger"
	"github.com/0sansnom/load/internal/proxy"
	"github.com/0sansnom/load/internal/server"
)

func main() {
	// Initialize logger
	log, err := logger.New(true) // Development mode
	if err != nil {
		log.Fatal("could not initialize logger", err)
	}

	// Configure the list of backends (you can load this dynamically from a config)
	backends := []string{
		"http://localhost:9001",
		"http://localhost:9002",
	}

	// Create the load balancer (Round Robin)
	rb, err := balancer.NewRoundRobin(backends)
	if err != nil {
		log.Fatal("could not create load balancer", err)
	}

	// Create health checker and start it
	checker := health.NewHealthChecker(rb, log, 1*time.Second, 2*time.Second)
	checker.Start()

	// Create the proxy (it will use the load balancer)
	proxy := proxy.New(rb, log)

	// Create the server
	httpServer := server.New(rb, checker, log)

	// Start the server
	log.Info("starting server on port 8080")
	if err := httpServer.Start(); err != nil {
		log.Fatal("server failed", err)
	}
}
