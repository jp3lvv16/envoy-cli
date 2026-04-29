package env_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/env"
)

func makeInterpolateSet(t *testing.T, name string, pairs map[string]string) *env.Set {
	t.Helper()
	s, err := env.NewSet(name)
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

func TestInterpolateNilSetReturnsError(t *testing.T) {
	_, err := env.Interpolate(nil)
	if err == nil {
		t.Fatal("expected error for nil set, got nil")
	}
}

func TestInterpolateNoReferences(t *testing.T) {
	s := makeInterpolateSet(t, "plain", map[string]string{
		"FOO": "hello",
		"BAR": "world",
	})
	out, err := env.Interpolate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := out.Get("FOO"); v != "hello" {
		t.Errorf("FOO: got %q, want %q", v, "hello")
	}
	if v, _ := out.Get("BAR"); v != "world" {
		t.Errorf("BAR: got %q, want %q", v, "world")
	}
}

func TestInterpolateResolvesReference(t *testing.T) {
	s := makeInterpolateSet(t, "refs", map[string]string{
		"BASE": "/app",
		"DATA": "${BASE}/data",
	})
	out, err := env.Interpolate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := out.Get("DATA"); v != "/app/data" {
		t.Errorf("DATA: got %q, want %q", v, "/app/data")
	}
}

func TestInterpolateChainedReferences(t *testing.T) {
	s := makeInterpolateSet(t, "chain", map[string]string{
		"A": "alpha",
		"B": "${A}-beta",
		"C": "${B}-gamma",
	})
	out, err := env.Interpolate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := out.Get("C"); v != "alpha-beta-gamma" {
		t.Errorf("C: got %q, want %q", v, "alpha-beta-gamma")
	}
}

func TestInterpolateMissingRefLeftInPlace(t *testing.T) {
	s := makeInterpolateSet(t, "missing", map[string]string{
		"FOO": "${UNDEFINED}_suffix",
	})
	out, err := env.Interpolate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := out.Get("FOO"); v != "${UNDEFINED}_suffix" {
		t.Errorf("FOO: got %q, want %q", v, "${UNDEFINED}_suffix")
	}
}

func TestInterpolateWithExternalVars(t *testing.T) {
	s := makeInterpolateSet(t, "external", map[string]string{
		"GREETING": "${SALUTATION} world",
	})
	ext := map[string]string{"SALUTATION": "hello"}
	out, err := env.InterpolateWith(s, ext)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := out.Get("GREETING"); v != "hello world" {
		t.Errorf("GREETING: got %q, want %q", v, "hello world")
	}
}

func TestInterpolateWithNilSetReturnsError(t *testing.T) {
	_, err := env.InterpolateWith(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil set, got nil")
	}
}

func TestInterpolateDoesNotMutateOriginal(t *testing.T) {
	s := makeInterpolateSet(t, "immutable", map[string]string{
		"BASE": "/opt",
		"PATH": "${BASE}/bin",
	})
	_, err := env.Interpolate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Original should remain unexpanded
	if v, _ := s.Get("PATH"); v != "${BASE}/bin" {
		t.Errorf("original PATH mutated: got %q, want %q", v, "${BASE}/bin")
	}
}

func TestInterpolateCyclicReferenceReturnsError(t *testing.T) {
	s := makeInterpolateSet(t, "cyclic", map[string]string{
		"X": "${Y}",
		"Y": "${X}",
	})
	_, err := env.Interpolate(s)
	if err == nil {
		t.Fatal("expected error for cyclic reference, got nil")
	}
}
