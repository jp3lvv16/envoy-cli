package env

import (
	"strings"
	"testing"
)

func TestImportNilSetReturnsError(t *testing.T) {
	err := Import(nil, ImportFormatDotenv, strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestImportUnknownFormatReturnsError(t *testing.T) {
	s, _ := NewSet("test")
	err := Import(s, ImportFormat("xml"), strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestImportDotenv(t *testing.T) {
	input := `# comment
FOO=bar
BAZ="hello world"
`
	s, _ := NewSet("test")
	if err := Import(s, ImportFormatDotenv, strings.NewReader(input)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := s.Get("FOO"); v != "bar" {
		t.Errorf("expected bar, got %q", v)
	}
	if v, _ := s.Get("BAZ"); v != "hello world" {
		t.Errorf("expected 'hello world', got %q", v)
	}
}

func TestImportDotenvInvalidLine(t *testing.T) {
	s, _ := NewSet("test")
	err := Import(s, ImportFormatDotenv, strings.NewReader("INVALID_LINE"))
	if err == nil {
		t.Fatal("expected error for invalid line")
	}
}

func TestImportJSON(t *testing.T) {
	input := `{"KEY":"value","NUM":"42"}`
	s, _ := NewSet("test")
	if err := Import(s, ImportFormatJSON, strings.NewReader(input)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := s.Get("KEY"); v != "value" {
		t.Errorf("expected value, got %q", v)
	}
	if v, _ := s.Get("NUM"); v != "42" {
		t.Errorf("expected 42, got %q", v)
	}
}

func TestImportJSONInvalid(t *testing.T) {
	s, _ := NewSet("test")
	err := Import(s, ImportFormatJSON, strings.NewReader("not-json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
