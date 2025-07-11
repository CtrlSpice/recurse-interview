package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server http.Server
	store  *Store
}

func NewServer(endpoint string) *Server {
	s := &Server{
		server: http.Server{
			Addr: endpoint,
		},
		store: NewStore(),
	}
	s.server.Handler = s.Handler()
	return s
}

func (s *Server) Start() {
	log.Printf("Server starting on %s", s.server.Addr)

	// Channel that listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// We listen for requests in a separate goroutine.
	go func() {
		serverErrors <- s.server.ListenAndServe()
	}()

	// Block main and wait for shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)
	case <-s.shutdown():
		log.Println("Shutting down...")

		// Give requests a second to complete
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Gracefully shutdown the server.
		if err := s.server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop server gracefully: %v", err)
			if err := s.server.Close(); err != nil {
				log.Fatalf("Could not force close server: %v", err)
			}
		}
	}
}

func (s *Server) Handler() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/set", s.setHandler)
	router.HandleFunc("/get", s.getHandler)
	return router
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Getting value for key '%s'", key)
	value := s.store.Get(key)
	if value == nil {
		log.Printf("Key '%s' not found", key)
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	log.Printf("Retrieved value '%s' for key '%s'", value, key)
	fmt.Fprintf(w, "%s", value)
}

func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	s.store.Set(key, value)

	response := fmt.Sprintf("Set key '%s' to value '%s'", key, value)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", response)
}

// shutdown returns a channel that listens for an interrupt/terminate signal from the OS,
// and will be closed when shutdown is requested.
func (s *Server) shutdown() chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	return shutdown
}
