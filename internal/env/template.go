package env

import (
	"fmt"
	"regexp"
	"strings"
)

// templateVarRe matches {{VAR_NAME}} placeholders.
var templateVarRe = regexp.MustCompile(`\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}`)

// RenderResult holds the output of a template rendering operation.
type RenderResult struct {
	Output   string
	Missing  []string
}

// Render replaces {{KEY}} placeholders in tmpl with values from s.
// Missing keys are collected in RenderResult.Missing; no error is returned
// for absent keys — the placeholder is left as-is.
func Render(s *Set, tmpl string) (*RenderResult, error) {
	if s == nil {
		return nil, fmt.Errorf("render: set must not be nil")
	}
	if tmpl == "" {
		return &RenderResult{Output: ""}, nil
	}

	seen := map[string]bool{}
	var missing []string

	output := templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := strings.TrimSuffix(strings.TrimPrefix(match, "{{"), "}}")
		val, err := s.Get(key)
		if err != nil {
			if !seen[key] {
				missing = append(missing, key)
				seen[key] = true
			}
			return match
		}
		return val
	})

	return &RenderResult{Output: output, Missing: missing}, nil
}

// RenderStrict is like Render but returns an error if any placeholder is unresolved.
func RenderStrict(s *Set, tmpl string) (string, error) {
	res, err := Render(s, tmpl)
	if err != nil {
		return "", err
	}
	if len(res.Missing) > 0 {
		return "", fmt.Errorf("render: unresolved placeholders: %s", strings.Join(res.Missing, ", "))
	}
	return res.Output, nil
}
