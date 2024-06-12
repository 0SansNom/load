package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
)

// BackendServer représente un serveur backend
type BackendServer struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

// ServerPool contient les serveurs backend
type ServerPool struct {
	servers []*BackendServer
	current uint64
}

// AddServer ajoute un serveur backend à la pool
func (s *ServerPool) AddServer(server *BackendServer) {
	s.servers = append(s.servers, server)
}

// NextServer retourne le prochain serveur backend à utiliser en utilisant le round-robin
func (s *ServerPool) NextServer() *BackendServer {
	next := atomic.AddUint64(&s.current, uint64(1))
	return s.servers[next%uint64(len(s.servers))]
}

// ProxyHandler est le gestionnaire HTTP qui fait office de proxy
func ProxyHandler(pool *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server := pool.NextServer()
		if server != nil {
			server.ReverseProxy.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
	}
}

func main() {
	// URLs des serveurs backend
	backendURLs := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	var servers []*BackendServer

	for _, backendURL := range backendURLs {
		url, err := url.Parse(backendURL)
		if err != nil {
			log.Fatalf("Could not parse URL %s: %v", backendURL, err)
		}
		servers = append(servers, &BackendServer{
			URL:          url,
			Alive:        true,
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
		})
	}

	serverPool := &ServerPool{servers: servers}

	http.HandleFunc("/", ProxyHandler(serverPool))

	log.Println("Starting proxy server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
