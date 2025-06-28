package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)
// Backend represents a backend server
type Backend struct {
	URL    *url.URL
	Alive  bool
	Mux    sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive sets the alive status of the backend
func (b *Backend) SetAlive(alive bool) {
	b.Mux.Lock()
	b.Alive = alive
	b.Mux.Unlock()
}

// IsAlive returns the alive status of the backend
func (b *Backend) IsAlive() (alive bool) {
	b.Mux.RLock()
	alive = b.Alive
	b.Mux.RUnlock()
	return
}

// ServerPool holds information about reachable backend servers
type ServerPool struct {
	backends []*Backend
	current  uint64
	Mux      sync.RWMutex
}

// AddBackend adds a new backend to the server pool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.Mux.Lock()
	s.backends = append(s.backends, backend)
	s.Mux.Unlock()
}

// NextIndex returns the next active backend to use for a request
func (s *ServerPool) NextIndex() int {
	s.Mux.Lock()
	defer s.Mux.Unlock()

	// Simple round-robin
	s.current = (s.current + 1) % uint64(len(s.backends))
	return int(s.current)
}

// GetNextPeer returns the next active backend to use for a request
func (s *ServerPool) GetNextPeer() *Backend {
	// Loop until we find an alive backend or exhaust all backends
	for i := 0; i < len(s.backends); i++ {
		next := s.NextIndex()
		backend := s.backends[next]
		if backend.IsAlive() {
			return backend
		}
	}
	return nil // No alive backends
}

// HealthCheck checks the health of backend servers
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"
		// Implement a simple health check by trying to connect to the backend
		res, err := http.Get(b.URL.String() + "/health") // Assuming a /health endpoint
		if err != nil || res.StatusCode != http.StatusOK {
			log.Printf("Backend %s is down: %v\n", b.URL, err)
			b.SetAlive(false)
			status = "down"
		} else {
			b.SetAlive(true)
		}
		log.Printf("Backend %s status: %s\n", b.URL, status)
	}
}

// lbHandler is the main handler for the load balancer
func lbHandler(w http.ResponseWriter, r *http.Request) {
	peer := serverPool.GetNextPeer()
	if peer == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	log.Printf("Forwarding request to backend: %s\n", peer.URL)
	peer.ReverseProxy.ServeHTTP(w, r)
}

var serverPool ServerPool

func main() {
	// Define backend URLs
	backendURLs := []string{
		"http://localhost:9001",
		"http://localhost:9002",
		"http://localhost:9003",
	}

	for _, u := range backendURLs {
		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Fatalf("Failed to parse URL %s: %v", u, err)
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("Error proxying to %s: %v\n", parsedURL, e)
			// Mark backend as down on error
			for _, b := range serverPool.backends {
				if b.URL.String() == parsedURL.String() {
					b.SetAlive(false)
					break
				}
			}
			lbHandler(writer, request) // Try next backend
		}

		serverPool.AddBackend(&Backend{
			URL:    parsedURL,
			Alive:  true, // Assume alive initially
			ReverseProxy: proxy,
		})
	}

	// Start health checking in a goroutine
	go func() {
		for {
			serverPool.HealthCheck()
			time.Sleep(10 * time.Second) // Check every 10 seconds
		}
	}()

	// Start the load balancer server
	log.Println("Starting load balancer on :8080")
	http.HandleFunc("/", lbHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}


