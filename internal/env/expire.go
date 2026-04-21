package env

import (
	"errors"
	"fmt"
	"time"
)

// Expiry holds an expiration timestamp for a named env set.
type Expiry struct {
	SetName   string    `json:"set_name"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ExpiryIndex maps set names to their expiry records.
type ExpiryIndex struct {
	entries map[string]Expiry
}

// NewExpiryIndex returns an initialised ExpiryIndex.
func NewExpiryIndex() *ExpiryIndex {
	return &ExpiryIndex{entries: make(map[string]Expiry)}
}

// Set registers or updates the expiry for the given set name.
func (x *ExpiryIndex) Set(setName string, ttl time.Duration) error {
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	if ttl <= 0 {
		return errors.New("ttl must be positive")
	}
	x.entries[setName] = Expiry{
		SetName:   setName,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}
	return nil
}

// Remove deletes the expiry entry for the given set name.
func (x *ExpiryIndex) Remove(setName string) error {
	if _, ok := x.entries[setName]; !ok {
		return fmt.Errorf("no expiry found for set %q", setName)
	}
	delete(x.entries, setName)
	return nil
}

// IsExpired reports whether the named set has passed its expiry time.
// Returns false if no expiry is registered.
func (x *ExpiryIndex) IsExpired(setName string) bool {
	e, ok := x.entries[setName]
	if !ok {
		return false
	}
	return time.Now().UTC().After(e.ExpiresAt)
}

// Get returns the Expiry for the given set name and whether it exists.
func (x *ExpiryIndex) Get(setName string) (Expiry, bool) {
	e, ok := x.entries[setName]
	return e, ok
}

// All returns a copy of all registered expiry entries.
func (x *ExpiryIndex) All() []Expiry {
	out := make([]Expiry, 0, len(x.entries))
	for _, e := range x.entries {
		out = append(out, e)
	}
	return out
}
