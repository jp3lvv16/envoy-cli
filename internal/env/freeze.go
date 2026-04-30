package env

import (
	"errors"
	"fmt"
)

// FreezeIndex tracks which env sets are frozen (read-only).
type FreezeIndex struct {
	frozen map[string]bool
}

// NewFreezeIndex returns an empty FreezeIndex.
func NewFreezeIndex() *FreezeIndex {
	return &FreezeIndex{frozen: make(map[string]bool)}
}

// Freeze marks the named set as frozen.
func (f *FreezeIndex) Freeze(name string) error {
	if name == "" {
		return errors.New("freeze: set name must not be empty")
	}
	f.frozen[name] = true
	return nil
}

// Unfreeze removes the frozen status from the named set.
func (f *FreezeIndex) Unfreeze(name string) error {
	if name == "" {
		return errors.New("unfreeze: set name must not be empty")
	}
	if !f.frozen[name] {
		return fmt.Errorf("unfreeze: set %q is not frozen", name)
	}
	delete(f.frozen, name)
	return nil
}

// IsFrozen reports whether the named set is currently frozen.
func (f *FreezeIndex) IsFrozen(name string) bool {
	return f.frozen[name]
}

// List returns all currently frozen set names.
func (f *FreezeIndex) List() []string {
	names := make([]string, 0, len(f.frozen))
	for k := range f.frozen {
		names = append(names, k)
	}
	return names
}

// AssertMutable returns an error if the named set is frozen.
func (f *FreezeIndex) AssertMutable(name string) error {
	if f.IsFrozen(name) {
		return fmt.Errorf("set %q is frozen and cannot be modified", name)
	}
	return nil
}
