package api

import (
	"encoding/json"
	"net/http"
	"time"
)

type retentionWatcher interface {
	PruneHistory(maxAge time.Duration) int
}

// makeHandleJobRetention returns an HTTP handler that triggers an immediate
// history prune with the supplied max-age query parameter (e.g. ?max_age=48h).
func makeHandleJobRetention(w retentionWatcher) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rawAge := r.URL.Query().Get("max_age")
		if rawAge == "" {
			http.Error(rw, "max_age query parameter required", http.StatusBadRequest)
			return
		}

		maxAge, err := time.ParseDuration(rawAge)
		if err != nil || maxAge <= 0 {
			http.Error(rw, "invalid max_age value", http.StatusBadRequest)
			return
		}

		removed := w.PruneHistory(maxAge)

		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(map[string]int{"removed": removed})
	}
}
