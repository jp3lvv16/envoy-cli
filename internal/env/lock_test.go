package env

import (
	"testing"
)

func TestLockEmptyNameReturnsError(t *testing.T) {
	l := NewLockIndex()
	if err := l.Lock("", "alice"); err == nil {
		t.Fatal("expected error for empty set name")
	}
}

func TestLockEmptyOwnerReturnsError(t *testing.T) {
	l := NewLockIndex()
	if err := l.Lock("prod", ""); err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestLockAndIsLocked(t *testing.T) {
	l := NewLockIndex()
	if err := l.Lock("prod", "alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !l.IsLocked("prod") {
		t.Fatal("expected set to be locked")
	}
}

func TestLockAlreadyLockedReturnsError(t *testing.T) {
	l := NewLockIndex()
	_ = l.Lock("prod", "alice")
	if err := l.Lock("prod", "bob"); err == nil {
		t.Fatal("expected error when locking already locked set")
	}
}

func TestUnlockNotLockedReturnsError(t *testing.T) {
	l := NewLockIndex()
	if err := l.Unlock("prod"); err == nil {
		t.Fatal("expected error when unlocking non-locked set")
	}
}

func TestUnlockRemovesLock(t *testing.T) {
	l := NewLockIndex()
	_ = l.Lock("prod", "alice")
	if err := l.Unlock("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.IsLocked("prod") {
		t.Fatal("expected set to be unlocked")
	}
}

func TestGetLockReturnsEntry(t *testing.T) {
	l := NewLockIndex()
	_ = l.Lock("staging", "bob")
	e, ok := l.GetLock("staging")
	if !ok {
		t.Fatal("expected lock entry to exist")
	}
	if e.LockedBy != "bob" {
		t.Fatalf("expected owner bob, got %s", e.LockedBy)
	}
	if e.LockedAt.IsZero() {
		t.Fatal("expected non-zero lock time")
	}
}

func TestGetLockMissingReturnsNotFound(t *testing.T) {
	l := NewLockIndex()
	_, ok := l.GetLock("nonexistent")
	if ok {
		t.Fatal("expected ok=false for missing lock entry")
	}
}

func TestAllReturnsCopy(t *testing.T) {
	l := NewLockIndex()
	_ = l.Lock("prod", "alice")
	_ = l.Lock("staging", "bob")
	all := l.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the copy must not affect the index.
	delete(all, "prod")
	if !l.IsLocked("prod") {
		t.Fatal("mutating copy should not affect index")
	}
}
