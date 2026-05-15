package api

import (
	"net/http"
)

// RegisterRoutes attaches all API routes to the given mux.
func RegisterRoutes(mux *http.ServeMux, s *Server) {
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/status", s.handleStatus)
	mux.HandleFunc("/jobs/", s.handleJobDetail)
}
