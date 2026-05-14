package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestWebhookServer(t *testing.T, statusCode int, capturedBody *map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("failed to decode webhook body: %v", err)
		}
		*capturedBody = payload
		w.WriteHeader(statusCode)
	}))
}

func TestWebhook_MissedJob_SendsCorrectPayload(t *testing.T) {
	var captured map[string]interface{}
	server := newTestWebhookServer(t, http.StatusOK, &captured)
	defer server.Close()

	w := NewWebhook(server.URL)
	expectedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

	if err := w.MissedJob("backup", expectedAt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured["event"] != "missed_job" {
		t.Errorf("expected event 'missed_job', got %q", captured["event"])
	}
	if captured["job"] != "backup" {
		t.Errorf("expected job 'backup', got %q", captured["job"])
	}
	if captured["at"] != expectedAt.Format(time.RFC3339) {
		t.Errorf("expected at %q, got %q", expectedAt.Format(time.RFC3339), captured["at"])
	}
}

func TestWebhook_LongRunningJob_SendsCorrectPayload(t *testing.T) {
	var captured map[string]interface{}
	server := newTestWebhookServer(t, http.StatusOK, &captured)
	defer server.Close()

	w := NewWebhook(server.URL)

	if err := w.LongRunningJob("report", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if captured["event"] != "long_running_job" {
		t.Errorf("expected event 'long_running_job', got %q", captured["event"])
	}
	if captured["job"] != "report" {
		t.Errorf("expected job 'report', got %q", captured["job"])
	}
}

func TestWebhook_NonOKStatus_ReturnsError(t *testing.T) {
	var captured map[string]interface{}
	server := newTestWebhookServer(t, http.StatusInternalServerError, &captured)
	defer server.Close()

	w := NewWebhook(server.URL)

	err := w.MissedJob("cleanup", time.Now())
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}
