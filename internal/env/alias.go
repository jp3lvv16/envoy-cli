package env

import (
	"errors"
	"fmt"
	"strings"
)

// AliasIndex maps short alias names to fully-qualified set names.
type AliasIndex struct {
	aliases map[string]string // alias -> set name
}

// NewAliasIndex returns an empty AliasIndex.
func NewAliasIndex() *AliasIndex {
	return &AliasIndex{aliases: make(map[string]string)}
}

// Add registers an alias for the given set name.
// Both alias and name must be non-empty. An alias may not shadow an existing
// alias unless it points to the same set.
func (a *AliasIndex) Add(alias, setName string) error {
	alias = strings.TrimSpace(alias)
	setName = strings.TrimSpace(setName)
	if alias == "" {
		return errors.New("alias: alias must not be empty")
	}
	if setName == "" {
		return errors.New("alias: set name must not be empty")
	}
	if existing, ok := a.aliases[alias]; ok && existing != setName {
		return fmt.Errorf("alias: %q already points to %q", alias, existing)
	}
	a.aliases[alias] = setName
	return nil
}

// Remove deletes an alias. Returns an error if the alias does not exist.
func (a *AliasIndex) Remove(alias string) error {
	if _, ok := a.aliases[alias]; !ok {
		return fmt.Errorf("alias: %q not found", alias)
	}
	delete(a.aliases, alias)
	return nil
}

// Resolve returns the set name for the given alias, or an error if not found.
func (a *AliasIndex) Resolve(alias string) (string, error) {
	if name, ok := a.aliases[alias]; ok {
		return name, nil
	}
	return "", fmt.Errorf("alias: %q not found", alias)
}

// All returns a copy of the current alias map.
func (a *AliasIndex) All() map[string]string {
	copy := make(map[string]string, len(a.aliases))
	for k, v := range a.aliases {
		copy[k] = v
	}
	return copy
}
