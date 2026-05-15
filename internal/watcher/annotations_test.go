package watcher_test

import (
	"testing"

	"github.com/cronwatcher/cronwatcher/internal/watcher"
)

func TestSetAnnotations_StoresAnnotations(t *testing.T) {
	w := watcher.NewForTest([]string{"backup"})

	err := w.SetAnnotations("backup", map[string]string{"owner": "ops", "env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	anns, ok := w.GetAnnotations("backup")
	if !ok {
		t.Fatal("expected annotations to exist")
	}
	if anns["owner"] != "ops" || anns["env"] != "prod" {
		t.Errorf("unexpected annotations: %v", anns)
	}
}

func TestSetAnnotations_MergesWithExisting(t *testing.T) {
	w := watcher.NewForTest([]string{"backup"})

	_ = w.SetAnnotations("backup", map[string]string{"owner": "ops"})
	_ = w.SetAnnotations("backup", map[string]string{"env": "prod"})

	anns, _ := w.GetAnnotations("backup")
	if anns["owner"] != "ops" || anns["env"] != "prod" {
		t.Errorf("expected merged annotations, got: %v", anns)
	}
}

func TestSetAnnotations_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest([]string{})

	err := w.SetAnnotations("ghost", map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}

func TestGetAnnotations_UnknownJob_ReturnsFalse(t *testing.T) {
	w := watcher.NewForTest([]string{})

	_, ok := w.GetAnnotations("ghost")
	if ok {
		t.Fatal("expected false for unknown job")
	}
}

func TestGetAnnotations_ReturnsCopy(t *testing.T) {
	w := watcher.NewForTest([]string{"backup"})
	_ = w.SetAnnotations("backup", map[string]string{"owner": "ops"})

	anns, _ := w.GetAnnotations("backup")
	anns["owner"] = "mutated"

	anns2, _ := w.GetAnnotations("backup")
	if anns2["owner"] != "ops" {
		t.Errorf("expected original value, got %q", anns2["owner"])
	}
}

func TestDeleteAnnotation_RemovesKey(t *testing.T) {
	w := watcher.NewForTest([]string{"backup"})
	_ = w.SetAnnotations("backup", map[string]string{"owner": "ops", "env": "prod"})

	err := w.DeleteAnnotation("backup", "owner")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	anns, _ := w.GetAnnotations("backup")
	if _, exists := anns["owner"]; exists {
		t.Error("expected 'owner' annotation to be deleted")
	}
	if anns["env"] != "prod" {
		t.Errorf("expected 'env' to remain, got: %v", anns)
	}
}

func TestDeleteAnnotation_UnknownJob_ReturnsError(t *testing.T) {
	w := watcher.NewForTest([]string{})

	err := w.DeleteAnnotation("ghost", "key")
	if err == nil {
		t.Fatal("expected error for unknown job")
	}
}
