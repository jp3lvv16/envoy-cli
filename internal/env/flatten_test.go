package env

import (
	"testing"
)

func makeFlattenSet(t *testing.T, name string, pairs map[string]string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet(%q): %v", name, err)
	}
	for k, v := range pairs {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put(%q, %q): %v", k, v, err)
		}
	}
	return s
}

func TestFlattenEmptyNameReturnsError(t *testing.T) {
	s := makeFlattenSet(t, "a", map[string]string{"K": "v"})
	_, err := Flatten("", s)
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestFlattenNilSetReturnsError(t *testing.T) {
	_, err := Flatten("out", nil)
	if err == nil {
		t.Fatal("expected error for nil set, got nil")
	}
}

func TestFlattenMergesAllSets(t *testing.T) {
	a := makeFlattenSet(t, "a", map[string]string{"FOO": "1", "BAR": "2"})
	b := makeFlattenSet(t, "b", map[string]string{"BAZ": "3"})

	res, err := Flatten("merged", a, b)
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	if res.Merged != 3 {
		t.Errorf("Merged = %d, want 3", res.Merged)
	}
	for _, k := range []string{"FOO", "BAR", "BAZ"} {
		if _, err := res.Set.Get(k); err != nil {
			t.Errorf("key %q missing from flattened set", k)
		}
	}
}

func TestFlattenLaterSetOverwritesEarlier(t *testing.T) {
	a := makeFlattenSet(t, "a", map[string]string{"KEY": "old"})
	b := makeFlattenSet(t, "b", map[string]string{"KEY": "new"})

	res, err := Flatten("merged", a, b)
	if err != nil {
		t.Fatalf("Flatten: %v", err)
	}
	v, _ := res.Set.Get("KEY")
	if v != "new" {
		t.Errorf("KEY = %q, want %q", v, "new")
	}
}

func TestFlattenPrefixEmptyNameReturnsError(t *testing.T) {
	s := makeFlattenSet(t, "a", map[string]string{"K": "v"})
	_, err := FlattenPrefix("", "_", s)
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestFlattenPrefixNilSetReturnsError(t *testing.T) {
	_, err := FlattenPrefix("out", "_", nil)
	if err == nil {
		t.Fatal("expected error for nil set, got nil")
	}
}

func TestFlattenPrefixKeysAreNamespaced(t *testing.T) {
	prod := makeFlattenSet(t, "prod", map[string]string{"HOST": "prod.example.com"})
	stg := makeFlattenSet(t, "stg", map[string]string{"HOST": "stg.example.com"})

	res, err := FlattenPrefix("all", "_", prod, stg)
	if err != nil {
		t.Fatalf("FlattenPrefix: %v", err)
	}

	for _, k := range []string{"PROD_HOST", "STG_HOST"} {
		if _, err := res.Set.Get(k); err != nil {
			t.Errorf("key %q missing from prefixed set", k)
		}
	}
	if res.Merged != 2 {
		t.Errorf("Merged = %d, want 2", res.Merged)
	}
}

func TestFlattenEmptySetsProducesEmptyResult(t *testing.T) {
	res, err := Flatten("empty")
	if err != nil {
		t.Fatalf("Flatten with no sets: %v", err)
	}
	if res.Merged != 0 {
		t.Errorf("Merged = %d, want 0", res.Merged)
	}
}
