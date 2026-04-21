package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const historyFile = "history.json"

// HistoryRecord is the persisted form of a single history entry.
type HistoryRecord struct {
	SetName   string            `json:"set_name"`
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Vars      map[string]string `json:"vars"`
}

// LoadHistory reads all persisted history records from the config directory.
func LoadHistory(dir string) ([]HistoryRecord, error) {
	path := filepath.Join(dir, historyFile)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []HistoryRecord{}, nil
	}
	if err != nil {
		return nil, err
	}
	var records []HistoryRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// SaveHistory persists history records to the config directory.
func SaveHistory(dir string, records []HistoryRecord) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, historyFile), data, 0o644)
}

// AppendHistory appends a new record and saves the updated list.
func AppendHistory(dir string, rec HistoryRecord) error {
	records, err := LoadHistory(dir)
	if err != nil {
		return err
	}
	records = append(records, rec)
	return SaveHistory(dir, records)
}

// HistoryFor returns all records for a given set name.
func HistoryFor(records []HistoryRecord, setName string) []HistoryRecord {
	var out []HistoryRecord
	for _, r := range records {
		if r.SetName == setName {
			out = append(out, r)
		}
	}
	return out
}
