package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
)

type BackendServer struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}


type ServerPool struct {
	servers []*BackendServer
	current uint64
}


func (s *ServerPool) AddServer(server *BackendServer) {
	s.servers = append(s.servers, server)
}


func (s *ServerPool) NextServer() *BackendServer {
	next := atomic.AddUint64(&s.current, uint64(1))
	return s.servers[next%uint64(len(s.servers))]
}

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
