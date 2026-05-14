package alert

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/user/cronwatcher/internal/config"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	h := slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	return slog.New(h)
}

func TestLogger_MissedJob(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(newTestLogger(&buf))

	job := config.Job{
		Name:     "cleanup",
		MaxDelay: 5 * time.Minute,
	}
	l.MissedJob(job, time.Now().Add(-10*time.Minute))

	out := buf.String()
	if !strings.Contains(out, "missed job") {
		t.Errorf("expected 'missed job' in output, got: %s", out)
	}
	if !strings.Contains(out, "cleanup") {
		t.Errorf("expected job name in output, got: %s", out)
	}
}

func TestLogger_LongRunningJob(t *testing.T) {
	var buf bytes.Buffer
	l := NewLogger(newTestLogger(&buf))

	job := config.Job{
		Name:        "report",
		MaxDuration: 10 * time.Minute,
	}
	l.LongRunningJob(job, 15*time.Minute)

	out := buf.String()
	if !strings.Contains(out, "long-running job") {
		t.Errorf("expected 'long-running job' in output, got: %s", out)
	}
	if !strings.Contains(out, "report") {
		t.Errorf("expected job name in output, got: %s", out)
	}
}
