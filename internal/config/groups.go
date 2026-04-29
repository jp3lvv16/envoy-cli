package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/user/envoy-cli/internal/env"
)

const groupsFile = "groups.json"

// LoadGroups reads the group index from the config directory.
// Returns an empty index if the file does not exist.
func LoadGroups(dir string) (*env.GroupIndex, error) {
	path := filepath.Join(dir, groupsFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return env.NewGroupIndex(), nil
	}
	if err != nil {
		return nil, err
	}
	var raw map[string][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	index := env.NewGroupIndex()
	for group, sets := range raw {
		for _, s := range sets {
			if err := index.Add(group, s); err != nil {
				return nil, err
			}
		}
	}
	return index, nil
}

// SaveGroups writes the group index to the config directory.
func SaveGroups(dir string, index *env.GroupIndex) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	raw := make(map[string][]string)
	for _, g := range index.Groups() {
		members, err := index.Members(g)
		if err != nil {
			return err
		}
		raw[g] = members
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, groupsFile), data, 0o644)
}

// AddGroup adds setName to groupName in the persisted index.
func AddGroup(dir, groupName, setName string) error {
	index, err := LoadGroups(dir)
	if err != nil {
		return err
	}
	if err := index.Add(groupName, setName); err != nil {
		return err
	}
	return SaveGroups(dir, index)
}

// RemoveGroup removes setName from groupName in the persisted index.
func RemoveGroup(dir, groupName, setName string) error {
	index, err := LoadGroups(dir)
	if err != nil {
		return err
	}
	if err := index.Remove(groupName, setName); err != nil {
		return err
	}
	return SaveGroups(dir, index)
}
