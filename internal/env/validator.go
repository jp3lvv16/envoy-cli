package env

import (
	"fmt"
	"regexp"
)

// validKeyRe matches POSIX-style environment variable names.
var validKeyRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidationError holds all violations found in a Set.
type ValidationError struct {
	Violations []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed with %d violation(s): %v", len(e.Violations), e.Violations)
}

// Validate checks every key in s against POSIX naming rules and ensures no
// value is empty (unless the caller explicitly allows it via allowEmpty).
func Validate(s *Set, allowEmpty bool) error {
	if s == nil {
		return fmt.Errorf("validate: set must not be nil")
	}

	var violations []string

	for _, k := range s.Keys() {
		if !validKeyRe.MatchString(k) {
			violations = append(violations, fmt.Sprintf("invalid key %q: must match [A-Za-z_][A-Za-z0-9_]*", k))
		}
		if !allowEmpty {
			v, _ := s.Get(k)
			if v == "" {
				violations = append(violations, fmt.Sprintf("key %q has an empty value", k))
			}
		}
	}

	if len(violations) > 0 {
		return &ValidationError{Violations: violations}
	}
	return nil
}
