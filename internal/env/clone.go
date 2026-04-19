package env

import "fmt"

// Clone creates a deep copy of src with a new name.
// The original set is left unchanged.
func Clone(src *Set, newName string) (*Set, error) {
	if src == nil {
		return nil, fmt.Errorf("clone: source set is nil")
	}
	if newName == "" {
		return nil, fmt.Errorf("clone: new name must not be empty")
	}

	dst, err := NewSet(newName)
	if err != nil {
		return nil, fmt.Errorf("clone: %w", err)
	}

	for k, v := range src.Vars() {
		if err := dst.Put(k, v); err != nil {
			return nil, fmt.Errorf("clone: copying key %q: %w", k, err)
		}
	}

	return dst, nil
}
