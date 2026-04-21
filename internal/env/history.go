package env

import (
	"errors"
	"time"
)

// HistoryEntry records a snapshot of a Set at a point in time.
type HistoryEntry struct {
	Timestamp time.Time
	Label     string
	Snapshot  map[string]string
}

// History maintains an ordered list of HistoryEntry values for a named set.
type History struct {
	SetName string
	entries []HistoryEntry
}

// NewHistory creates a History for the given set name.
func NewHistory(setName string) (*History, error) {
	if setName == "" {
		return nil, errors.New("history: set name must not be empty")
	}
	return &History{SetName: setName}, nil
}

// Record captures the current state of s into the history with an optional label.
func (h *History) Record(s *Set, label string) error {
	if s == nil {
		return errors.New("history: set must not be nil")
	}
	snap, err := TakeSnapshot(s)
	if err != nil {
		return err
	}
	h.entries = append(h.entries, HistoryEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Snapshot:  snap.Vars,
	})
	return nil
}

// Entries returns a copy of all recorded history entries.
func (h *History) Entries() []HistoryEntry {
	out := make([]HistoryEntry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the number of recorded entries.
func (h *History) Len() int { return len(h.entries) }

// At returns the HistoryEntry at index i.
func (h *History) At(i int) (HistoryEntry, error) {
	if i < 0 || i >= len(h.entries) {
		return HistoryEntry{}, errors.New("history: index out of range")
	}
	return h.entries[i], nil
}

// Rollback restores s to the state captured at index i.
func (h *History) Rollback(s *Set, i int) error {
	entry, err := h.At(i)
	if err != nil {
		return err
	}
	snap := &Snapshot{SetName: h.SetName, Vars: entry.Snapshot}
	return RestoreSnapshot(s, snap)
}
