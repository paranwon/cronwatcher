package watcher_test

import (
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

const (
	testAlertMissed      = "missed"
	testAlertLongRunning = "long_running"
)

func TestWatcher_AlertState_SetAndCheck(t *testing.T) {
	w := watcher.NewForTest([]string{"deploy"})
	if w.HasAlertFired("deploy", testAlertMissed) {
		t.Fatal("expected alert not fired initially")
	}
	w.SetAlertFired("deploy", testAlertMissed)
	if !w.HasAlertFired("deploy", testAlertMissed) {
		t.Fatal("expected alert fired after SetAlertFired")
	}
}

func TestWatcher_AlertState_ClearResetsState(t *testing.T) {
	w := watcher.NewForTest([]string{"deploy"})
	w.SetAlertFired("deploy", testAlertMissed)
	w.ClearAlert("deploy", testAlertMissed)
	if w.HasAlertFired("deploy", testAlertMissed) {
		t.Fatal("expected alert cleared after ClearAlert")
	}
}

func TestWatcher_AlertState_ClearAllResetsAllTypes(t *testing.T) {
	w := watcher.NewForTest([]string{"deploy"})
	w.SetAlertFired("deploy", testAlertMissed)
	w.SetAlertFired("deploy", testAlertLongRunning)
	w.ClearAllAlerts("deploy")
	if w.HasAlertFired("deploy", testAlertMissed) {
		t.Fatal("expected missed alert cleared")
	}
	if w.HasAlertFired("deploy", testAlertLongRunning) {
		t.Fatal("expected long_running alert cleared")
	}
}

func TestWatcher_AlertState_IndependentPerJob(t *testing.T) {
	w := watcher.NewForTest([]string{"jobA", "jobB"})
	w.SetAlertFired("jobA", testAlertMissed)
	if w.HasAlertFired("jobB", testAlertMissed) {
		t.Fatal("alert state should be independent per job")
	}
}
