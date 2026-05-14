package watcher

import (
	"sync"
	"testing"
	"time"

	"github.com/user/cronwatcher/internal/config"
)

type mockSink struct {
	mu      sync.Mutex
	missed  []string
	long    []string
}

func (m *mockSink) MissedJob(job config.Job, _ time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.missed = append(m.missed, job.Name)
}

func (m *mockSink) LongRunningJob(job config.Job, _ time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.long = append(m.long, job.Name)
}

func testConfig() *config.Config {
	return &config.Config{
		CheckInterval: 100 * time.Millisecond,
		Jobs: []config.Job{
			{
				Name:        "backup",
				MaxDelay:    5 * time.Minute,
				MaxDuration: 10 * time.Minute,
			},
		},
	}
}

func TestRecordStart_SetsRunning(t *testing.T) {
	sink := &mockSink{}
	w := New(testConfig(), sink)
	w.RecordStart("backup", time.Now())

	s := w.states["backup"]
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if !s.Running {
		t.Error("expected job to be marked as running")
	}
}

func TestRecordFinish_ClearsRunning(t *testing.T) {
	sink := &mockSink{}
	w := New(testConfig(), sink)
	w.RecordStart("backup", time.Now())
	w.RecordFinish("backup", time.Now())

	s := w.states["backup"]
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if s.Running {
		t.Error("expected job to be marked as not running")
	}
}

func TestCheck_MissedJob(t *testing.T) {
	sink := &mockSink{}
	w := New(testConfig(), sink)

	past := time.Now().Add(-10 * time.Minute)
	w.RecordStart("backup", past)
	w.RecordFinish("backup", past)

	w.check(time.Now())

	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.missed) == 0 {
		t.Error("expected missed job alert")
	}
}

func TestCheck_LongRunningJob(t *testing.T) {
	sink := &mockSink{}
	w := New(testConfig(), sink)

	past := time.Now().Add(-15 * time.Minute)
	w.RecordStart("backup", past)

	w.check(time.Now())

	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.long) == 0 {
		t.Error("expected long-running job alert")
	}
}

func TestCheck_NoAlertWhenHealthy(t *testing.T) {
	sink := &mockSink{}
	w := New(testConfig(), sink)

	now := time.Now()
	w.RecordStart("backup", now.Add(-1*time.Minute))
	w.RecordFinish("backup", now.Add(-30*time.Second))

	w.check(now)

	sink.mu.Lock()
	defer sink.mu.Unlock()
	if len(sink.missed) != 0 || len(sink.long) != 0 {
		t.Error("expected no alerts for healthy job")
	}
}
