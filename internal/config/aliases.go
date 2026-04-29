package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const aliasFile = "aliases.json"

// LoadAliases reads the alias map from the config directory.
// Returns an empty map when the file does not exist yet.
func LoadAliases(dir string) (map[string]string, error) {
	path := filepath.Join(dir, aliasFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("aliases: read %s: %w", path, err)
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("aliases: parse %s: %w", path, err)
	}
	if m == nil {
		m = map[string]string{}
	}
	return m, nil
}

// SaveAliases persists the alias map to the config directory.
func SaveAliases(dir string, aliases map[string]string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("aliases: mkdir %s: %w", dir, err)
	}
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return fmt.Errorf("aliases: marshal: %w", err)
	}
	path := filepath.Join(dir, aliasFile)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("aliases: write %s: %w", path, err)
	}
	return nil
}

// AddAlias adds or validates an alias entry and saves the updated map.
func AddAlias(dir, alias, setName string) error {
	if alias == "" {
		return errors.New("aliases: alias must not be empty")
	}
	if setName == "" {
		return errors.New("aliases: set name must not be empty")
	}
	m, err := LoadAliases(dir)
	if err != nil {
		return err
	}
	if existing, ok := m[alias]; ok && existing != setName {
		return fmt.Errorf("aliases: %q already points to %q", alias, existing)
	}
	m[alias] = setName
	return SaveAliases(dir, m)
}

// RemoveAlias deletes an alias entry and saves the updated map.
func RemoveAlias(dir, alias string) error {
	m, err := LoadAliases(dir)
	if err != nil {
		return err
	}
	if _, ok := m[alias]; !ok {
		return fmt.Errorf("aliases: %q not found", alias)
	}
	delete(m, alias)
	return SaveAliases(dir, m)
}
