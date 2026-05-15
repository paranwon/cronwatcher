package api

import (
	"encoding/json"
	"net/http"
)

// handleJobAnnotations handles GET and POST/DELETE for job annotations.
//
// GET  /api/jobs/annotations?name=<job>         — returns all annotations
// POST /api/jobs/annotations?name=<job>         — merges annotations (JSON body map[string]string)
// DELETE /api/jobs/annotations?name=<job>&key=  — removes a single annotation key
func (s *Server) handleJobAnnotations(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing 'name' query parameter", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		anns, ok := s.watcher.GetAnnotations(name)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(anns)

	case http.MethodPost:
		var payload map[string]string
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		if err := s.watcher.SetAnnotations(name, payload); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	case http.MethodDelete:
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing 'key' query parameter", http.StatusBadRequest)
			return
		}
		if err := s.watcher.DeleteAnnotation(name, key); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
