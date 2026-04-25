package env

import (
	"strings"
	"testing"
)

func makeLintSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("lint-test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	return s
}

func TestLintNilSetReturnsError(t *testing.T) {
	_, err := Lint(nil)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestLintCleanSetNoIssues(t *testing.T) {
	s := makeLintSet(t)
	s.Put("HOST", "localhost")
	s.Put("PORT", "8080")

	results, err := Lint(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no issues, got %d: %v", len(results), results)
	}
}

func TestLintDetectsLowercaseKey(t *testing.T) {
	s := makeLintSet(t)
	s.Put("host", "localhost")

	results, err := Lint(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(results, "uppercase-keys") {
		t.Error("expected uppercase-keys violation")
	}
}

func TestLintDetectsWhitespacePadding(t *testing.T) {
	s := makeLintSet(t)
	s.Put("HOST", "  localhost  ")

	results, err := Lint(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(results, "no-whitespace-padding") {
		t.Error("expected no-whitespace-padding violation")
	}
}

func TestLintDetectsEmptyValue(t *testing.T) {
	s := makeLintSet(t)
	s.Put("HOST", "")

	results, err := Lint(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(results, "no-empty-values") {
		t.Error("expected no-empty-values violation")
	}
}

func TestFormatLintNoIssues(t *testing.T) {
	out := FormatLint(nil)
	if !strings.Contains(out, "no lint issues") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatLintShowsViolations(t *testing.T) {
	s := makeLintSet(t)
	s.Put("bad key", "")

	results, _ := Lint(s)
	out := FormatLint(results)
	if out == "" {
		t.Error("expected non-empty lint output")
	}
}

// helper
func containsRule(results []LintResult, rule string) bool {
	for _, r := range results {
		if r.Rule == rule {
			return true
		}
	}
	return false
}
