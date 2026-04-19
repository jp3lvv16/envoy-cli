package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const tagsFileName = "tags.json"

// tagStore is the serialised form of a TagIndex.
type tagStore map[string][]string

// LoadTags reads the tag index from the config directory.
// Returns an empty index if the file does not exist.
func LoadTags(dir string) (tagStore, error) {
	path := filepath.Join(dir, tagsFileName)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return make(tagStore), nil
	}
	if err != nil {
		return nil, err
	}
	var store tagStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, err
	}
	return store, nil
}

// SaveTags persists the tag index to the config directory.
func SaveTags(dir string, store tagStore) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, tagsFileName), data, 0o644)
}

// AddTag adds setName to tag in the store.
func AddTag(store tagStore, tag, setName string) error {
	if tag == "" {
		return errors.New("tag must not be empty")
	}
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	for _, s := range store[tag] {
		if s == setName {
			return nil
		}
	}
	store[tag] = append(store[tag], setName)
	return nil
}

// RemoveTag removes setName from tag in the store.
func RemoveTag(store tagStore, tag, setName string) error {
	sets, ok := store[tag]
	if !ok {
		return errors.New("tag not found: " + tag)
	}
	new := sets[:0]
	for _, s := range sets {
		if s != setName {
			new = append(new, s)
		}
	}
	if len(new) == 0 {
		delete(store, tag)
	} else {
		store[tag] = new
	}
	return nil
}
