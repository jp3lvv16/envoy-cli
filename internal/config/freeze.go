package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const freezeFile = "freeze.json"

type freezeData struct {
	Frozen []string `json:"frozen"`
}

// LoadFreezes reads the frozen-set list from the config directory.
// Returns an empty slice if the file does not exist.
func LoadFreezes(dir string) ([]string, error) {
	path := filepath.Join(dir, freezeFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	var fd freezeData
	if err := json.Unmarshal(data, &fd); err != nil {
		return nil, err
	}
	return fd.Frozen, nil
}

// SaveFreezes writes the frozen-set list to the config directory.
func SaveFreezes(dir string, names []string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	fd := freezeData{Frozen: names}
	data, err := json.MarshalIndent(fd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, freezeFile), data, 0o644)
}

// AddFreeze appends name to the frozen list if not already present.
func AddFreeze(dir, name string) error {
	if name == "" {
		return errors.New("freeze: name must not be empty")
	}
	names, err := LoadFreezes(dir)
	if err != nil {
		return err
	}
	for _, n := range names {
		if n == name {
			return nil // idempotent
		}
	}
	return SaveFreezes(dir, append(names, name))
}

// RemoveFreeze removes name from the frozen list.
func RemoveFreeze(dir, name string) error {
	if name == "" {
		return errors.New("unfreeze: name must not be empty")
	}
	names, err := LoadFreezes(dir)
	if err != nil {
		return err
	}
	filtered := names[:0]
	for _, n := range names {
		if n != name {
			filtered = append(filtered, n)
		}
	}
	if len(filtered) == len(names) {
		return errors.New("unfreeze: set not found in freeze list")
	}
	return SaveFreezes(dir, filtered)
}
