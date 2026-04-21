package env

import (
	"testing"
	"time"
)

func TestSetExpiryEmptyNameReturnsError(t *testing.T) {
	x := NewExpiryIndex()
	if err := x.Set("", time.Minute); err == nil {
		t.Fatal("expected error for empty set name")
	}
}

func TestSetExpiryZeroTTLReturnsError(t *testing.T) {
	x := NewExpiryIndex()
	if err := x.Set("prod", 0); err == nil {
		t.Fatal("expected error for zero ttl")
	}
}

func TestSetExpiryNegativeTTLReturnsError(t *testing.T) {
	x := NewExpiryIndex()
	if err := x.Set("prod", -time.Second); err == nil {
		t.Fatal("expected error for negative ttl")
	}
}

func TestSetAndGetExpiry(t *testing.T) {
	x := NewExpiryIndex()
	if err := x.Set("prod", time.Hour); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := x.Get("prod")
	if !ok {
		t.Fatal("expected expiry to exist")
	}
	if e.SetName != "prod" {
		t.Errorf("expected set name 'prod', got %q", e.SetName)
	}
	if e.ExpiresAt.Before(time.Now()) {
		t.Error("expected expiry to be in the future")
	}
}

func TestIsExpiredFutureReturnsFalse(t *testing.T) {
	x := NewExpiryIndex()
	_ = x.Set("staging", time.Hour)
	if x.IsExpired("staging") {
		t.Error("expected set not to be expired yet")
	}
}

func TestIsExpiredPastReturnsTrue(t *testing.T) {
	x := NewExpiryIndex()
	_ = x.Set("old", time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !x.IsExpired("old") {
		t.Error("expected set to be expired")
	}
}

func TestIsExpiredMissingSetReturnsFalse(t *testing.T) {
	x := NewExpiryIndex()
	if x.IsExpired("nonexistent") {
		t.Error("expected false for unknown set")
	}
}

func TestRemoveExpiryDeletesEntry(t *testing.T) {
	x := NewExpiryIndex()
	_ = x.Set("dev", time.Minute)
	if err := x.Remove("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := x.Get("dev"); ok {
		t.Error("expected expiry to be removed")
	}
}

func TestRemoveExpiryMissingReturnsError(t *testing.T) {
	x := NewExpiryIndex()
	if err := x.Remove("ghost"); err == nil {
		t.Fatal("expected error removing non-existent entry")
	}
}

func TestAllReturnsAllEntries(t *testing.T) {
	x := NewExpiryIndex()
	_ = x.Set("a", time.Minute)
	_ = x.Set("b", time.Hour)
	all := x.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
