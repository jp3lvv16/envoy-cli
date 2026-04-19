package env

import "fmt"

// CopyOptions controls behaviour when copying variables between sets.
type CopyOptions struct {
	// Overwrite existing keys in the destination set.
	Overwrite bool
}

// Copy copies variables from src into dst.
// If opts.Overwrite is false, existing keys in dst are left unchanged.
// Returns the number of variables copied.
func Copy(src, dst *Set, opts CopyOptions) (int, error) {
	if src == nil {
		return 0, fmt.Errorf("copy: source set must not be nil")
	}
	if dst == nil {
		return 0, fmt.Errorf("copy: destination set must not be nil")
	}

	copied := 0
	for k, v := range src.vars {
		if !opts.Overwrite {
			if _, err := dst.Get(k); err == nil {
				// key already exists, skip
				continue
			}
		}
		if err := dst.Put(k, v); err != nil {
			return copied, fmt.Errorf("copy: failed to set key %q: %w", k, err)
		}
		copied++
	}
	return copied, nil
}

// Merge is an alias for Copy with Overwrite set to true.
func Merge(src, dst *Set) (int, error) {
	return Copy(src, dst, CopyOptions{Overwrite: true})
}
