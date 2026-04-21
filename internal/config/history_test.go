package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadHistoryMissingFileReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	records, err := LoadHistory(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty, got %d records", len(records))
	}
}

func TestSaveAndLoadHistory(t *testing.T) {
	dir := t.TempDir()
	recs := []HistoryRecord{
		{SetName: "prod", Label: "initial", Timestamp: time.Now().UTC(), Vars: map[string]string{"A": "1"}},
		{SetName: "prod", Label: "updated", Timestamp: time.Now().UTC(), Vars: map[string]string{"A": "2"}},
	}
	if err := SaveHistory(dir, recs); err != nil {
		t.Fatalf("SaveHistory: %v", err)
	}
	loaded, err := LoadHistory(dir)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2, got %d", len(loaded))
	}
	if loaded[1].Vars["A"] != "2" {
		t.Errorf("unexpected var value: %s", loaded[1].Vars["A"])
	}
}

func TestAppendHistory(t *testing.T) {
	dir := t.TempDir()
	rec := HistoryRecord{SetName: "staging", Label: "snap", Timestamp: time.Now().UTC(), Vars: map[string]string{"X": "y"}}
	if err := AppendHistory(dir, rec); err != nil {
		t.Fatalf("AppendHistory: %v", err)
	}
	if err := AppendHistory(dir, rec); err != nil {
		t.Fatalf("AppendHistory second: %v", err)
	}
	loaded, _ := LoadHistory(dir)
	if len(loaded) != 2 {
		t.Fatalf("expected 2 records, got %d", len(loaded))
	}
}

func TestHistoryFor(t *testing.T) {
	recs := []HistoryRecord{
		{SetName: "prod", Label: "a"},
		{SetName: "staging", Label: "b"},
		{SetName: "prod", Label: "c"},
	}
	out := HistoryFor(recs, "prod")
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestSaveHistoryCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/history.json"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("file should not exist yet")
	}
	_ = SaveHistory(dir, []HistoryRecord{})
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file should exist after save: %v", err)
	}
}
