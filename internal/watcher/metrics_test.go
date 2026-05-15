package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestGetMetrics_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_, ok := w.GetMetrics("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetMetrics_AfterFinish_ReturnsStats(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	w.RecordStart("backup", "run-1")
	w.RecordFinish("backup", "run-1", nil)

	m, ok := w.GetMetrics("backup")
	if !ok {
		t.Fatal("expected metrics to exist after finish")
	}
	if m.RunCount != 1 {
		t.Errorf("expected RunCount=1, got %d", m.RunCount)
	}
	if m.ErrorCount != 0 {
		t.Errorf("expected ErrorCount=0, got %d", m.ErrorCount)
	}
}

func TestGetMetrics_RecordsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	w.RecordStart("backup", "run-2")
	w.RecordFinish("backup", "run-2", &jobError{"boom"})

	m, _ := w.GetMetrics("backup")
	if m.ErrorCount != 1 {
		t.Errorf("expected ErrorCount=1, got %d", m.ErrorCount)
	}
}

func TestGetMetrics_MinMaxAvg(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	for i, id := range []string{"r1", "r2", "r3"} {
		w.RecordStart("backup", id)
		_ = i
		w.RecordFinish("backup", id, nil)
	}

	m, _ := w.GetMetrics("backup")
	if m.RunCount != 3 {
		t.Errorf("expected RunCount=3, got %d", m.RunCount)
	}
	if m.MinDuration > m.MaxDuration {
		t.Error("MinDuration should be <= MaxDuration")
	}
	if m.AvgDuration() == 0 && m.TotalTime > 0 {
		t.Error("AvgDuration should not be zero when TotalTime > 0")
	}
}

func TestAllMetrics_ReturnsAllJobs(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	w.RecordStart("backup", "r1")
	w.RecordFinish("backup", "r1", nil)

	all := w.AllMetrics()
	if len(all) == 0 {
		t.Fatal("expected at least one metrics entry")
	}
}

type jobError struct{ msg string }

func (e *jobError) Error() string { return e.msg }

func TestGetMetrics_LastRun_IsRecent(t *testing.T) {
	before := time.Now()
	w := watcher.NewForTest(testConfig())
	w.RecordStart("backup", "r1")
	w.RecordFinish("backup", "r1", nil)

	m, _ := w.GetMetrics("backup")
	if m.LastRun.Before(before) {
		t.Error("LastRun should be after test start")
	}
}
