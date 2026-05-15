package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

type fakeSuccessRateWatcher struct {
	rates map[string]watcher.SuccessRate
}

func (f *fakeSuccessRateWatcher) GetSuccessRate(name string) (watcher.SuccessRate, bool) {
	r, ok := f.rates[name]
	return r, ok
}

func newFakeSuccessRateWatcher() *fakeSuccessRateWatcher {
	return &fakeSuccessRateWatcher{
		rates: map[string]watcher.SuccessRate{
			"backup": {Total: 10, Success: 8, Failure: 2, Rate: 0.8},
		},
	}
}

func TestHandleJobSuccessRate_ReturnsRate(t *testing.T) {
	h := makeHandleJobSuccessRate(newFakeSuccessRateWatcher())
	req := httptest.NewRequest(http.MethodGet, "/job/successrate?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var got watcher.SuccessRate
	if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if got.Total != 10 || got.Success != 8 || got.Failure != 2 {
		t.Fatalf("unexpected rate data: %+v", got)
	}
}

func TestHandleJobSuccessRate_MissingName_ReturnsBadRequest(t *testing.T) {
	h := makeHandleJobSuccessRate(newFakeSuccessRateWatcher())
	req := httptest.NewRequest(http.MethodGet, "/job/successrate", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleJobSuccessRate_NotFound_Returns404(t *testing.T) {
	h := makeHandleJobSuccessRate(newFakeSuccessRateWatcher())
	req := httptest.NewRequest(http.MethodGet, "/job/successrate?name=ghost", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
}

func TestHandleJobSuccessRate_MethodNotAllowed(t *testing.T) {
	h := makeHandleJobSuccessRate(newFakeSuccessRateWatcher())
	req := httptest.NewRequest(http.MethodPost, "/job/successrate?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}
