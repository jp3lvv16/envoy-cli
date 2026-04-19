package env

import (
	"bytes"
	"testing"
)

// TestRoundTripDotenv exports a set to dotenv format and re-imports it,
// verifying that all key-value pairs survive the round trip.
func TestRoundTripDotenv(t *testing.T) {
	original, _ := NewSet("rt")
	_ = original.Put("ALPHA", "one")
	_ = original.Put("BETA", "two")

	var buf bytes.Buffer
	if err := Export(original, ExportFormatDotenv, &buf); err != nil {
		t.Fatalf("export error: %v", err)
	}

	restored, _ := NewSet("rt")
	if err := Import(restored, ImportFormatDotenv, &buf); err != nil {
		t.Fatalf("import error: %v", err)
	}

	for _, key := range []string{"ALPHA", "BETA"} {
		ov, _ := original.Get(key)
		rv, _ := restored.Get(key)
		if ov != rv {
			t.Errorf("key %s: original=%q restored=%q", key, ov, rv)
		}
	}
}

// TestRoundTripJSON exports a set to JSON format and re-imports it.
func TestRoundTripJSON(t *testing.T) {
	original, _ := NewSet("rt")
	_ = original.Put("HOST", "localhost")
	_ = original.Put("PORT", "8080")

	var buf bytes.Buffer
	if err := Export(original, ExportFormatJSON, &buf); err != nil {
		t.Fatalf("export error: %v", err)
	}

	restored, _ := NewSet("rt")
	if err := Import(restored, ImportFormatJSON, &buf); err != nil {
		t.Fatalf("import error: %v", err)
	}

	for _, key := range []string{"HOST", "PORT"} {
		ov, _ := original.Get(key)
		rv, _ := restored.Get(key)
		if ov != rv {
			t.Errorf("key %s: original=%q restored=%q", key, ov, rv)
		}
	}
}
