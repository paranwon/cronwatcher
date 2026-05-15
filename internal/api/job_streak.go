package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type streakReader interface {
	GetStreak(name string) (watcher.StreakEntry, bool)
}

func makeHandleJobStreak(w streakReader) http.HandlerFunc {
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
		entry, ok := w.GetStreak(name)
		if !ok {
			http.Error(rw, "job not found", http.StatusNotFound)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(rw).Encode(map[string]int{
			"success_streak": entry.SuccessStreak,
			"failure_streak": entry.FailureStreak,
		})
	}
}
