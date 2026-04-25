package env

import (
	"fmt"
	"strings"
)

// LintRule represents a single lint check applied to an env set.
type LintRule struct {
	Name    string
	Message string
}

// LintResult holds the outcome of a lint check.
type LintResult struct {
	Key  string
	Rule string
	Msg  string
}

// Lint runs a series of style and correctness checks on the given Set.
// It returns a slice of LintResult (one per violation) and any fatal error.
func Lint(s *Set) ([]LintResult, error) {
	if s == nil {
		return nil, fmt.Errorf("lint: set must not be nil")
	}

	var results []LintResult

	for _, k := range s.Keys() {
		v, _ := s.Get(k)

		// Rule: key should be uppercase
		if k != strings.ToUpper(k) {
			results = append(results, LintResult{
				Key:  k,
				Rule: "uppercase-keys",
				Msg:  fmt.Sprintf("key %q is not uppercase", k),
			})
		}

		// Rule: no leading/trailing whitespace in values
		if v != strings.TrimSpace(v) {
			results = append(results, LintResult{
				Key:  k,
				Rule: "no-whitespace-padding",
				Msg:  fmt.Sprintf("value for key %q has leading or trailing whitespace", k),
			})
		}

		// Rule: no empty values
		if strings.TrimSpace(v) == "" {
			results = append(results, LintResult{
				Key:  k,
				Rule: "no-empty-values",
				Msg:  fmt.Sprintf("key %q has an empty value", k),
			})
		}

		// Rule: key must not contain spaces
		if strings.Contains(k, " ") {
			results = append(results, LintResult{
				Key:  k,
				Rule: "no-spaces-in-keys",
				Msg:  fmt.Sprintf("key %q contains spaces", k),
			})
		}
	}

	return results, nil
}

// FormatLint returns a human-readable string for a slice of LintResult.
func FormatLint(results []LintResult) string {
	if len(results) == 0 {
		return "no lint issues found"
	}
	var sb strings.Builder
	for _, r := range results {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", r.Rule, r.Key, r.Msg))
	}
	return strings.TrimRight(sb.String(), "\n")
}
