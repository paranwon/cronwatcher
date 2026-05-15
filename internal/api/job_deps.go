package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type depsWatcher interface {
	SetDependencies(job string, deps []watcher.Dependency) error
	GetDependencies(job string) ([]watcher.Dependency, bool)
	CheckDependencies(job string) []string
}

type depPayload struct {
	RequiredJob  string `json:"required_job"`
	MaxStaleness string `json:"max_staleness"`
}

func makeHandleJobDeps(w depsWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPut {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			deps, ok := w.GetDependencies(name)
			if !ok {
				http.Error(rw, "job not found", http.StatusNotFound)
				return
			}
			unsatisfied := w.CheckDependencies(name)
			rw.Header().Set("Content-Type", "application/json")
			json.NewEncoder(rw).Encode(map[string]any{
				"dependencies": deps,
				"unsatisfied":  unsatisfied,
			})

		case http.MethodPut:
			var payload []depPayload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(rw, "invalid body", http.StatusBadRequest)
				return
			}
			deps := make([]watcher.Dependency, 0, len(payload))
			for _, p := range payload {
				d, err := time.ParseDuration(p.MaxStaleness)
				if err != nil {
					http.Error(rw, "invalid max_staleness: "+p.MaxStaleness, http.StatusBadRequest)
					return
				}
				deps = append(deps, watcher.Dependency{RequiredJob: p.RequiredJob, MaxStaleness: d})
			}
			if err := w.SetDependencies(name, deps); err != nil {
				http.Error(rw, "job not found", http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
		}
	}
}
