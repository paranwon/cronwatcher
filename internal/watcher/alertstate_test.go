package watcher

import (
	"testing"
)

const (
	alertMissed      = "missed"
	alertLongRunning = "long_running"
)

func TestAlertState_HasFired_FalseInitially(t *testing.T) {
	s := newAlertState([]string{"backup"})
	if s.HasFired("backup", alertMissed) {
		t.Error("expected HasFired to return false initially")
	}
}

func TestAlertState_SetFired_MarksAsFired(t *testing.T) {
	s := newAlertState([]string{"backup"})
	s.SetFired("backup", alertMissed)
	if !s.HasFired("backup", alertMissed) {
		t.Error("expected HasFired to return true after SetFired")
	}
}

func TestAlertState_Clear_ResetsState(t *testing.T) {
	s := newAlertState([]string{"backup"})
	s.SetFired("backup", alertMissed)
	s.Clear("backup", alertMissed)
	if s.HasFired("backup", alertMissed) {
		t.Error("expected HasFired to return false after Clear")
	}
}

func TestAlertState_ClearAll_ResetsAllTypes(t *testing.T) {
	s := newAlertState([]string{"backup"})
	s.SetFired("backup", alertMissed)
	s.SetFired("backup", alertLongRunning)
	s.ClearAll("backup")
	if s.HasFired("backup", alertMissed) || s.HasFired("backup", alertLongRunning) {
		t.Error("expected all states cleared after ClearAll")
	}
}

func TestAlertState_UnknownJob_ReturnsFalse(t *testing.T) {
	s := newAlertState([]string{})
	if s.HasFired("unknown", alertMissed) {
		t.Error("expected HasFired to return false for unknown job")
	}
}

func TestAlertState_SetFired_UnknownJob_CreatesEntry(t *testing.T) {
	s := newAlertState([]string{})
	s.SetFired("newjob", alertMissed)
	if !s.HasFired("newjob", alertMissed) {
		t.Error("expected HasFired to return true after SetFired on unknown job")
	}
}

func TestAlertState_IndependentPerType(t *testing.T) {
	s := newAlertState([]string{"backup"})
	s.SetFired("backup", alertMissed)
	if s.HasFired("backup", alertLongRunning) {
		t.Error("expected long_running alert to be independent of missed alert")
	}
}
