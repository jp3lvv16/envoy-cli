package env

import (
	"fmt"
	"sort"
)

// Set holds a named collection of environment variables.
type Set struct {
	name string
	vars map[string]string
}

// NewSet creates an empty Set with the given name.
func NewSet(name string) (*Set, error) {
	if name == "" {
		return nil, fmt.Errorf("set: name must not be empty")
	}
	return &Set{name: name, vars: make(map[string]string)}, nil
}

// Name returns the set's name.
func (s *Set) Name() string { return s.name }

// Put stores key=value. Key must not be empty.
func (s *Set) Put(key, value string) error {
	if key == "" {
		return fmt.Errorf("set: key must not be empty")
	}
	s.vars[key] = value
	return nil
}

// Get retrieves the value for key, returning an error if absent.
func (s *Set) Get(key string) (string, error) {
	v, ok := s.vars[key]
	if !ok {
		return "", fmt.Errorf("set: key %q not found", key)
	}
	return v, nil
}

// Delete removes key from the set.
func (s *Set) Delete(key string) { delete(s.vars, key) }

// Keys returns all keys in sorted order.
func (s *Set) Keys() []string {
	keys := make([]string, 0, len(s.vars))
	for k := range s.vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Len returns the number of variables in the set.
func (s *Set) Len() int { return len(s.vars) }

// All returns a shallow copy of the underlying map.
func (s *Set) All() map[string]string {
	out := make(map[string]string, len(s.vars))
	for k, v := range s.vars {
		out[k] = v
	}
	return out
}
