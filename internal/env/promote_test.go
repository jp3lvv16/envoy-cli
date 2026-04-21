package env

import (
	"testing"
)

func makePromoteSet(t *testing.T, name string, pairs map[string]string) *Set {
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

func TestPromoteNilSrcReturnsError(t *testing.T) {
	dst := makePromoteSet(t, "dst", nil)
	_, err := Promote(nil, dst, false)
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestPromoteNilDstReturnsError(t *testing.T) {
	src := makePromoteSet(t, "src", nil)
	_, err := Promote(src, nil, false)
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestPromoteNoOverwriteSkipsExisting(t *testing.T) {
	src := makePromoteSet(t, "src", map[string]string{"A": "1", "B": "2"})
	dst := makePromoteSet(t, "dst", map[string]string{"A": "original"})

	res, err := Promote(src, dst, false)
	if err != nil {
		t.Fatalf("Promote: %v", err)
	}

	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}

	// A must retain original value
	v, _ := dst.Get("A")
	if v != "original" {
		t.Errorf("expected A=original, got %q", v)
	}

	// B must be promoted
	v, _ = dst.Get("B")
	if v != "2" {
		t.Errorf("expected B=2, got %q", v)
	}
}

func TestPromoteWithOverwriteUpdatesAll(t *testing.T) {
	src := makePromoteSet(t, "src", map[string]string{"A": "new", "B": "2"})
	dst := makePromoteSet(t, "dst", map[string]string{"A": "old"})

	res, err := Promote(src, dst, true)
	if err != nil {
		t.Fatalf("Promote: %v", err)
	}

	if len(res.Skipped) != 0 {
		t.Errorf("expected no skips, got %v", res.Skipped)
	}

	v, _ := dst.Get("A")
	if v != "new" {
		t.Errorf("expected A=new, got %q", v)
	}
}

func TestPromoteKeysEmptyListReturnsError(t *testing.T) {
	src := makePromoteSet(t, "src", map[string]string{"A": "1"})
	dst := makePromoteSet(t, "dst", nil)
	_, err := PromoteKeys(src, dst, []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestPromoteKeysMissingKeyReturnsError(t *testing.T) {
	src := makePromoteSet(t, "src", map[string]string{"A": "1"})
	dst := makePromoteSet(t, "dst", nil)
	_, err := PromoteKeys(src, dst, []string{"MISSING"}, false)
	if err == nil {
		t.Fatal("expected error for missing key in src")
	}
}

func TestPromoteKeysOnlySpecifiedKeys(t *testing.T) {
	src := makePromoteSet(t, "src", map[string]string{"A": "1", "B": "2", "C": "3"})
	dst := makePromoteSet(t, "dst", nil)

	res, err := PromoteKeys(src, dst, []string{"A", "C"}, false)
	if err != nil {
		t.Fatalf("PromoteKeys: %v", err)
	}

	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}

	if _, err := dst.Get("B"); err == nil {
		t.Error("B should not have been promoted")
	}
}
