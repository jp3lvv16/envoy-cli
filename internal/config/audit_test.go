package config

import (
	"os"
	"testing"

	"github.com/user/envoy-cli/internal/env"
)

func TestLoadAuditMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	store, err := LoadAudit(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store) != 0 {
		t.Fatalf("expected empty store, got %d entries", len(store))
	}
}

func TestSaveAndLoadAudit(t *testing.T) {
	dir := t.TempDir()
	log, _ := env.NewAuditLog("staging")
	_ = log.Record("alice", env.AuditPut, "KEY", "", "val")
	store := auditStore{"staging": log}

	if err := SaveAudit(dir, store); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadAudit(dir)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, ok := loaded["staging"]; !ok {
		t.Fatal("expected staging log to be present")
	}
}

func TestAppendAuditCreatesEntry(t *testing.T) {
	dir := t.TempDir()
	err := AppendAudit(dir, "prod", "bob", env.AuditDelete, "SECRET", "old", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	log, err := AuditFor(dir, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log == nil || len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %v", log)
	}
	if log.Entries[0].Actor != "bob" {
		t.Errorf("expected actor bob, got %s", log.Entries[0].Actor)
	}
}

func TestAppendAuditAccumulatesEntries(t *testing.T) {
	dir := t.TempDir()
	_ = AppendAudit(dir, "dev", "alice", env.AuditPut, "A", "", "1")
	_ = AppendAudit(dir, "dev", "alice", env.AuditPut, "B", "", "2")

	log, err := AuditFor(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
}

func TestAuditForMissingSetReturnsNil(t *testing.T) {
	dir := t.TempDir()
	log, err := AuditFor(dir, "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log != nil {
		t.Fatal("expected nil log for unknown set")
	}
}

func TestSaveAuditCreatesFile(t *testing.T) {
	dir := t.TempDir()
	store := auditStore{}
	if err := SaveAudit(dir, store); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir + "/audit.json"); err != nil {
		t.Fatalf("expected audit.json to exist: %v", err)
	}
}
