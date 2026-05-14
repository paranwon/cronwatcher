package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatcher/internal/watcher"
)

// Server exposes a lightweight HTTP API for querying job status.
type Server struct {
	watcher *watcher.Watcher
	addr    string
	server  *http.Server
}

// JobStatus represents the current status of a monitored cron job.
type JobStatus struct {
	Name      string     `json:"name"`
	Running   bool       `json:"running"`
	LastStart *time.Time `json:"last_start,omitempty"`
	LastEnd   *time.Time `json:"last_end,omitempty"`
	LastDuration *string `json:"last_duration,omitempty"`
}

// New creates a new API Server.
func New(w *watcher.Watcher, addr string) *Server {
	s := &Server{watcher: w, addr: addr}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/status", s.handleStatus)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return s
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop() error {
	return s.server.Close()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	statuses := s.watcher.Status()
	result := make([]JobStatus, 0, len(statuses))
	for _, js := range statuses {
		result = append(result, js)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
