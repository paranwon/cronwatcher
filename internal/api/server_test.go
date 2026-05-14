package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/internal/api"
	"github.com/cronwatcher/internal/config"
	"github.com/cronwatcher/internal/watcher"
)

func newTestWatcher(t *testing.T) *watcher.Watcher {
	t.Helper()
	cfg := &config.Config{
		Jobs: []config.Job{
			{Name: "backup", Schedule: "@daily", MaxDuration: "10m"},
		},
	}
	w, err := watcher.New(cfg, nil)
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}
	return w
}

func TestHandleHealth_ReturnsOK(t *testing.T) {
	w := newTestWatcher(t)
	s := api.New(w, ":0")

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %q", body["status"])
	}
}

func TestHandleStatus_ReturnsJobList(t *testing.T) {
	w := newTestWatcher(t)
	s := api.New(w, ":0")

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var jobs []api.JobStatus
	if err := json.NewDecoder(rec.Body).Decode(&jobs); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(jobs) != 1 {
		t.Errorf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Name != "backup" {
		t.Errorf("expected job name 'backup', got %q", jobs[0].Name)
	}
	if jobs[0].Running {
		t.Errorf("expected job not running")
	}
}
