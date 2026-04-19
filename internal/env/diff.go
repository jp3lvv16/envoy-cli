package env

import "fmt"

// DiffResult holds the differences between two Sets.
type DiffResult struct {
	Added   map[string]string // keys present in b but not a
	Removed map[string]string // keys present in a but not b
	Changed map[string][2]string // keys in both but with different values [old, new]
}

// Diff compares two Sets and returns a DiffResult describing the changes
// needed to transform a into b.
func Diff(a, b *Set) (*DiffResult, error) {
	if a == nil {
		return nil, fmt.Errorf("diff: source set must not be nil")
	}
	if b == nil {
		return nil, fmt.Errorf("diff: target set must not be nil")
	}

	result := &DiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	// Find removed and changed
	for k, v := range a.Vars {
		if bv, ok := b.Vars[k]; !ok {
			result.Removed[k] = v
		} else if bv != v {
			result.Changed[k] = [2]string{v, bv}
		}
	}

	// Find added
	for k, v := range b.Vars {
		if _, ok := a.Vars[k]; !ok {
			result.Added[k] = v
		}
	}

	return result, nil
}

// IsEmpty returns true when there are no differences.
func (d *DiffResult) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}
