package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

// jobDetailResponse holds the response payload for a single job status.
type jobDetailResponse struct {
	Name        string `json:"name"`
	Running     bool   `json:"running"`
	LastSeen    string `json:"last_seen,omitempty"`
	LastElapsed string `json:"last_elapsed,omitempty"`
}

// handleJobDetail returns the status of a single job identified by name in the URL path.
// Path format: /jobs/{name}
func (s *Server) handleJobDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/jobs/")
	name = strings.TrimSpace(name)
	if name == "" {
		http.Error(w, "job name is required", http.StatusBadRequest)
		return
	}

	statuses := s.watcher.Status()
	for _, js := range statuses {
		if js.Name == name {
			resp := jobDetailResponse{
				Name:    js.Name,
				Running: js.Running,
			}
			if !js.LastSeen.IsZero() {
				resp.LastSeen = js.LastSeen.UTC().Format("2006-01-02T15:04:05Z")
			}
			if js.LastElapsed > 0 {
				resp.LastElapsed = js.LastElapsed.String()
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	http.Error(w, "job not found", http.StatusNotFound)
}
