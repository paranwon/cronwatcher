package watcher

import (
	"strings"
	"testing"
)

func TestRecordOutput_StoresEntry(t *testing.T) {
	s := newOutputStore()
	s.initJob("backup")

	if err := s.recordOutput("backup", "done", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := s.getOutput("backup")
	if !ok {
		t.Fatal("expected entry, got none")
	}
	if e.Stdout != "done" {
		t.Errorf("expected stdout 'done', got %q", e.Stdout)
	}
	if e.Truncated {
		t.Error("expected not truncated")
	}
}

func TestRecordOutput_UnknownJob_ReturnsError(t *testing.T) {
	s := newOutputStore()
	err := s.recordOutput("ghost", "out", "err")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetOutput_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newOutputStore()
	_, ok := s.getOutput("missing")
	if ok {
		t.Error("expected false for unknown job")
	}
}

func TestRecordOutput_TruncatesLargeStdout(t *testing.T) {
	s := newOutputStore()
	s.initJob("bigout")

	large := strings.Repeat("x", maxOutputBytes+100)
	if err := s.recordOutput("bigout", large, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, ok := s.getOutput("bigout")
	if !ok {
		t.Fatal("expected entry")
	}
	if len(e.Stdout) != maxOutputBytes {
		t.Errorf("expected stdout truncated to %d, got %d", maxOutputBytes, len(e.Stdout))
	}
	if !e.Truncated {
		t.Error("expected truncated flag to be true")
	}
}

func TestRecordOutput_TruncatesLargeStderr(t *testing.T) {
	s := newOutputStore()
	s.initJob("bigerr")

	large := strings.Repeat("e", maxOutputBytes+50)
	if err := s.recordOutput("bigerr", "", large); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	e, _ := s.getOutput("bigerr")
	if len(e.Stderr) != maxOutputBytes {
		t.Errorf("expected stderr truncated to %d, got %d", maxOutputBytes, len(e.Stderr))
	}
	if !e.Truncated {
		t.Error("expected truncated flag")
	}
}

func TestRecordOutput_OverwritesPrevious(t *testing.T) {
	s := newOutputStore()
	s.initJob("job1")

	_ = s.recordOutput("job1", "first", "")
	_ = s.recordOutput("job1", "second", "")

	e, _ := s.getOutput("job1")
	if e.Stdout != "second" {
		t.Errorf("expected 'second', got %q", e.Stdout)
	}
}
