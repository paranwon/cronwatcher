package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/api"
	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type fakeStreakWatcher struct {
	streaks map[string]watcher.StreakEntry
}

func (f *fakeStreakWatcher) GetStreak(name string) (watcher.StreakEntry, bool) {
	e, ok := f.streaks[name]
	return e, ok
}

func newStreakWatcher() *fakeStreakWatcher {
	return &fakeStreakWatcher{
		streaks: map[string]watcher.StreakEntry{
			"backup": {SuccessStreak: 3, FailureStreak: 0},
		},
	}
}

func TestHandleJobStreak_ReturnsStreak(t *testing.T) {
	h := api.MakeHandleJobStreakExported(newStreakWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/streak?name=backup", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]int
	_ = json.NewDecoder(rec.Body).Decode(&body)
	if body["success_streak"] != 3 {
		t.Errorf("expected success_streak=3, got %d", body["success_streak"])
	}
}

func TestHandleJobStreak_MissingName_ReturnsBadRequest(t *testing.T) {
	h := api.MakeHandleJobStreakExported(newStreakWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/streak", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobStreak_NotFound_Returns404(t *testing.T) {
	h := api.MakeHandleJobStreakExported(newStreakWatcher())
	req := httptest.NewRequest(http.MethodGet, "/jobs/streak?name=ghost", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleJobStreak_MethodNotAllowed(t *testing.T) {
	h := api.MakeHandleJobStreakExported(newStreakWatcher())
	req := httptest.NewRequest(http.MethodPost, "/jobs/streak?name=backup", nil)
	rec := httptest.NewRecorder()
	h(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
