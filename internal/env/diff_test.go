package env

import (
	"testing"
)

func makeDiffSet(t *testing.T, name string, vars map[string]string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	for k, v := range vars {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put: %v", err)
		}
	}
	return s
}

func TestDiffNilSrcReturnsError(t *testing.T) {
	b := makeDiffSet(t, "b", nil)
	_, err := Diff(nil, b)
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestDiffNilDstReturnsError(t *testing.T) {
	a := makeDiffSet(t, "a", nil)
	_, err := Diff(a, nil)
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestDiffIdenticalSetsIsEmpty(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	a := makeDiffSet(t, "a", vars)
	b := makeDiffSet(t, "b", vars)
	d, err := Diff(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !d.IsEmpty() {
		t.Errorf("expected empty diff, got added=%v removed=%v changed=%v", d.Added, d.Removed, d.Changed)
	}
}

func TestDiffDetectsAdded(t *testing.T) {
	a := makeDiffSet(t, "a", map[string]string{"FOO": "1"})
	b := makeDiffSet(t, "b", map[string]string{"FOO": "1", "BAR": "2"})
	d, err := Diff(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := d.Added["BAR"]; !ok || v != "2" {
		t.Errorf("expected BAR=2 in Added, got %v", d.Added)
	}
}

func TestDiffDetectsRemoved(t *testing.T) {
	a := makeDiffSet(t, "a", map[string]string{"FOO": "1", "BAR": "2"})
	b := makeDiffSet(t, "b", map[string]string{"FOO": "1"})
	d, _ := Diff(a, b)
	if v, ok := d.Removed["BAR"]; !ok || v != "2" {
		t.Errorf("expected BAR in Removed, got %v", d.Removed)
	}
}

func TestDiffDetectsChanged(t *testing.T) {
	a := makeDiffSet(t, "a", map[string]string{"FOO": "old"})
	b := makeDiffSet(t, "b", map[string]string{"FOO": "new"})
	d, _ := Diff(a, b)
	if pair, ok := d.Changed["FOO"]; !ok || pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected FOO changed old->new, got %v", d.Changed)
	}
}
