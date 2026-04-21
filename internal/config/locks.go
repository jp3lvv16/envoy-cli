package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const locksFile = "locks.json"

// LockRecord is the persistent form of a lock entry.
type LockRecord struct {
	LockedAt time.Time `json:"locked_at"`
	LockedBy string    `json:"locked_by"`
}

// LoadLocks reads lock records from the config directory.
// Returns an empty map if the file does not exist.
func LoadLocks(dir string) (map[string]LockRecord, error) {
	path := filepath.Join(dir, locksFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(map[string]LockRecord), nil
	}
	if err != nil {
		return nil, err
	}
	var records map[string]LockRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// SaveLocks writes lock records to the config directory.
func SaveLocks(dir string, records map[string]LockRecord) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, locksFile), data, 0o644)
}

// AddLock persists a lock entry for the named set.
func AddLock(dir, setName, owner string) error {
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	records, err := LoadLocks(dir)
	if err != nil {
		return err
	}
	if _, exists := records[setName]; exists {
		return errors.New("set is already locked")
	}
	records[setName] = LockRecord{LockedAt: time.Now().UTC(), LockedBy: owner}
	return SaveLocks(dir, records)
}

// RemoveLock removes the persisted lock for the named set.
func RemoveLock(dir, setName string) error {
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	records, err := LoadLocks(dir)
	if err != nil {
		return err
	}
	if _, exists := records[setName]; !exists {
		return errors.New("set is not locked")
	}
	delete(records, setName)
	return SaveLocks(dir, records)
}
