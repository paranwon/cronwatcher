package watcher

import (
	"testing"
)

func TestRunCount_InitiallyZero(t *testing.T) {
	w := NewForTest(testConfig())
	name := testConfig().Jobs[0].Name

	count, ok := w.GetRunCount(name)
	if !ok {
		t.Fatalf("expected run count to exist for %q", name)
	}
	if count.Success != 0 || count.Failure != 0 || count.Total != 0 {
		t.Errorf("expected zero counts, got %+v", count)
	}
}

func TestRunCount_RecordsSuccess(t *testing.T) {
	w := NewForTest(testConfig())
	name := testConfig().Jobs[0].Name

	w.RecordStart(name)
	w.RecordFinish(name, nil)

	count, ok := w.GetRunCount(name)
	if !ok {
		t.Fatalf("expected run count to exist")
	}
	if count.Success != 1 {
		t.Errorf("expected 1 success, got %d", count.Success)
	}
	if count.Failure != 0 {
		t.Errorf("expected 0 failures, got %d", count.Failure)
	}
	if count.Total != 1 {
		t.Errorf("expected total 1, got %d", count.Total)
	}
}

func TestRunCount_RecordsFailure(t *testing.T) {
	w := NewForTest(testConfig())
	name := testConfig().Jobs[0].Name

	w.RecordStart(name)
	w.RecordFinish(name, fmt.Errorf("job failed"))

	count, ok := w.GetRunCount(name)
	if !ok {
		t.Fatalf("expected run count to exist")
	}
	if count.Failure != 1 {
		t.Errorf("expected 1 failure, got %d", count.Failure)
	}
	if count.Success != 0 {
		t.Errorf("expected 0 successes, got %d", count.Success)
	}
}

func TestRunCount_UnknownJob_ReturnsFalse(t *testing.T) {
	w := NewForTest(testConfig())

	_, ok := w.GetRunCount("nonexistent")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestRunCount_AccumulatesMultipleRuns(t *testing.T) {
	w := NewForTest(testConfig())
	name := testConfig().Jobs[0].Name

	for i := 0; i < 3; i++ {
		w.RecordStart(name)
		w.RecordFinish(name, nil)
	}
	w.RecordStart(name)
	w.RecordFinish(name, fmt.Errorf("err"))

	count, _ := w.GetRunCount(name)
	if count.Success != 3 {
		t.Errorf("expected 3 successes, got %d", count.Success)
	}
	if count.Failure != 1 {
		t.Errorf("expected 1 failure, got %d", count.Failure)
	}
	if count.Total != 4 {
		t.Errorf("expected total 4, got %d", count.Total)
	}
}
