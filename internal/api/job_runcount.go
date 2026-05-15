package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type runCountWatcher interface {
	GetRunCount(name string) (watcher.RunCount, bool)
}

func makeHandleJobRunCount(w runCountWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(rw, "missing name", http.StatusBadRequest)
			return
		}

		count, ok := w.GetRunCount(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"job":     name,
			"success": count.Success,
			"failure": count.Failure,
			"total":   count.Total,
		})
	}
}
