package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/api"
	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type fakeDepsWatcher struct {
	deps        map[string][]watcher.Dependency
	unsatisfied []string
}

func (f *fakeDepsWatcher) SetDependencies(job string, deps []watcher.Dependency) error {
	if _, ok := f.deps[job]; !ok {
		return watcher.ErrUnknownJob
	}
	f.deps[job] = deps
	return nil
}
func (f *fakeDepsWatcher) GetDependencies(job string) ([]watcher.Dependency, bool) {
	d, ok := f.deps[job]
	return d, ok
}
func (f *fakeDepsWatcher) CheckDependencies(job string) []string { return f.unsatisfied }

func newFakeDepsWatcher() *fakeDepsWatcher {
	return &fakeDepsWatcher{
		deps: map[string][]watcher.Dependency{
			"job-a": {{RequiredJob: "job-b", MaxStaleness: time.Hour}},
		},
	}
}

func TestHandleJobDeps_GET_ReturnsDeps(t *testing.T) {
	fw := newFakeDepsWatcher()
	h := api.MakeHandleJobDepsExported(fw)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/jobs/deps?name=job-a", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var body map[string]any
	json.NewDecoder(rr.Body).Decode(&body)
	if _, ok := body["dependencies"]; !ok {
		t.Error("expected dependencies key in response")
	}
}

func TestHandleJobDeps_GET_NotFound(t *testing.T) {
	fw := newFakeDepsWatcher()
	h := api.MakeHandleJobDepsExported(fw)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/jobs/deps?name=missing", nil))
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestHandleJobDeps_PUT_SetsDeps(t *testing.T) {
	fw := newFakeDepsWatcher()
	h := api.MakeHandleJobDepsExported(fw)
	body, _ := json.Marshal([]map[string]string{{"required_job": "job-c", "max_staleness": "2h"}})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/jobs/deps?name=job-a", bytes.NewReader(body)))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestHandleJobDeps_PUT_InvalidDuration(t *testing.T) {
	fw := newFakeDepsWatcher()
	h := api.MakeHandleJobDepsExported(fw)
	body, _ := json.Marshal([]map[string]string{{"required_job": "job-c", "max_staleness": "bad"}})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/jobs/deps?name=job-a", bytes.NewReader(body)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHandleJobDeps_MethodNotAllowed(t *testing.T) {
	fw := newFakeDepsWatcher()
	h := api.MakeHandleJobDepsExported(fw)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/jobs/deps?name=job-a", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
