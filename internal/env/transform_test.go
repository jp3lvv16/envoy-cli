package env

import (
	"strings"
	"testing"
)

func makeTransformSet(t *testing.T) *Set {
	t.Helper()
	s, _ := NewSet("xform")
	_ = s.Put("HOST", "localhost")
	_ = s.Put("ENV", "dev")
	return s
}

func TestTransformNilSetReturnsError(t *testing.T) {
	_, err := Transform(nil, UppercaseValues())
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestTransformNilFuncReturnsError(t *testing.T) {
	s := makeTransformSet(t)
	_, err := Transform(s, nil)
	if err == nil {
		t.Fatal("expected error for nil func")
	}
}

func TestTransformUppercaseValues(t *testing.T) {
	s := makeTransformSet(t)
	out, err := Transform(s, UppercaseValues())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("HOST")
	if v != "LOCALHOST" {
		t.Errorf("expected LOCALHOST, got %s", v)
	}
}

func TestTransformPrefixValues(t *testing.T) {
	s := makeTransformSet(t)
	out, err := Transform(s, PrefixValues("pfx_"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("ENV")
	if !strings.HasPrefix(v, "pfx_") {
		t.Errorf("expected prefix, got %s", v)
	}
}

func TestTransformDoesNotMutateSrc(t *testing.T) {
	s := makeTransformSet(t)
	_, _ = Transform(s, UppercaseValues())
	v, _ := s.Get("HOST")
	if v != "localhost" {
		t.Errorf("src mutated, expected localhost got %s", v)
	}
}

func TestTransformPreservesSetName(t *testing.T) {
	s := makeTransformSet(t)
	out, _ := Transform(s, UppercaseValues())
	if out.Name() != s.Name() {
		t.Errorf("expected name %s, got %s", s.Name(), out.Name())
	}
}
