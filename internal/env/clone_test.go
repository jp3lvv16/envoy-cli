package env

import (
	"testing"
)

func makeCloneSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("original")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = s.Put("HOST", "localhost")
	_ = s.Put("PORT", "8080")
	return s
}

func TestCloneNilSrcReturnsError(t *testing.T) {
	_, err := Clone(nil, "copy")
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestCloneEmptyNameReturnsError(t *testing.T) {
	src := makeCloneSet(t)
	_, err := Clone(src, "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCloneCopiesAllVars(t *testing.T) {
	src := makeCloneSet(t)
	dst, err := Clone(src, "copy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst.Name() != "copy" {
		t.Errorf("expected name %q, got %q", "copy", dst.Name())
	}
	for k, v := range src.Vars() {
		got, err := dst.Get(k)
		if err != nil {
			t.Errorf("missing key %q in clone", k)
			continue
		}
		if got != v {
			t.Errorf("key %q: expected %q, got %q", k, v, got)
		}
	}
}

func TestCloneIsIndependent(t *testing.T) {
	src := makeCloneSet(t)
	dst, err := Clone(src, "copy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = dst.Put("HOST", "changed")
	v, _ := src.Get("HOST")
	if v != "localhost" {
		t.Errorf("clone mutation affected source: got %q", v)
	}
}

func TestCloneOriginalNameUnchanged(t *testing.T) {
	src := makeCloneSet(t)
	_, _ = Clone(src, "copy")
	if src.Name() != "original" {
		t.Errorf("source name changed to %q", src.Name())
	}
}
