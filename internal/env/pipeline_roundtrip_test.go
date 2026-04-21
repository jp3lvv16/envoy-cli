package env

import (
	"strings"
	"testing"
)

// TestPipelineRoundTripViaExportImport runs a pipeline that transforms a set,
// exports it as dotenv, re-imports it, and verifies the values survived.
func TestPipelineRoundTripViaExportImport(t *testing.T) {
	src, err := NewSet("roundtrip")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = src.Put("api_key", "secret")
	_ = src.Put("base_url", "http://example.com")

	p, _ := NewPipeline("rtrip")
	_ = p.AddStep(func(s *Set) (*Set, error) {
		return Transform(s, func(_, v string) string { return strings.ToUpper(v) })
	})

	transformed, err := p.Run(src)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	data, err := Export(transformed, "dotenv")
	if err != nil {
		t.Fatalf("Export: %v", err)
	}

	dst, err := NewSet("roundtrip")
	if err != nil {
		t.Fatalf("NewSet dst: %v", err)
	}
	if err := Import(dst, data, "dotenv"); err != nil {
		t.Fatalf("Import: %v", err)
	}

	v, err := dst.Get("api_key")
	if err != nil {
		t.Fatalf("Get api_key: %v", err)
	}
	if v != "SECRET" {
		t.Fatalf("expected SECRET, got %q", v)
	}
}
