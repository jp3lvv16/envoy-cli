package env

import (
	"testing"
)

func makeIntersectSet(t *testing.T, name string, pairs map[string]string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	for k, v := range pairs {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put(%q, %q): %v", k, v, err)
		}
	}
	return s
}

func TestIntersectNilSrcReturnsError(t *testing.T) {
	dst := makeIntersectSet(t, "dst", map[string]string{"A": "1"})
	_, err := Intersect(nil, dst, "result")
	if err == nil {
		t.Fatal("expected error for nil src, got nil")
	}
}

func TestIntersectNilDstReturnsError(t *testing.T) {
	src := makeIntersectSet(t, "src", map[string]string{"A": "1"})
	_, err := Intersect(src, nil, "result")
	if err == nil {
		t.Fatal("expected error for nil dst, got nil")
	}
}

func TestIntersectEmptyNameReturnsError(t *testing.T) {
	src := makeIntersectSet(t, "src", map[string]string{"A": "1"})
	dst := makeIntersectSet(t, "dst", map[string]string{"A": "1"})
	_, err := Intersect(src, dst, "")
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestIntersectReturnsCommonKeys(t *testing.T) {
	src := makeIntersectSet(t, "src", map[string]string{"A": "1", "B": "2", "C": "3"})
	dst := makeIntersectSet(t, "dst", map[string]string{"B": "99", "C": "99", "D": "4"})

	result, err := Intersect(src, dst, "common")
	if err != nil {
		t.Fatalf("Intersect: %v", err)
	}

	for _, key := range []string{"B", "C"} {
		if _, err := result.Get(key); err != nil {
			t.Errorf("expected key %q in result, but missing", key)
		}
	}
	if _, err := result.Get("A"); err == nil {
		t.Error("key A should not be in result")
	}
	if _, err := result.Get("D"); err == nil {
		t.Error("key D should not be in result")
	}
}

func TestIntersectValuesFromSrc(t *testing.T) {
	src := makeIntersectSet(t, "src", map[string]string{"X": "src-val"})
	dst := makeIntersectSet(t, "dst", map[string]string{"X": "dst-val"})

	result, err := Intersect(src, dst, "res")
	if err != nil {
		t.Fatalf("Intersect: %v", err)
	}
	v, _ := result.Get("X")
	if v != "src-val" {
		t.Errorf("expected src-val, got %q", v)
	}
}

func TestSubtractNilSrcReturnsError(t *testing.T) {
	dst := makeIntersectSet(t, "dst", map[string]string{"A": "1"})
	_, err := Subtract(nil, dst, "result")
	if err == nil {
		t.Fatal("expected error for nil src, got nil")
	}
}

func TestSubtractReturnsOnlyUniqueKeys(t *testing.T) {
	src := makeIntersectSet(t, "src", map[string]string{"A": "1", "B": "2", "C": "3"})
	dst := makeIntersectSet(t, "dst", map[string]string{"B": "99"})

	result, err := Subtract(src, dst, "diff")
	if err != nil {
		t.Fatalf("Subtract: %v", err)
	}

	for _, key := range []string{"A", "C"} {
		if _, err := result.Get(key); err != nil {
			t.Errorf("expected key %q in result, but missing", key)
		}
	}
	if _, err := result.Get("B"); err == nil {
		t.Error("key B should not be in result")
	}
}
