package env

import (
	"errors"
	"fmt"
)

// Set represents a named collection of environment variables.
type Set struct {
	Name   string            `json:"name"`
	Vars   map[string]string `json:"vars"`
}

// NewSet creates a new empty Set with the given name.
func NewSet(name string) (*Set, error) {
	if name == "" {
		return nil, errors.New("set name cannot be empty")
	}
	return &Set{
		Name: name,
		Vars: make(map[string]string),
	}, nil
}

// Put adds or updates a key-value pair in the set.
func (s *Set) Put(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	s.Vars[key] = value
	return nil
}

// Get retrieves a value by key. Returns an error if the key does not exist.
func (s *Set) Get(key string) (string, error) {
	v, ok := s.Vars[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in set %q", key, s.Name)
	}
	return v, nil
}

// Delete removes a key from the set. Returns an error if the key does not exist.
func (s *Set) Delete(key string) error {
	if _, ok := s.Vars[key]; !ok {
		return fmt.Errorf("key %q not found in set %q", key, s.Name)
	}
	delete(s.Vars, key)
	return nil
}

// List returns all key-value pairs as a slice of formatted strings.
func (s *Set) List() []string {
	out := make([]string, 0, len(s.Vars))
	for k, v := range s.Vars {
		out = append(out, fmt.Sprintf("%s=%s", k, v))
	}
	return out
}
