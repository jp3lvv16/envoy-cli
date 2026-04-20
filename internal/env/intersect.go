package env

import "fmt"

// Intersect returns a new Set containing only the keys that exist in both src
// and dst. The values are taken from src. The resulting set is named after the
// provided name argument.
func Intersect(src, dst *Set, name string) (*Set, error) {
	if src == nil {
		return nil, fmt.Errorf("intersect: src set must not be nil")
	}
	if dst == nil {
		return nil, fmt.Errorf("intersect: dst set must not be nil")
	}
	if name == "" {
		return nil, fmt.Errorf("intersect: name must not be empty")
	}

	result, err := NewSet(name)
	if err != nil {
		return nil, fmt.Errorf("intersect: %w", err)
	}

	src.mu.RLock()
	defer src.mu.RUnlock()
	dst.mu.RLock()
	defer dst.mu.RUnlock()

	for k, v := range src.vars {
		if _, exists := dst.vars[k]; exists {
			if err := result.Put(k, v); err != nil {
				return nil, fmt.Errorf("intersect: %w", err)
			}
		}
	}

	return result, nil
}

// Subtract returns a new Set containing only the keys from src that do NOT
// exist in dst. The resulting set is named after the provided name argument.
func Subtract(src, dst *Set, name string) (*Set, error) {
	if src == nil {
		return nil, fmt.Errorf("subtract: src set must not be nil")
	}
	if dst == nil {
		return nil, fmt.Errorf("subtract: dst set must not be nil")
	}
	if name == "" {
		return nil, fmt.Errorf("subtract: name must not be empty")
	}

	result, err := NewSet(name)
	if err != nil {
		return nil, fmt.Errorf("subtract: %w", err)
	}

	src.mu.RLock()
	defer src.mu.RUnlock()
	dst.mu.RLock()
	defer dst.mu.RUnlock()

	for k, v := range src.vars {
		if _, exists := dst.vars[k]; !exists {
			if err := result.Put(k, v); err != nil {
				return nil, fmt.Errorf("subtract: %w", err)
			}
		}
	}

	return result, nil
}
