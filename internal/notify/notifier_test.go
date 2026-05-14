package notify

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cronwatcher/internal/watcher"
)

type mockAlerter struct {
	missedCalls      []string
	longRunningCalls []string
	returnErr        error
}

func (m *mockAlerter) MissedJob(name string, lastRun time.Time) error {
	m.missedCalls = append(m.missedCalls, name)
	return m.returnErr
}

func (m *mockAlerter) LongRunningJob(name string, duration time.Duration) error {
	m.longRunningCalls = append(m.longRunningCalls, name)
	return m.returnErr
}

func newTestNotifier(w *watcher.Watcher, a *mockAlerter) *Notifier {
	return New(w, a, log.New(os.Stdout, "", 0))
}

func TestCheckMissed_AlertsOnMissedJob(t *testing.T) {
	w := watcher.NewForTest()
	w.InjectMissed("backup", time.Now().Add(-2*time.Hour))

	m := &mockAlerter{}
	n := newTestNotifier(w, m)
	n.CheckMissed()

	if len(m.missedCalls) != 1 || m.missedCalls[0] != "backup" {
		t.Errorf("expected missed alert for 'backup', got %v", m.missedCalls)
	}
}

func TestCheckMissed_AlertError_DoesNotPanic(t *testing.T) {
	w := watcher.NewForTest()
	w.InjectMissed("backup", time.Now().Add(-1*time.Hour))

	m := &mockAlerter{returnErr: errors.New("send failed")}
	n := newTestNotifier(w, m)

	// should not panic
	n.CheckMissed()
}

func TestCheckLongRunning_AlertsOnExceedingJob(t *testing.T) {
	w := watcher.NewForTest()
	w.InjectLongRunning("report", 10*time.Minute)

	m := &mockAlerter{}
	n := newTestNotifier(w, m)
	n.CheckLongRunning()

	if len(m.longRunningCalls) != 1 || m.longRunningCalls[0] != "report" {
		t.Errorf("expected long-running alert for 'report', got %v", m.longRunningCalls)
	}
}

func TestCheckLongRunning_NoExceedingJobs_NoAlert(t *testing.T) {
	w := watcher.NewForTest()

	m := &mockAlerter{}
	n := newTestNotifier(w, m)
	n.CheckLongRunning()

	if len(m.longRunningCalls) != 0 {
		t.Errorf("expected no alerts, got %v", m.longRunningCalls)
	}
}
