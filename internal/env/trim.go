package env

import (
	"errors"
	"strings"
)

// TrimResult holds the result of a trim operation on a single key.
type TrimResult struct {
	Key      string
	OldValue string
	NewValue string
	Changed  bool
}

// Trim removes leading and trailing whitespace from all values in the set.
// Returns a slice of TrimResult describing what changed.
func Trim(s *Set) ([]TrimResult, error) {
	if s == nil {
		return nil, errors.New("trim: set must not be nil")
	}

	keys, err := SortedKeys(s, true)
	if err != nil {
		return nil, err
	}

	var results []TrimResult
	for _, k := range keys {
		val, err := s.Get(k)
		if err != nil {
			continue
		}
		trimmed := strings.TrimSpace(val)
		changed := trimmed != val
		if changed {
			if err := s.Put(k, trimmed); err != nil {
				return nil, err
			}
		}
		results = append(results, TrimResult{
			Key:      k,
			OldValue: val,
			NewValue: trimmed,
			Changed:  changed,
		})
	}
	return results, nil
}

// TrimKeys returns only the TrimResults where a change was made.
func TrimKeys(s *Set) ([]TrimResult, error) {
	all, err := Trim(s)
	if err != nil {
		return nil, err
	}
	var changed []TrimResult
	for _, r := range all {
		if r.Changed {
			changed = append(changed, r)
		}
	}
	return changed, nil
}
