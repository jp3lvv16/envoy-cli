package env

import (
	"testing"
)

func TestNewSetEmptyNameReturnsError(t *testing.T) {
	_, err := NewSet("")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestNewSetCreatesEmptyVars(t *testing.T) {
	s, err := NewSet("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Vars) != 0 {
		t.Errorf("expected empty vars, got %d entries", len(s.Vars))
	}
}

func TestPutAndGet(t *testing.T) {
	s, _ := NewSet("staging")
	if err := s.Put("APP_ENV", "staging"); err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	v, err := s.Get("APP_ENV")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if v != "staging" {
		t.Errorf("expected %q, got %q", "staging", v)
	}
}

func TestPutEmptyKeyReturnsError(t *testing.T) {
	s, _ := NewSet("dev")
	if err := s.Put("", "value"); err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestGetMissingKeyReturnsError(t *testing.T) {
	s, _ := NewSet("dev")
	_, err := s.Get("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestDelete(t *testing.T) {
	s, _ := NewSet("dev")
	_ = s.Put("KEY", "val")
	if err := s.Delete("KEY"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := s.Get("KEY"); err == nil {
		t.Fatal("expected key to be deleted")
	}
}

func TestDeleteMissingKeyReturnsError(t *testing.T) {
	s, _ := NewSet("dev")
	if err := s.Delete("NOPE"); err == nil {
		t.Fatal("expected error deleting missing key, got nil")
	}
}

func TestList(t *testing.T) {
	s, _ := NewSet("dev")
	_ = s.Put("A", "1")
	_ = s.Put("B", "2")
	lines := s.List()
	if len(lines) != 2 {
		t.Errorf("expected 2 entries, got %d", len(lines))
	}
}
