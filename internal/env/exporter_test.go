package env

import (
	"strings"
	"testing"
)

func makeTestSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("FOO", "bar")
	_ = s.Put("BAZ", "qux")
	return s
}

func TestExportNilSetReturnsError(t *testing.T) {
	_, err := Export(nil, FormatDotenv)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestExportUnknownFormatReturnsError(t *testing.T) {
	s := makeTestSet(t)
	_, err := Export(s, ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExportDotenv(t *testing.T) {
	s := makeTestSet(t)
	out, err := Export(s, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in dotenv output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in dotenv output, got:\n%s", out)
	}
}

func TestExportShell(t *testing.T) {
	s := makeTestSet(t)
	out, err := Export(s, FormatExport)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export FOO=") {
		t.Errorf("expected export FOO= in shell output, got:\n%s", out)
	}
}

func TestExportJSON(t *testing.T) {
	s := makeTestSet(t)
	out, err := Export(s, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "{") || !strings.Contains(out, "}") {
		t.Errorf("expected JSON braces in output, got:\n%s", out)
	}
	if !strings.Contains(out, "\"FOO\"") {
		t.Errorf("expected FOO key in JSON output, got:\n%s", out)
	}
}
