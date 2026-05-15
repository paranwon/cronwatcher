package watcher_test

import (
	"testing"
	"time"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestSetDependencies_StoresDeps(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	deps := []watcher.Dependency{
		{RequiredJob: "job-b", MaxStaleness: time.Hour},
	}
	if err := w.SetDependencies("job-a", deps); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := w.GetDependencies("job-a")
	if !ok {
		t.Fatal("expected deps to be found")
	}
	if len(got) != 1 || got[0].RequiredJob != "job-b" {
		t.Errorf("unexpected deps: %+v", got)
	}
}

func TestSetDependencies_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	err := w.SetDependencies("no-such-job", nil)
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetDependencies_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_, ok := w.GetDependencies("no-such-job")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetDependencies_ReturnsCopy(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	deps := []watcher.Dependency{{RequiredJob: "job-b", MaxStaleness: time.Hour}}
	_ = w.SetDependencies("job-a", deps)

	got, _ := w.GetDependencies("job-a")
	got[0].RequiredJob = "mutated"

	again, _ := w.GetDependencies("job-a")
	if again[0].RequiredJob == "mutated" {
		t.Error("GetDependencies should return a copy, not a reference")
	}
}

func TestCheckDependencies_SatisfiedWhenRecentSuccess(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	w.RecordStart("job-b")
	w.RecordFinish("job-b", nil)

	_ = w.SetDependencies("job-a", []watcher.Dependency{
		{RequiredJob: "job-b", MaxStaleness: time.Hour},
	})
	unsatisfied := w.CheckDependencies("job-a")
	if len(unsatisfied) != 0 {
		t.Errorf("expected no unsatisfied deps, got %v", unsatisfied)
	}
}

func TestCheckDependencies_UnsatisfiedWhenNeverRun(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.SetDependencies("job-a", []watcher.Dependency{
		{RequiredJob: "job-b", MaxStaleness: time.Hour},
	})
	unsatisfied := w.CheckDependencies("job-a")
	if len(unsatisfied) != 1 || unsatisfied[0] != "job-b" {
		t.Errorf("expected job-b unsatisfied, got %v", unsatisfied)
	}
}
