package env

import (
	"testing"
)

func TestFreezeEmptyNameReturnsError(t *testing.T) {
	f := NewFreezeIndex()
	if err := f.Freeze(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestUnfreezeEmptyNameReturnsError(t *testing.T) {
	f := NewFreezeIndex()
	if err := f.Unfreeze(""); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestFreezeAndIsFrozen(t *testing.T) {
	f := NewFreezeIndex()
	if err := f.Freeze("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsFrozen("prod") {
		t.Fatal("expected prod to be frozen")
	}
}

func TestIsFrozenUnknownSetReturnsFalse(t *testing.T) {
	f := NewFreezeIndex()
	if f.IsFrozen("staging") {
		t.Fatal("expected false for unknown set")
	}
}

func TestUnfreezeNotFrozenReturnsError(t *testing.T) {
	f := NewFreezeIndex()
	if err := f.Unfreeze("prod"); err == nil {
		t.Fatal("expected error when unfreezing a non-frozen set")
	}
}

func TestUnfreezeRemovesFrozenStatus(t *testing.T) {
	f := NewFreezeIndex()
	_ = f.Freeze("prod")
	if err := f.Unfreeze("prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.IsFrozen("prod") {
		t.Fatal("expected prod to be unfrozen")
	}
}

func TestFreezeListReturnsAllFrozen(t *testing.T) {
	f := NewFreezeIndex()
	_ = f.Freeze("prod")
	_ = f.Freeze("staging")
	list := f.List()
	if len(list) != 2 {
		t.Fatalf("expected 2 frozen sets, got %d", len(list))
	}
}

func TestAssertMutableFrozenReturnsError(t *testing.T) {
	f := NewFreezeIndex()
	_ = f.Freeze("prod")
	if err := f.AssertMutable("prod"); err == nil {
		t.Fatal("expected error for frozen set")
	}
}

func TestAssertMutableUnfrozenReturnsNil(t *testing.T) {
	f := NewFreezeIndex()
	if err := f.AssertMutable("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
