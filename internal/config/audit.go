package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/envoy-cli/internal/env"
)

const auditFile = "audit.json"

type auditStore map[string]*env.AuditLog

// LoadAudit reads the audit store from disk. Returns an empty store if the
// file does not exist yet.
func LoadAudit(dir string) (auditStore, error) {
	path := filepath.Join(dir, auditFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return auditStore{}, nil
	}
	if err != nil {
		return nil, err
	}
	var store auditStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// SaveAudit persists the audit store to disk.
func SaveAudit(dir string, store auditStore) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, auditFile), data, 0o644)
}

// AppendAudit records an audit entry for the given set name and persists it.
func AppendAudit(dir, setName, actor string, action env.AuditAction, key, oldVal, newVal string) error {
	store, err := LoadAudit(dir)
	if err != nil {
		return err
	}
	log, ok := store[setName]
	if !ok {
		log, err = env.NewAuditLog(setName)
		if err != nil {
			return err
		}
		store[setName] = log
	}
	if err := log.Record(actor, action, key, oldVal, newVal); err != nil {
		return err
	}
	return SaveAudit(dir, store)
}

// AuditFor returns the audit log for a specific set, or nil if none exists.
func AuditFor(dir, setName string) (*env.AuditLog, error) {
	store, err := LoadAudit(dir)
	if err != nil {
		return nil, err
	}
	log, ok := store[setName]
	if !ok {
		return nil, nil
	}
	return log, nil
}
