package watcher_test

import (
	"testing"

	"github.com/example/cronwatcher/internal/watcher"
)

func TestSetTags_StoresTags(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	err := w.SetTags("backup", map[string]string{"team": "infra", "env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tags, ok := w.GetTags("backup")
	if !ok {
		t.Fatal("expected tags to be found")
	}
	if tags["team"] != "infra" {
		t.Errorf("expected team=infra, got %q", tags["team"])
	}
	if tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %q", tags["env"])
	}
}

func TestSetTags_MergesWithExistingTags(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	_ = w.SetTags("backup", map[string]string{"team": "infra"})
	_ = w.SetTags("backup", map[string]string{"env": "staging"})

	tags, _ := w.GetTags("backup")
	if tags["team"] != "infra" {
		t.Errorf("expected team=infra, got %q", tags["team"])
	}
	if tags["env"] != "staging" {
		t.Errorf("expected env=staging, got %q", tags["env"])
	}
}

func TestSetTags_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	err := w.SetTags("nonexistent", map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error for unknown job, got nil")
	}
}

func TestGetTags_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest(testConfig())

	_, ok := w.GetTags("nonexistent")
	if ok {
		t.Fatal("expected ok=false for unknown job")
	}
}

func TestGetTags_ReturnsCopy(t *testing.T) {
	w := watcher.NewForTest(testConfig())
	_ = w.SetTags("backup", map[string]string{"team": "infra"})

	tags, _ := w.GetTags("backup")
	tags["team"] = "mutated"

	original, _ := w.GetTags("backup")
	if original["team"] != "infra" {
		t.Errorf("expected original tag to be unmodified, got %q", original["team"])
	}
}
