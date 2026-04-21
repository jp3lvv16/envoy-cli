package config

import (
	"os"
	"testing"
)

func TestLoadLocksMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	records, err := LoadLocks(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(records))
	}
}

func TestSaveAndLoadLocks(t *testing.T) {
	dir := t.TempDir()
	if err := AddLock(dir, "prod", "alice"); err != nil {
		t.Fatalf("AddLock: %v", err)
	}
	records, err := LoadLocks(dir)
	if err != nil {
		t.Fatalf("LoadLocks: %v", err)
	}
	r, ok := records["prod"]
	if !ok {
		t.Fatal("expected prod lock to exist")
	}
	if r.LockedBy != "alice" {
		t.Fatalf("expected alice, got %s", r.LockedBy)
	}
	if r.LockedAt.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestAddLockDuplicateReturnsError(t *testing.T) {
	dir := t.TempDir()
	_ = AddLock(dir, "prod", "alice")
	if err := AddLock(dir, "prod", "bob"); err == nil {
		t.Fatal("expected error on duplicate lock")
	}
}

func TestRemoveLockDeletesEntry(t *testing.T) {
	dir := t.TempDir()
	_ = AddLock(dir, "prod", "alice")
	if err := RemoveLock(dir, "prod"); err != nil {
		t.Fatalf("RemoveLock: %v", err)
	}
	records, _ := LoadLocks(dir)
	if _, ok := records["prod"]; ok {
		t.Fatal("expected lock to be removed")
	}
}

func TestRemoveLockNotLockedReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := RemoveLock(dir, "prod"); err == nil {
		t.Fatal("expected error when removing non-existent lock")
	}
}

func TestSaveLocksCreatesFile(t *testing.T) {
	dir := t.TempDir()
	_ = AddLock(dir, "staging", "ci")
	if _, err := os.Stat(dir + "/locks.json"); err != nil {
		t.Fatalf("expected locks.json to exist: %v", err)
	}
}
