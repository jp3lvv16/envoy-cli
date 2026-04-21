package env

import (
	"errors"
	"sync"
	"time"
)

// LockEntry records when a set was locked and by whom.
type LockEntry struct {
	LockedAt time.Time
	LockedBy string
}

// LockIndex tracks which env sets are locked.
type LockIndex struct {
	mu      sync.RWMutex
	entries map[string]LockEntry
}

// NewLockIndex returns an empty LockIndex.
func NewLockIndex() *LockIndex {
	return &LockIndex{entries: make(map[string]LockEntry)}
}

// Lock marks the named set as locked by the given owner.
// Returns an error if the set is already locked.
func (l *LockIndex) Lock(setName, owner string) error {
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	if owner == "" {
		return errors.New("owner must not be empty")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.entries[setName]; exists {
		return errors.New("set is already locked")
	}
	l.entries[setName] = LockEntry{LockedAt: time.Now().UTC(), LockedBy: owner}
	return nil
}

// Unlock removes the lock on the named set.
// Returns an error if the set is not locked.
func (l *LockIndex) Unlock(setName string) error {
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.entries[setName]; !exists {
		return errors.New("set is not locked")
	}
	delete(l.entries, setName)
	return nil
}

// IsLocked reports whether the named set is currently locked.
func (l *LockIndex) IsLocked(setName string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, exists := l.entries[setName]
	return exists
}

// GetLock returns the LockEntry for the named set and true if locked,
// or a zero LockEntry and false if not locked.
func (l *LockIndex) GetLock(setName string) (LockEntry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.entries[setName]
	return e, ok
}

// All returns a copy of all current lock entries keyed by set name.
func (l *LockIndex) All() map[string]LockEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	copy := make(map[string]LockEntry, len(l.entries))
	for k, v := range l.entries {
		copy[k] = v
	}
	return copy
}
