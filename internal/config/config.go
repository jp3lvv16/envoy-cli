package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".envoy.json"

// EnvSet represents a named set of environment variables for a deployment target.
type EnvSet struct {
	Name      string            `json:"name"`
	Target    string            `json:"target"`
	Variables map[string]string `json:"variables"`
}

// Config holds all environment sets managed by envoy-cli.
type Config struct {
	Version string   `json:"version"`
	Sets    []EnvSet `json:"sets"`
}

// Load reads the config file from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigFile
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Version: "1"}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the config to the given path, creating directories as needed.
func Save(path string, cfg *Config) error {
	if path == "" {
		path = defaultConfigFile
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// GetSet returns the EnvSet with the given name, or nil if not found.
func (c *Config) GetSet(name string) *EnvSet {
	for i := range c.Sets {
		if c.Sets[i].Name == name {
			return &c.Sets[i]
		}
	}
	return nil
}

// AddOrUpdateSet upserts an EnvSet by name.
func (c *Config) AddOrUpdateSet(set EnvSet) {
	for i := range c.Sets {
		if c.Sets[i].Name == set.Name {
			c.Sets[i] = set
			return
		}
	}
	c.Sets = append(c.Sets, set)
}
