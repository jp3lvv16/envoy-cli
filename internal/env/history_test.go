package env

import (
	"testing"
)

func makeHistorySet(t *testing.T, name string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	return s
}

func TestNewHistoryEmptyNameReturnsError(t *testing.T) {
	_, err := NewHistory("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewHistoryCreatesEmptyHistory(t *testing.T) {
	h, err := NewHistory("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", h.Len())
	}
}

func TestRecordNilSetReturnsError(t *testing.T) {
	h, _ := NewHistory("prod")
	if err := h.Record(nil, ""); err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestRecordCapturesState(t *testing.T) {
	s := makeHistorySet(t, "prod")
	_ = s.Put("KEY", "v1")
	h, _ := NewHistory("prod")
	_ = h.Record(s, "initial")
	_ = s.Put("KEY", "v2")
	_ = h.Record(s, "updated")

	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
	e0, _ := h.At(0)
	if e0.Snapshot["KEY"] != "v1" {
		t.Errorf("entry 0: expected v1, got %s", e0.Snapshot["KEY"])
	}
	e1, _ := h.At(1)
	if e1.Snapshot["KEY"] != "v2" {
		t.Errorf("entry 1: expected v2, got %s", e1.Snapshot["KEY"])
	}
}

func TestAtOutOfRangeReturnsError(t *testing.T) {
	h, _ := NewHistory("prod")
	_, err := h.At(0)
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestRollbackRestoresState(t *testing.T) {
	s := makeHistorySet(t, "prod")
	_ = s.Put("KEY", "v1")
	h, _ := NewHistory("prod")
	_ = h.Record(s, "initial")
	_ = s.Put("KEY", "v2")

	if err := h.Rollback(s, 0); err != nil {
		t.Fatalf("Rollback error: %v", err)
	}
	v, _ := s.Get("KEY")
	if v != "v1" {
		t.Errorf("expected v1 after rollback, got %s", v)
	}
}

func TestEntriesReturnsCopy(t *testing.T) {
	s := makeHistorySet(t, "prod")
	_ = s.Put("A", "1")
	h, _ := NewHistory("prod")
	_ = h.Record(s, "snap")
	entries := h.Entries()
	entries[0].Label = "mutated"
	orig, _ := h.At(0)
	if orig.Label == "mutated" {
		t.Error("Entries should return a copy, not a reference")
	}
}
