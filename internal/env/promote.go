package env

import "fmt"

// PromoteResult holds the outcome of a promotion between two sets.
type PromoteResult struct {
	Promoted []string // keys copied to destination
	Skipped  []string // keys skipped because they already exist and overwrite=false
}

// Promote copies all variables from src into dst.
// When overwrite is false, keys that already exist in dst are skipped and
// recorded in PromoteResult.Skipped. When overwrite is true every key from
// src is written to dst regardless.
func Promote(src, dst *Set, overwrite bool) (*PromoteResult, error) {
	if src == nil {
		return nil, fmt.Errorf("promote: src set must not be nil")
	}
	if dst == nil {
		return nil, fmt.Errorf("promote: dst set must not be nil")
	}

	result := &PromoteResult{}

	src.mu.RLock()
	defer src.mu.RUnlock()

	for k, v := range src.vars {
		_, err := dst.Get(k)
		exists := err == nil

		if exists && !overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if putErr := dst.Put(k, v); putErr != nil {
			return nil, fmt.Errorf("promote: writing key %q: %w", k, putErr)
		}
		result.Promoted = append(result.Promoted, k)
	}

	return result, nil
}

// PromoteKeys copies only the specified keys from src into dst, obeying the
// same overwrite semantics as Promote.
func PromoteKeys(src, dst *Set, keys []PromoteResult, error) {
	if src == nil {
		return nil, fmt.Errorf("promote: src set must not be nil")
	}
	if dst == nil {
		return nil, fmt.Errorf("promote: dst set must not be nil")
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("promote: keys list must not be empty")
	}

	result := &PromoteResult{}

	for _, k := range keys {
		v, err := src.Get(k)
		if err != nil {
			return nil, fmt.Errorf("promote: key %q not found in src: %w", k, err)
		}

		_, dstErr := dst.Get(k)
		exists := dstErr == nil

		if exists && !overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		if putErr := dst.Put(k, v); putErr != nil {
			return nil, fmt.Errorf("promote: writing key %q: %w", k, putErr)
		}
		result.Promoted = append(result.Promoted, k)
	}

	return result, nil
}
