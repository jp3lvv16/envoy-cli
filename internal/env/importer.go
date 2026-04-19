package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ImportFormat represents the format of the input to import.
type ImportFormat string

const (
	ImportFormatDotenv ImportFormat = "dotenv"
	ImportFormatJSON   ImportFormat = "json"
)

// Import reads environment variables from r in the given format and populates s.
func Import(s *Set, format ImportFormat, r io.Reader) error {
	if s == nil {
		return fmt.Errorf("import: set must not be nil")
	}
	switch format {
	case ImportFormatDotenv:
		return importDotenv(s, r)
	case ImportFormatJSON:
		return importJSON(s, r)
	default:
		return fmt.Errorf("import: unknown format %q", format)
	}
}

func importDotenv(s *Set, r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("import: invalid line %q", line)
		}
		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"`)
		if err := s.Put(key, val); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func importJSON(s *Set, r io.Reader) error {
	var data map[string]string
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Errorf("import: invalid JSON: %w", err)
	}
	for k, v := range data {
		if err := s.Put(k, v); err != nil {
			return err
		}
	}
	return nil
}
