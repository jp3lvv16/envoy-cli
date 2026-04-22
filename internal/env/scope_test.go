package env

import (
	"testing"
)

func makeScopeSet(t *testing.T, name string, pairs map[string]string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	for k, v := range pairs {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put(%q): %v", k, err)
		}
	}
	return s
}

func TestNewScopeResolverNoScopesReturnsError(t *testing.T) {
	_, err := NewScopeResolver(nil)
	if err == nil {
		t.Fatal("expected error for empty scopes")
	}
}

func TestNewScopeResolverNilSetReturnsError(t *testing.T) {
	_, err := NewScopeResolver([]*Scope{{Name: "dev", Priority: 1, Set: nil}})
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestResolveEmptyKeyReturnsError(t *testing.T) {
	s := makeScopeSet(t, "dev", map[string]string{"A": "1"})
	r, _ := NewScopeResolver([]*Scope{{Name: "dev", Priority: 1, Set: s}})
	_, _, err := r.Resolve("")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestResolveMissingKeyReturnsError(t *testing.T) {
	s := makeScopeSet(t, "dev", map[string]string{"A": "1"})
	r, _ := NewScopeResolver([]*Scope{{Name: "dev", Priority: 1, Set: s}})
	_, _, err := r.Resolve("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestResolvePicksHighestPriority(t *testing.T) {
	dev := makeScopeSet(t, "dev", map[string]string{"DB_HOST": "localhost"})
	prod := makeScopeSet(t, "prod", map[string]string{"DB_HOST": "prod.db.example.com"})
	r, err := NewScopeResolver([]*Scope{
		{Name: "dev", Priority: 1, Set: dev},
		{Name: "prod", Priority: 10, Set: prod},
	})
	if err != nil {
		t.Fatalf("NewScopeResolver: %v", err)
	}
	val, scope, err := r.Resolve("DB_HOST")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if scope != "prod" {
		t.Errorf("expected winning scope prod, got %q", scope)
	}
	if val != "prod.db.example.com" {
		t.Errorf("unexpected value %q", val)
	}
}

func TestResolveAllMergesScopes(t *testing.T) {
	base := makeScopeSet(t, "base", map[string]string{"A": "base_a", "B": "base_b"})
	override := makeScopeSet(t, "override", map[string]string{"B": "over_b", "C": "over_c"})
	r, err := NewScopeResolver([]*Scope{
		{Name: "base", Priority: 1, Set: base},
		{Name: "override", Priority: 5, Set: override},
	})
	if err != nil {
		t.Fatalf("NewScopeResolver: %v", err)
	}
	out, err := r.ResolveAll("merged")
	if err != nil {
		t.Fatalf("ResolveAll: %v", err)
	}
	check := func(key, want string) {
		t.Helper()
		v, err := out.Get(key)
		if err != nil {
			t.Fatalf("Get(%q): %v", key, err)
		}
		if v != want {
			t.Errorf("key %q: want %q got %q", key, want, v)
		}
	}
	check("A", "base_a")
	check("B", "over_b")
	check("C", "over_c")
}

func TestResolveAllEmptyNameReturnsError(t *testing.T) {
	s := makeScopeSet(t, "dev", map[string]string{"X": "1"})
	r, _ := NewScopeResolver([]*Scope{{Name: "dev", Priority: 1, Set: s}})
	_, err := r.ResolveAll("")
	if err == nil {
		t.Fatal("expected error for empty result name")
	}
}
