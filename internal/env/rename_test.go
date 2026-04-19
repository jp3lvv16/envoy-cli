package env

import (
	"testing"
)

func TestRenameNilSetReturnsError(t *testing.T) {
	_, err := Rename(nil, "new")
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestRenameEmptyNameReturnsError(t *testing.T) {
	s, _ := NewSet("original")
	_, err := Rename(s, "")
	if err == nil {
		t.Fatal("expected error for empty new name")
	}
}

func TestRenamePreservesVars(t *testing.T) {
	s, _ := NewSet("original")
	_ = s.Put("KEY", "value")
	_ = s.Put("FOO", "bar")

	renamed, err := Rename(s, "renamed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if renamed.Name() != "renamed" {
		t.Errorf("expected name %q, got %q", "renamed", renamed.Name())
	}

	for _, k := range []string{"KEY", "FOO"} {
		v, err := renamed.Get(k)
		if err != nil {
			t.Errorf("missing key %q after rename", k)
		}
		orig, _ := s.Get(k)
		if v != orig {
			t.Errorf("key %q: expected %q, got %q", k, orig, v)
		}
	}
}

func TestRenameIsIndependent(t *testing.T) {
	s, _ := NewSet("original")
	_ = s.Put("KEY", "before")

	renamed, _ := Rename(s, "copy")
	_ = s.Put("KEY", "after")

	v, _ := renamed.Get("KEY")
	if v != "before" {
		t.Errorf("rename should be independent of source; got %q", v)
	}
}

func TestRenameOriginalUnchanged(t *testing.T) {
	s, _ := NewSet("original")
	_, _ = Rename(s, "newname")

	if s.Name() != "original" {
		t.Errorf("source name should remain %q, got %q", "original", s.Name())
	}
}
