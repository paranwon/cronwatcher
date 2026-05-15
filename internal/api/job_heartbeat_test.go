package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatcher/internal/watcher"
)

type fakeHeartbeatWatcher struct {
	*watcher.Watcher
	recordErr error
	lastJob   string
}

func (f *fakeHeartbeatWatcher) RecordHeartbeat(job string) error {
	f.lastJob = job
	return f.recordErr
}

func newHeartbeatWatcher(t *testing.T) *fakeHeartbeatWatcher {
	t.Helper()
	return &fakeHeartbeatWatcher{}
}

func TestHandleJobHeartbeat_RecordsHeartbeat(t *testing.T) {
	fw := newHeartbeatWatcher(t)
	h := makeHandleJobHeartbeat(fw)

	req := httptest.NewRequest(http.MethodPost, "/jobs/heartbeat?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	if fw.lastJob != "backup" {
		t.Fatalf("expected job 'backup', got %q", fw.lastJob)
	}
}

func TestHandleJobHeartbeat_MissingName_ReturnsBadRequest(t *testing.T) {
	fw := newHeartbeatWatcher(t)
	h := makeHandleJobHeartbeat(fw)

	req := httptest.NewRequest(http.MethodPost, "/jobs/heartbeat", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandleJobHeartbeat_UnknownJob_Returns404(t *testing.T) {
	fw := newHeartbeatWatcher(t)
	fw.recordErr = errors.New("unknown job")
	h := makeHandleJobHeartbeat(fw)

	req := httptest.NewRequest(http.MethodPost, "/jobs/heartbeat?name=ghost", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rw.Code)
	}
}

func TestHandleJobHeartbeat_MethodNotAllowed(t *testing.T) {
	fw := newHeartbeatWatcher(t)
	h := makeHandleJobHeartbeat(fw)

	req := httptest.NewRequest(http.MethodGet, "/jobs/heartbeat?name=backup", nil)
	rw := httptest.NewRecorder()
	h(rw, req)

	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}

var _ = time.Now // suppress unused import
