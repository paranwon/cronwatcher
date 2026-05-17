package api

import (
	"encoding/json"
	"net/http"

	"github.com/dkimot/cronwatcher/internal/watcher"
)

type blockedWatcher interface {
	GetBlocked(jobName string) ([]watcher.BlockedEntry, bool)
	ClearBlocked(jobName string) error
}

func makeHandleJobBlocked(w blockedWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodDelete {
			if err := w.ClearBlocked(name); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		entries, ok := w.GetBlocked(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(map[string]interface{}{
			"job":     name,
			"blocked": entries,
		})
	}
}
