package alert_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/yourorg/cronwatcher/internal/alert"
)

func newTestLogger(buf *bytes.Buffer) *alert.Logger {
	l := log.New(buf, "", 0)
	return alert.NewLogger(l)
}

func TestLogger_MissedJob(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	logger.MissedJob("cleanup")

	if !strings.Contains(buf.String(), "cleanup") {
		t.Errorf("expected log output to mention job name, got: %s", buf.String())
	}
	if !strings.Contains(strings.ToLower(buf.String()), "missed") {
		t.Errorf("expected log output to mention 'missed', got: %s", buf.String())
	}
}

func TestLogger_LongRunningJob(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	logger.LongRunning("report", 95.3)

	if !strings.Contains(buf.String(), "report") {
		t.Errorf("expected log output to mention job name, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "95.3") {
		t.Errorf("expected log output to include duration, got: %s", buf.String())
	}
}
