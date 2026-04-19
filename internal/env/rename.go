package env

import "fmt"

// Rename changes the name of an existing Set, returning a new Set with the
// updated name and the same variables. The original Set is not modified.
func Rename(s *Set, newName string) (*Set, error) {
	if s == nil {
		return nil, fmt.Errorf("rename: source set must not be nil")
	}
	if newName == "" {
		return nil, fmt.Errorf("rename: new name must not be empty")
	}

	dst, err := NewSet(newName)
	if err != nil {
		return nil, fmt.Errorf("rename: %w", err)
	}

	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		if err := dst.Put(k, v); err != nil {
			return nil, fmt.Errorf("rename: copy key %q: %w", k, err)
		}
	}

	return dst, nil
}
