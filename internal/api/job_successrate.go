package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type successRateWatcher interface {
	GetSuccessRate(name string) (watcher.SuccessRate, bool)
}

func makeHandleJobSuccessRate(w successRateWatcher) http.HandlerFunc {
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
		rate, ok := w.GetSuccessRate(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(rate)
	}
}
