package env

import "testing"

// TestCloneRoundTripViaExportImport clones a set, exports it as dotenv,
// re-imports it and verifies all variables survive the round-trip.
func TestCloneRoundTripViaExportImport(t *testing.T) {
	src, _ := NewSet("base")
	_ = src.Put("APP_ENV", "production")
	_ = src.Put("DB_URL", "postgres://localhost/mydb")

	cloned, err := Clone(src, "base-clone")
	if err != nil {
		t.Fatalf("Clone: %v", err)
	}

	data, err := Export(cloned, "dotenv")
	if err != nil {
		t.Fatalf("Export: %v", err)
	}

	reimported, _ := NewSet("reimported")
	if err := Import(reimported, data, "dotenv"); err != nil {
		t.Fatalf("Import: %v", err)
	}

	for k, v := range src.Vars() {
		got, err := reimported.Get(k)
		if err != nil {
			t.Errorf("missing key %q after round-trip", k)
			continue
		}
		if got != v {
			t.Errorf("key %q: want %q, got %q", k, v, got)
		}
	}
}
