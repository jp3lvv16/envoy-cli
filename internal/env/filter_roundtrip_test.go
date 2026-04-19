package env

import (
	"strings"
	"testing"
)

// TestFilterRoundTripViaExportImport filters a set, exports it as dotenv,
// re-imports it and verifies the values survive the round-trip.
func TestFilterRoundTripViaExportImport(t *testing.T) {
	src := makeFilterSet(t)

	filtered, err := FilterByPrefix(src, "APP_")
	if err != nil {
		t.Fatalf("filter: %v", err)
	}

	raw, err := Export(filtered, "dotenv")
	if err != nil {
		t.Fatalf("export: %v", err)
	}

	dst, err := NewSet("imported")
	if err != nil {
		t.Fatalf("new set: %v", err)
	}

	if err := Import(dst, strings.NewReader(raw), "dotenv"); err != nil {
		t.Fatalf("import: %v", err)
	}

	for _, k := range filtered.Keys() {
		want, _ := filtered.Get(k)
		got, err := dst.Get(k)
		if err != nil {
			t.Errorf("missing key %q after round-trip", k)
			continue
		}
		if got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}
