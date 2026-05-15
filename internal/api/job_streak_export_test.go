package api

import (
	"net/http"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

// MakeHandleJobStreakExported exposes makeHandleJobStreak for external tests.
func MakeHandleJobStreakExported(w interface {
	GetStreak(name string) (watcher.StreakEntry, bool)
}) http.HandlerFunc {
	return makeHandleJobStreak(w)
}
