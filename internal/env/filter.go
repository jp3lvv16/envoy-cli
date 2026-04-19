package env

import "fmt"

// FilterFunc is a predicate applied to each key-value pair in a Set.
type FilterFunc func(key, value string) bool

// Filter returns a new Set containing only the entries from src for which fn
// returns true. The new Set has the same name as src.
func Filter(src *Set, fn FilterFunc) (*Set, error) {
	if src == nil {
		return nil, fmt.Errorf("filter: src set must not be nil")
	}
	if fn == nil {
		return nil, fmt.Errorf("filter: filter function must not be nil")
	}

	dst, err := NewSet(src.Name())
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}

	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		if fn(k, v) {
			if err := dst.Put(k, v); err != nil {
				return nil, fmt.Errorf("filter: %w", err)
			}
		}
	}

	return dst, nil
}

// FilterByPrefix returns a new Set containing only entries whose keys start
// with the given prefix.
func FilterByPrefix(src *Set, prefix string) (*Set, error) {
	return Filter(src, func(key, _ string) bool {
		return len(key) >= len(prefix) && key[:len(prefix)] == prefix
	})
}
