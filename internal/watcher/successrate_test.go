package watcher

import (
	"testing"
)

func TestGetSuccessRate_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newSuccessRateStore([]string{"backup"})
	_, ok := s.get("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestRecordSuccess_UnknownJob_ReturnsError(t *testing.T) {
	s := newSuccessRateStore([]string{"backup"})
	if err := s.recordSuccess("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordFailure_UnknownJob_ReturnsError(t *testing.T) {
	s := newSuccessRateStore([]string{"backup"})
	if err := s.recordFailure("ghost"); err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestRecordSuccess_IncrementsCorrectly(t *testing.T) {
	s := newSuccessRateStore([]string{"backup"})
	_ = s.recordSuccess("backup")
	_ = s.recordSuccess("backup")
	r, ok := s.get("backup")
	if !ok {
		t.Fatal("expected result")
	}
	if r.Total != 2 || r.Success != 2 || r.Failure != 0 {
		t.Fatalf("unexpected counts: %+v", r)
	}
	if r.Rate != 1.0 {
		t.Fatalf("expected rate 1.0, got %f", r.Rate)
	}
}

func TestRecordFailure_IncrementsCorrectly(t *testing.T) {
	s := newSuccessRateStore([]string{"backup"})
	_ = s.recordSuccess("backup")
	_ = s.recordFailure("backup")
	_ = s.recordFailure("backup")
	r, _ := s.get("backup")
	if r.Total != 3 || r.Success != 1 || r.Failure != 2 {
		t.Fatalf("unexpected counts: %+v", r)
	}
	const want = 1.0 / 3.0
	if r.Rate < want-0.0001 || r.Rate > want+0.0001 {
		t.Fatalf("expected rate ~0.333, got %f", r.Rate)
	}
}

func TestGetSuccessRate_InitiallyZero(t *testing.T) {
	s := newSuccessRateStore([]string{"backup"})
	r, ok := s.get("backup")
	if !ok {
		t.Fatal("expected result for known job")
	}
	if r.Total != 0 || r.Rate != 0 {
		t.Fatalf("expected zeroed struct, got %+v", r)
	}
}
