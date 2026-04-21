package env

import (
	"strings"
	"testing"
)

func makeFilterSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("filter-test")
	if err != nil {
		t.Fatalf("makeFilterSet: %v", err)
	}
	_ = s.Put("APP_HOST", "localhost")
	_ = s.Put("APP_PORT", "8080")
	_ = s.Put("DB_HOST", "db.local")
	_ = s.Put("DB_PORT", "5432")
	return s
}

func TestFilterNilSetReturnsError(t *testing.T) {
	_, err := Filter(nil, func(k, v string) bool { return true })
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestFilterNilFuncReturnsError(t *testing.T) {
	s := makeFilterSet(t)
	_, err := Filter(s, nil)
	if err == nil {
		t.Fatal("expected error for nil func")
	}
}

func TestFilterKeepsMatchingEntries(t *testing.T) {
	s := makeFilterSet(t)
	out, err := Filter(s, func(k, _ string) bool {
		return strings.HasPrefix(k, "APP_")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Keys()) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out.Keys()))
	}
	if v, _ := out.Get("APP_HOST"); v != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", v)
	}
}

func TestFilterEmptyResultIsValid(t *testing.T) {
	s := makeFilterSet(t)
	out, err := Filter(s, func(k, v string) bool { return false })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Keys()) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(out.Keys()))
	}
}

func TestFilterByPrefix(t *testing.T) {
	s := makeFilterSet(t)
	out, err := FilterByPrefix(s, "DB_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Keys()) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out.Keys()))
	}
	if v, _ := out.Get("DB_PORT"); v != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", v)
	}
}

func TestFilterPreservesSetName(t *testing.T) {
	s := makeFilterSet(t)
	out, _ := FilterByPrefix(s, "APP_")
	if out.Name() != s.Name() {
		t.Errorf("expected name %q, got %q", s.Name(), out.Name())
	}
}

func TestFilterByPrefixEmptyPrefixMatchesAll(t *testing.T) {
	s := makeFilterSet(t)
	out, err := FilterByPrefix(s, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Keys()) != len(s.Keys()) {
		t.Fatalf("expected %d keys, got %d", len(s.Keys()), len(out.Keys()))
	}
}
