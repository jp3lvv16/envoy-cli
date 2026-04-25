package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	cfg, err := Load("/tmp/envoy_nonexistent_12345.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("expected version '1', got %q", cfg.Version)
	}
	if len(cfg.Sets) != 0 {
		t.Errorf("expected empty sets, got %d", len(cfg.Sets))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envoy.json")

	cfg := &Config{
		Version: "1",
		Sets: []EnvSet{
			{
				Name:   "production",
				Target: "prod",
				Variables: map[string]string{"DB_HOST": "prod.db", "PORT": "5432"},
			},
		},
	}

	if err := Save(path, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Sets) != 1 {
		t.Fatalf("expected 1 set, got %d", len(loaded.Sets))
	}
	if loaded.Sets[0].Variables["DB_HOST"] != "prod.db" {
		t.Errorf("unexpected DB_HOST value: %q", loaded.Sets[0].Variables["DB_HOST"])
	}
}

func TestGetSet(t *testing.T) {
	cfg := &Config{Sets: []EnvSet{{Name: "staging", Target: "stg"}}}
	if s := cfg.GetSet("staging"); s == nil {
		t.Error("expected to find 'staging' set")
	}
	if s := cfg.GetSet("missing"); s != nil {
		t.Error("expected nil for missing set")
	}
}

func TestAddOrUpdateSet(t *testing.T) {
	cfg := &Config{}
	cfg.AddOrUpdateSet(EnvSet{Name: "dev", Target: "development", Variables: map[string]string{"DEBUG": "true"}})
	if len(cfg.Sets) != 1 {
		t.Fatalf("expected 1 set after add, got %d", len(cfg.Sets))
	}
	cfg.AddOrUpdateSet(EnvSet{Name: "dev", Target: "development", Variables: map[string]string{"DEBUG": "false"}})
	if len(cfg.Sets) != 1 {
		t.Fatalf("expected 1 set after update, got %d", len(cfg.Sets))
	}
	if cfg.Sets[0].Variables["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false after update")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", ".envoy.json")
	if err := Save(path, &Config{Version: "1"}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestRemoveSet(t *testing.T) {
	cfg := &Config{
		Sets: []EnvSet{
			{Name: "dev", Target: "development"},
			{Name: "prod", Target: "production"},
		},
	}
	cfg.RemoveSet("dev")
	if len(cfg.Sets) != 1 {
		t.Fatalf("expected 1 set after remove, got %d", len(cfg.Sets))
	}
	if cfg.Sets[0].Name != "prod" {
		t.Errorf("expected remaining set to be 'prod', got %q", cfg.Sets[0].Name)
	}
	// Removing a non-existent set should be a no-op
	cfg.RemoveSet("missing")
	if len(cfg.Sets) != 1 {
		t.Errorf("expected set count unchanged after removing missing set, got %d", len(cfg.Sets))
	}
}
