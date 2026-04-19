package env

import (
	"fmt"
	"time"
)

// Snapshot captures the state of a Set at a point in time.
type Snapshot struct {
	SetName   string
	TakenAt   time.Time
	Vars      map[string]string
}

// TakeSnapshot creates a Snapshot from the given Set.
func TakeSnapshot(s *Set) (*Snapshot, error) {
	if s == nil {
		return nil, fmt.Errorf("snapshot: set must not be nil")
	}
	vars := make(map[string]string)
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		vars[k] = v
	}
	return &Snapshot{
		SetName: s.Name(),
		TakenAt: time.Now().UTC(),
		Vars:    vars,
	}, nil
}

// RestoreSnapshot applies a Snapshot's variables to the destination Set,
// overwriting any existing keys.
func RestoreSnapshot(snap *Snapshot, dst *Set) error {
	if snap == nil {
		return fmt.Errorf("snapshot: snapshot must not be nil")
	}
	if dst == nil {
		return fmt.Errorf("snapshot: destination set must not be nil")
	}
	for k, v := range snap.Vars {
		if err := dst.Put(k, v); err != nil {
			return fmt.Errorf("snapshot: restore failed on key %q: %w", k, err)
		}
	}
	return nil
}
