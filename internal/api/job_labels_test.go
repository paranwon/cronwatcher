package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/api"
)

type fakeLabelWatcher struct {
	labels map[string]map[string]string
}

func newFakeLabelWatcher(jobs ...string) *fakeLabelWatcher {
	m := make(map[string]map[string]string)
	for _, j := range jobs {
		m[j] = make(map[string]string)
	}
	return &fakeLabelWatcher{labels: m}
}

func (f *fakeLabelWatcher) SetLabels(job string, labels map[string]string) error {
	if _, ok := f.labels[job]; !ok {
		return errors.New("unknown job: " + job)
	}
	for k, v := range labels {
		f.labels[job][k] = v
	}
	return nil
}

func (f *fakeLabelWatcher) GetLabels(job string) (map[string]string, bool) {
	v, ok := f.labels[job]
	return v, ok
}

func (f *fakeLabelWatcher) DeleteLabel(job, key string) error {
	if _, ok := f.labels[job]; !ok {
		return errors.New("unknown job: " + job)
	}
	delete(f.labels[job], key)
	return nil
}

func TestHandleJobLabels_GET_ReturnsLabels(t *testing.T) {
	w := newFakeLabelWatcher("job1")
	_ = w.SetLabels("job1", map[string]string{"env": "prod"})
	h := api.ExportMakeHandleJobLabels(w)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/labels?name=job1", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result map[string]string
	json.NewDecoder(rec.Body).Decode(&result)
	if result["env"] != "prod" {
		t.Errorf("unexpected labels: %v", result)
	}
}

func TestHandleJobLabels_GET_NotFound(t *testing.T) {
	w := newFakeLabelWatcher()
	h := api.ExportMakeHandleJobLabels(w)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/labels?name=ghost", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandleJobLabels_POST_SetsLabels(t *testing.T) {
	w := newFakeLabelWatcher("job1")
	h := api.ExportMakeHandleJobLabels(w)

	body, _ := json.Marshal(map[string]string{"team": "infra"})
	req := httptest.NewRequest(http.MethodPost, "/api/jobs/labels?name=job1", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if w.labels["job1"]["team"] != "infra" {
		t.Error("expected label to be set")
	}
}

func TestHandleJobLabels_DELETE_RemovesKey(t *testing.T) {
	w := newFakeLabelWatcher("job1")
	_ = w.SetLabels("job1", map[string]string{"env": "prod"})
	h := api.ExportMakeHandleJobLabels(w)

	req := httptest.NewRequest(http.MethodDelete, "/api/jobs/labels?name=job1&key=env", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if _, exists := w.labels["job1"]["env"]; exists {
		t.Error("expected label to be deleted")
	}
}

func TestHandleJobLabels_MethodNotAllowed(t *testing.T) {
	w := newFakeLabelWatcher("job1")
	h := api.ExportMakeHandleJobLabels(w)

	req := httptest.NewRequest(http.MethodPut, "/api/jobs/labels?name=job1", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
