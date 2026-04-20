package env

import (
	"strings"
	"testing"
)

// TestTransformRoundTripViaExportImport transforms a set, exports it as dotenv,
// re-imports it and verifies the transformed values survive the round-trip.
func TestTransformRoundTripViaExportImport(t *testing.T) {
	s, _ := NewSet("rt")
	_ = s.Put("KEY", "hello")
	_ = s.Put("OTHER", "world")

	xformed, err := Transform(s, UppercaseValues())
	if err != nil {
		t.Fatalf("transform: %v", err)
	}

	data, err := Export(xformed, "dotenv")
	if err != nil {
		t.Fatalf("export: %v", err)
	}

	dst, _ := NewSet("rt")
	if err := Import(dst, data, "dotenv"); err != nil {
		t.Fatalf("import: %v", err)
	}

	for _, k := range []string{"KEY", "OTHER"} {
		v, err := dst.Get(k)
		if err != nil {
			t.Fatalf("missing key %s: %v", k, err)
		}
		if !strings.ToUpper(v) == v {
			t.Errorf("expected uppercase value for %s, got %s", k, v)
		}
	}
}
