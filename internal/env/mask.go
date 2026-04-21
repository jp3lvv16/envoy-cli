package env

import (
	"errors"
	"strings"
)

// MaskResult holds the masked representation of a Set's variables.
type MaskResult struct {
	Name string
	Vars map[string]string
}

// defaultSensitiveSubstrings are common key substrings that indicate sensitive values.
var defaultSensitiveSubstrings = []string{
	"password", "passwd", "secret", "token", "apikey", "api_key",
	"auth", "credential", "private", "key",
}

const maskedValue = "********"

// Mask returns a MaskResult where values whose keys match any sensitive
// substring are replaced with a redacted placeholder. The original Set
// is never modified.
func Mask(s *Set, sensitiveSubstrings []string) (*MaskResult, error) {
	if s == nil {
		return nil, errors.New("mask: set must not be nil")
	}

	substrings := sensitiveSubstrings
	if len(substrings) == 0 {
		substrings = defaultSensitiveSubstrings
	}

	result := &MaskResult{
		Name: s.Name(),
		Vars: make(map[string]string),
	}

	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		if isSensitive(k, substrings) {
			result.Vars[k] = maskedValue
		} else {
			result.Vars[k] = v
		}
	}

	return result, nil
}

// isSensitive returns true if the key contains any of the given substrings
// (case-insensitive).
func isSensitive(key string, substrings []string) bool {
	lower := strings.ToLower(key)
	for _, sub := range substrings {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
