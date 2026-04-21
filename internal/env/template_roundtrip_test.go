package env

import (
	"strings"
	"testing"
)

// TestTemplateRoundTripViaImport imports a dotenv file and then renders a
// template using the imported values, verifying end-to-end correctness.
func TestTemplateRoundTripViaImport(t *testing.T) {
	dotenv := "APP_HOST=example.com\nAPP_PORT=443\nAPP_SCHEME=https\n"

	s, err := NewSet("imported")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}

	if err := Import(s, strings.NewReader(dotenv), "dotenv"); err != nil {
		t.Fatalf("Import: %v", err)
	}

	tmpl := "{{APP_SCHEME}}://{{APP_HOST}}:{{APP_PORT}}/api"
	out, err := RenderStrict(s, tmpl)
	if err != nil {
		t.Fatalf("RenderStrict: %v", err)
	}

	want := "https://example.com:443/api"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}
