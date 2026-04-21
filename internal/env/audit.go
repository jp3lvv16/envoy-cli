package env

import (
	"errors"
	"fmt"
	"time"
)

// AuditAction describes the kind of change recorded in an audit entry.
type AuditAction string

const (
	AuditPut    AuditAction = "put"
	AuditDelete AuditAction = "delete"
	AuditImport AuditAction = "import"
	AuditClear  AuditAction = "clear"
)

// AuditEntry records a single mutation on an environment set.
type AuditEntry struct {
	Timestamp time.Time   `json:"timestamp"`
	Actor     string      `json:"actor"`
	Action    AuditAction `json:"action"`
	Key       string      `json:"key,omitempty"`
	OldValue  string      `json:"old_value,omitempty"`
	NewValue  string      `json:"new_value,omitempty"`
}

// AuditLog holds an ordered list of audit entries for a named set.
type AuditLog struct {
	SetName string       `json:"set_name"`
	Entries []AuditEntry `json:"entries"`
}

// NewAuditLog creates an empty AuditLog for the given set name.
func NewAuditLog(setName string) (*AuditLog, error) {
	if setName == "" {
		return nil, errors.New("audit: set name must not be empty")
	}
	return &AuditLog{SetName: setName, Entries: []AuditEntry{}}, nil
}

// Record appends an audit entry to the log.
func (a *AuditLog) Record(actor string, action AuditAction, key, oldVal, newVal string) error {
	if actor == "" {
		return errors.New("audit: actor must not be empty")
	}
	if action == "" {
		return errors.New("audit: action must not be empty")
	}
	a.Entries = append(a.Entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Actor:     actor,
		Action:    action,
		Key:       key,
		OldValue:  oldVal,
		NewValue:  newVal,
	})
	return nil
}

// FilterByActor returns entries that match the given actor.
func (a *AuditLog) FilterByActor(actor string) []AuditEntry {
	var out []AuditEntry
	for _, e := range a.Entries {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}

// FilterByAction returns entries that match the given action.
func (a *AuditLog) FilterByAction(action AuditAction) []AuditEntry {
	var out []AuditEntry
	for _, e := range a.Entries {
		if e.Action == action {
			out = append(out, e)
		}
	}
	return out
}

// Summary returns a human-readable summary line for an entry.
func Summary(e AuditEntry) string {
	return fmt.Sprintf("%s | %-8s | actor=%-12s | key=%s",
		e.Timestamp.Format(time.RFC3339), e.Action, e.Actor, e.Key)
}
