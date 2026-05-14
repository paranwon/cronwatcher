package alert_test

import (
	"testing"

	"github.com/yourorg/cronwatcher/internal/alert"
)

// recordingAlerter captures calls for assertion in tests.
type recordingAlerter struct {
	missed      []string
	longRunning []string
}

func (r *recordingAlerter) MissedJob(jobName string) {
	r.missed = append(r.missed, jobName)
}

func (r *recordingAlerter) LongRunning(jobName string, _ float64) {
	r.longRunning = append(r.longRunning, jobName)
}

func TestMulti_MissedJob_DispatchesToAll(t *testing.T) {
	a1 := &recordingAlerter{}
	a2 := &recordingAlerter{}
	m := alert.NewMulti(a1, a2)

	m.MissedJob("backup")

	for _, a := range []*recordingAlerter{a1, a2} {
		if len(a.missed) != 1 || a.missed[0] != "backup" {
			t.Errorf("expected alerter to receive MissedJob(\"backup\"), got %v", a.missed)
		}
	}
}

func TestMulti_LongRunning_DispatchesToAll(t *testing.T) {
	a1 := &recordingAlerter{}
	a2 := &recordingAlerter{}
	m := alert.NewMulti(a1, a2)

	m.LongRunning("report", 120.5)

	for _, a := range []*recordingAlerter{a1, a2} {
		if len(a.longRunning) != 1 || a.longRunning[0] != "report" {
			t.Errorf("expected alerter to receive LongRunning(\"report\"), got %v", a.longRunning)
		}
	}
}

func TestMulti_NoAlerters_DoesNotPanic(t *testing.T) {
	m := alert.NewMulti()
	m.MissedJob("noop")
	m.LongRunning("noop", 0)
}
