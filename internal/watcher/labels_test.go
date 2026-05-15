package watcher_test

import (
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestSetLabels_StoresLabels(t *testing.T) {
	w := watcher.NewForTest([]string{"job1"})
	err := w.SetLabels("job1", map[string]string{"env": "prod", "team": "infra"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	labels, ok := w.GetLabels("job1")
	if !ok {
		t.Fatal("expected labels to exist")
	}
	if labels["env"] != "prod" || labels["team"] != "infra" {
		t.Errorf("unexpected labels: %v", labels)
	}
}

func TestSetLabels_MergesWithExisting(t *testing.T) {
	w := watcher.NewForTest([]string{"job1"})
	_ = w.SetLabels("job1", map[string]string{"env": "staging"})
	_ = w.SetLabels("job1", map[string]string{"team": "platform"})

	labels, _ := w.GetLabels("job1")
	if labels["env"] != "staging" || labels["team"] != "platform" {
		t.Errorf("expected merged labels, got: %v", labels)
	}
}

func TestSetLabels_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest([]string{})
	err := w.SetLabels("ghost", map[string]string{"x": "y"})
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetLabels_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest([]string{})
	_, ok := w.GetLabels("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetLabels_ReturnsCopy(t *testing.T) {
	w := watcher.NewForTest([]string{"job1"})
	_ = w.SetLabels("job1", map[string]string{"env": "prod"})

	labels, _ := w.GetLabels("job1")
	labels["env"] = "mutated"

	again, _ := w.GetLabels("job1")
	if again["env"] != "prod" {
		t.Errorf("GetLabels returned a reference, not a copy")
	}
}

func TestDeleteLabel_RemovesKey(t *testing.T) {
	w := watcher.NewForTest([]string{"job1"})
	_ = w.SetLabels("job1", map[string]string{"env": "prod", "team": "infra"})
	err := w.DeleteLabel("job1", "env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, _ := w.GetLabels("job1")
	if _, exists := labels["env"]; exists {
		t.Error("expected 'env' label to be deleted")
	}
	if labels["team"] != "infra" {
		t.Error("expected 'team' label to remain")
	}
}
