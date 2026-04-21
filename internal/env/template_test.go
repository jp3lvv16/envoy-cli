package env

import (
	"testing"
)

func makeTemplateSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("tpl")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("HOST", "localhost")
	_ = s.Put("PORT", "8080")
	_ = s.Put("ENV", "production")
	return s
}

func TestRenderNilSetReturnsError(t *testing.T) {
	_, err := Render(nil, "hello")
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestRenderEmptyTemplateReturnsEmpty(t *testing.T) {
	s := makeTemplateSet(t)
	res, err := Render(s, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "" {
		t.Errorf("expected empty output, got %q", res.Output)
	}
}

func TestRenderResolvesPlaceholders(t *testing.T) {
	s := makeTemplateSet(t)
	res, err := Render(s, "http://{{HOST}}:{{PORT}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "http://localhost:8080" {
		t.Errorf("unexpected output: %q", res.Output)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
}

func TestRenderMissingKeyLeftInPlace(t *testing.T) {
	s := makeTemplateSet(t)
	res, err := Render(s, "{{HOST}}:{{UNDEFINED}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "localhost:{{UNDEFINED}}" {
		t.Errorf("unexpected output: %q", res.Output)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "UNDEFINED" {
		t.Errorf("expected missing=[UNDEFINED], got %v", res.Missing)
	}
}

func TestRenderDuplicateMissingKeyReportedOnce(t *testing.T) {
	s := makeTemplateSet(t)
	res, err := Render(s, "{{MISSING}} and {{MISSING}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 missing entry, got %d", len(res.Missing))
	}
}

func TestRenderStrictReturnsErrorOnMissing(t *testing.T) {
	s := makeTemplateSet(t)
	_, err := RenderStrict(s, "{{HOST}}:{{MISSING}}")
	if err == nil {
		t.Fatal("expected error for unresolved placeholder")
	}
}

func TestRenderStrictSucceedsWhenAllResolved(t *testing.T) {
	s := makeTemplateSet(t)
	out, err := RenderStrict(s, "{{ENV}}-{{HOST}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "production-localhost" {
		t.Errorf("unexpected output: %q", out)
	}
}
