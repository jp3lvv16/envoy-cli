package env

import "fmt"

// MergeStrategy defines how conflicts are resolved when merging two sets.
type MergeStrategy int

const (
	// MergeStrategySkip keeps the destination value on conflict.
	MergeStrategySkip MergeStrategy = iota
	// MergeStrategyOverwrite replaces the destination value on conflict.
	MergeStrategyOverwrite
	// MergeStrategyError returns an error on conflict.
	MergeStrategyError
)

// MergeWithStrategy merges src into dst using the provided strategy to resolve
// key conflicts. dst is modified in-place; src is never modified.
func MergeWithStrategy(src, dst *Set, strategy MergeStrategy) error {
	if src == nil {
		return fmt.Errorf("merge: src set must not be nil")
	}
	if dst == nil {
		return fmt.Errorf("merge: dst set must not be nil")
	}

	src.mu.RLock()
	defer src.mu.RUnlock()

	for k, v := range src.vars {
		_, exists := dst.vars[k]
		switch {
		case !exists:
			if err := dst.Put(k, v); err != nil {
				return fmt.Errorf("merge: %w", err)
			}
		case strategy == MergeStrategyOverwrite:
			if err := dst.Put(k, v); err != nil {
				return fmt.Errorf("merge: %w", err)
			}
		case strategy == MergeStrategyError:
			return fmt.Errorf("merge: conflict on key %q", k)
		// MergeStrategySkip: do nothing
		}
	}
	return nil
}
