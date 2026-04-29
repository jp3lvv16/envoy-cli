package config

import (
	"os"
	"testing"
)

func TestLoadGroupsMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	index, err := LoadGroups(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(index.Groups()) != 0 {
		t.Fatal("expected empty group index")
	}
}

func TestSaveAndLoadGroups(t *testing.T) {
	dir := t.TempDir()
	index, _ := LoadGroups(dir)
	_ = index.Add("prod", "api")
	_ = index.Add("prod", "worker")
	if err := SaveGroups(dir, index); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := LoadGroups(dir)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	members, _ := loaded.Members("prod")
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestAddGroupEmptyGroupReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := AddGroup(dir, "", "api"); err == nil {
		t.Fatal("expected error for empty group name")
	}
}

func TestAddGroupDuplicateIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	_ = AddGroup(dir, "staging", "api")
	_ = AddGroup(dir, "staging", "api")
	index, _ := LoadGroups(dir)
	members, _ := index.Members("staging")
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}
}

func TestRemoveGroupDeletesEmptyGroup(t *testing.T) {
	dir := t.TempDir()
	_ = AddGroup(dir, "dev", "only-set")
	_ = RemoveGroup(dir, "dev", "only-set")
	index, _ := LoadGroups(dir)
	if len(index.Groups()) != 0 {
		t.Fatal("expected no groups after removing last member")
	}
}

func TestRemoveGroupNotFoundReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := RemoveGroup(dir, "ghost", "api"); err == nil {
		t.Fatal("expected error for missing group")
	}
}

func TestSaveGroupsCreatesFile(t *testing.T) {
	dir := t.TempDir()
	index, _ := LoadGroups(dir)
	_ = index.Add("ci", "build")
	_ = SaveGroups(dir, index)
	path := dir + "/groups.json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected groups.json to be created")
	}
}
