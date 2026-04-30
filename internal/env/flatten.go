package env

import (
	"fmt"
	"strings"
)

// FlattenResult holds the output of a Flatten operation.
type FlattenResult struct {
	Set     *Set
	Merged  int
	Skipped int
}

// Flatten merges multiple sets into a single new set using the given name.
// Keys from later sets overwrite keys from earlier sets.
// Returns an error if any set is nil or if name is empty.
func Flatten(name string, sets ...*Set) (*FlattenResult, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("flatten: name must not be empty")
	}
	for i, s := range sets {
		if s == nil {
			return nil, fmt.Errorf("flatten: set at index %d is nil", i)
		}
	}

	out, err := NewSet(name)
	if err != nil {
		return nil, fmt.Errorf("flatten: %w", err)
	}

	result := &FlattenResult{Set: out}

	for _, s := range sets {
		keys, err := SortedKeys(s, true)
		if err != nil {
			return nil, fmt.Errorf("flatten: %w", err)
		}
		for _, k := range keys {
			v, err := s.Get(k)
			if err != nil {
				result.Skipped++
				continue
			}
			_ = out.Put(k, v)
			result.Merged++
		}
	}

	return result, nil
}

// FlattenPrefix merges multiple sets into one, prefixing each key with the
// source set's name followed by sep (e.g. "PROD_KEY").
func FlattenPrefix(name, sep string, sets ...*Set) (*FlattenResult, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("flattenPrefix: name must not be empty")
	}
	for i, s := range sets {
		if s == nil {
			return nil, fmt.Errorf("flattenPrefix: set at index %d is nil", i)
		}
	}
	if sep == "" {
		sep = "_"
	}

	out, err := NewSet(name)
	if err != nil {
		return nil, fmt.Errorf("flattenPrefix: %w", err)
	}

	result := &FlattenResult{Set: out}

	for _, s := range sets {
		prefix := strings.ToUpper(s.Name()) + sep
		keys, err := SortedKeys(s, true)
		if err != nil {
			return nil, fmt.Errorf("flattenPrefix: %w", err)
		}
		for _, k := range keys {
			v, err := s.Get(k)
			if err != nil {
				result.Skipped++
				continue
			}
			_ = out.Put(prefix+k, v)
			result.Merged++
		}
	}

	return result, nil
}
