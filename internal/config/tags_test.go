package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTagsMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	store, err := LoadTags(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store) != 0 {
		t.Fatal("expected empty store")
	}
}

func TestSaveAndLoadTags(t *testing.T) {
	dir := t.TempDir()
	store := make(tagStore)
	_ = AddTag(store, "prod", "api")
	_ = AddTag(store, "prod", "db")
	if err := SaveTags(dir, store); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadTags(dir)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if len(loaded["prod"]) != 2 {
		t.Fatalf("expected 2 sets, got %d", len(loaded["prod"]))
	}
}

func TestAddTagEmptyTagReturnsError(t *testing.T) {
	store := make(tagStore)
	if err := AddTag(store, "", "set"); err == nil {
		t.Fatal("expected error")
	}
}

func TestAddTagDuplicateIsIdempotent(t *testing.T) {
	store := make(tagStore)
	_ = AddTag(store, "dev", "web")
	_ = AddTag(store, "dev", "web")
	if len(store["dev"]) != 1 {
		t.Fatalf("expected 1, got %d", len(store["dev"]))
	}
}

func TestRemoveTagDeletesEmptyTag(t *testing.T) {
	store := make(tagStore)
	_ = AddTag(store, "staging", "only")
	_ = RemoveTag(store, "staging", "only")
	if _, ok := store["staging"]; ok {
		t.Fatal("expected tag to be removed")
	}
}

func TestRemoveTagMissingReturnsError(t *testing.T) {
	store := make(tagStore)
	if err := RemoveTag(store, "ghost", "set"); err == nil {
		t.Fatal("expected error")
	}
}

func TestSaveTagsCreatesFile(t *testing.T) {
	dir := t.TempDir()
	store := make(tagStore)
	_ = AddTag(store, "x", "y")
	_ = SaveTags(dir, store)
	if _, err := os.Stat(filepath.Join(dir, tagsFileName)); err != nil {
		t.Fatal("expected file to exist")
	}
}
