package api

import (
	"encoding/json"
	"net/http"
)

// handleJobPause handles POST /jobs/{name}/pause and POST /jobs/{name}/resume
func (s *Server) handleJobPause(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}

	action := r.URL.Query().Get("action")
	if action != "pause" && action != "resume" {
		http.Error(w, "action must be pause or resume", http.StatusBadRequest)
		return
	}

	var err error
	if action == "pause" {
		err = s.watcher.PauseJob(name)
	} else {
		err = s.watcher.ResumeJob(name)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"job":    name,
		"status": action + "d",
	})
}
