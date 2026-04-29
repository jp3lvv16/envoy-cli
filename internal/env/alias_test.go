package env

import (
	"testing"
)

func TestAddEmptyAliasReturnsError(t *testing.T) {
	idx := NewAliasIndex()
	if err := idx.Add("", "prod"); err == nil {
		t.Fatal("expected error for empty alias")
	}
}

func TestAddEmptySetNameReturnsError(t *testing.T) {
	idx := NewAliasIndex()
	if err := idx.Add("p", ""); err == nil {
		t.Fatal("expected error for empty set name")
	}
}

func TestAddAndResolveAlias(t *testing.T) {
	idx := NewAliasIndex()
	if err := idx.Add("p", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	name, err := idx.Resolve("p")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "production" {
		t.Errorf("expected %q, got %q", "production", name)
	}
}

func TestAddDuplicateAliasIdempotent(t *testing.T) {
	idx := NewAliasIndex()
	_ = idx.Add("p", "production")
	if err := idx.Add("p", "production"); err != nil {
		t.Fatalf("duplicate same target should be idempotent, got: %v", err)
	}
}

func TestAddConflictingAliasReturnsError(t *testing.T) {
	idx := NewAliasIndex()
	_ = idx.Add("p", "production")
	if err := idx.Add("p", "staging"); err == nil {
		t.Fatal("expected error when alias points to different set")
	}
}

func TestRemoveAliasDeletesEntry(t *testing.T) {
	idx := NewAliasIndex()
	_ = idx.Add("p", "production")
	if err := idx.Remove("p"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := idx.Resolve("p"); err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveMissingAliasReturnsError(t *testing.T) {
	idx := NewAliasIndex()
	if err := idx.Remove("ghost"); err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestResolveMissingAliasReturnsError(t *testing.T) {
	idx := NewAliasIndex()
	if _, err := idx.Resolve("nope"); err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestAllReturnsCopy(t *testing.T) {
	idx := NewAliasIndex()
	_ = idx.Add("p", "production")
	_ = idx.Add("s", "staging")
	all := idx.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutating the copy must not affect the index.
	delete(all, "p")
	if _, err := idx.Resolve("p"); err != nil {
		t.Fatal("index should not be affected by mutation of All() copy")
	}
}
