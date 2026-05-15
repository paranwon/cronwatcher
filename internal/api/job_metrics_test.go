package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/api"
	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func newMetricsWatcher() *watcher.Watcher {
	return newTestWatcher()
}

func TestHandleJobMetrics_ReturnsMetrics(t *testing.T) {
	w := newMetricsWatcher()
	w.RecordStart("backup", "run-1")
	w.RecordFinish("backup", "run-1", nil)

	handler := api.MakeHandleJobMetricsExported(w)
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/metrics?name=backup", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp["job_name"] != "backup" {
		t.Errorf("expected job_name=backup, got %v", resp["job_name"])
	}
	if resp["run_count"].(float64) != 1 {
		t.Errorf("expected run_count=1, got %v", resp["run_count"])
	}
}

func TestHandleJobMetrics_MissingName_ReturnsBadRequest(t *testing.T) {
	w := newMetricsWatcher()
	handler := api.MakeHandleJobMetricsExported(w)
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/metrics", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleJobMetrics_NotFound_Returns404(t *testing.T) {
	w := newMetricsWatcher()
	handler := api.MakeHandleJobMetricsExported(w)
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/metrics?name=ghost", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleJobMetrics_MethodNotAllowed(t *testing.T) {
	w := newMetricsWatcher()
	handler := api.MakeHandleJobMetricsExported(w)
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/metrics?name=backup", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestHandleJobMetrics_LastRunPresent(t *testing.T) {
	w := newMetricsWatcher()
	w.RecordStart("backup", "run-x")
	w.RecordFinish("backup", "run-x", nil)

	handler := api.MakeHandleJobMetricsExported(w)
	req := httptest.NewRequest(http.MethodGet, "/api/jobs/metrics?name=backup", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	var resp map[string]interface{}
	json.NewDecoder(rec.Body).Decode(&resp)

	if _, ok := resp["last_run"]; !ok {
		t.Error("expected last_run field in response")
	}
	_ = time.Now() // ensure time package used
}
