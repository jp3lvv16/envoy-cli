package env

import (
	"testing"
)

func makeRollbackHistory(t *testing.T) (*History, *Set) {
	t.Helper()
	s, err := NewSet("prod")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	h, err := NewHistory("prod")
	if err != nil {
		t.Fatalf("NewHistory: %v", err)
	}
	return h, s
}

func TestRollbackNilHistoryReturnsError(t *testing.T) {
	_, err := Rollback(nil, 0)
	if err == nil {
		t.Fatal("expected error for nil history")
	}
}

func TestRollbackNegativeVersionReturnsError(t *testing.T) {
	h, s := makeRollbackHistory(t)
	_ = s.Put("K", "v")
	if err := h.Record(s); err != nil {
		t.Fatalf("Record: %v", err)
	}
	_, err := Rollback(h, -1)
	if err == nil {
		t.Fatal("expected error for negative version")
	}
}

func TestRollbackNoSnapshotsReturnsError(t *testing.T) {
	h, _ := makeRollbackHistory(t)
	_, err := Rollback(h, 0)
	if err == nil {
		t.Fatal("expected error when no snapshots exist")
	}
}

func TestRollbackVersionOutOfRangeReturnsError(t *testing.T) {
	h, s := makeRollbackHistory(t)
	_ = s.Put("A", "1")
	if err := h.Record(s); err != nil {
		t.Fatalf("Record: %v", err)
	}
	_, err := Rollback(h, 5)
	if err == nil {
		t.Fatal("expected error for out-of-range version")
	}
}

func TestRollbackReturnsCorrectState(t *testing.T) {
	h, s := makeRollbackHistory(t)

	_ = s.Put("KEY", "first")
	if err := h.Record(s); err != nil {
		t.Fatalf("Record v0: %v", err)
	}

	_ = s.Put("KEY", "second")
	if err := h.Record(s); err != nil {
		t.Fatalf("Record v1: %v", err)
	}

	// version 0 → most recent snapshot ("second")
	res, err := Rollback(h, 0)
	if err != nil {
		t.Fatalf("Rollback v0: %v", err)
	}
	val, _ := res.Set.Get("KEY")
	if val != "second" {
		t.Errorf("expected 'second', got %q", val)
	}

	// version 1 → older snapshot ("first")
	res, err = Rollback(h, 1)
	if err != nil {
		t.Fatalf("Rollback v1: %v", err)
	}
	val, _ = res.Set.Get("KEY")
	if val != "first" {
		t.Errorf("expected 'first', got %q", val)
	}
}

func TestRollbackToLatestEqualsVersionZero(t *testing.T) {
	h, s := makeRollbackHistory(t)
	_ = s.Put("X", "latest")
	if err := h.Record(s); err != nil {
		t.Fatalf("Record: %v", err)
	}
	res, err := RollbackToLatest(h)
	if err != nil {
		t.Fatalf("RollbackToLatest: %v", err)
	}
	if res.Version != 0 {
		t.Errorf("expected version 0, got %d", res.Version)
	}
	val, _ := res.Set.Get("X")
	if val != "latest" {
		t.Errorf("expected 'latest', got %q", val)
	}
}
