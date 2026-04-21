package env

import (
	"strings"
	"testing"
)

func TestFormatCompareNilReturnsEmpty(t *testing.T) {
	out := FormatCompare(nil)
	if out != "" {
		t.Errorf("expected empty string for nil, got %q", out)
	}
}

func TestFormatCompareShowsSameSection(t *testing.T) {
	cr := &CompareResult{
		Same:       map[string]string{"HOST": "localhost"},
		OnlyInSrc:  map[string]string{},
		OnlyInDst:  map[string]string{},
		Conflicted: map[string][2]string{},
	}
	out := FormatCompare(cr)
	if !strings.Contains(out, "Same") {
		t.Error("expected 'Same' section in output")
	}
	if !strings.Contains(out, "HOST=localhost") {
		t.Error("expected HOST=localhost in output")
	}
}

func TestFormatCompareShowsOnlyInSrc(t *testing.T) {
	cr := &CompareResult{
		Same:       map[string]string{},
		OnlyInSrc:  map[string]string{"EXTRA": "val"},
		OnlyInDst:  map[string]string{},
		Conflicted: map[string][2]string{},
	}
	out := FormatCompare(cr)
	if !strings.Contains(out, "Only in source") {
		t.Error("expected 'Only in source' section")
	}
	if !strings.Contains(out, "+ EXTRA=val") {
		t.Error("expected '+ EXTRA=val' in output")
	}
}

func TestFormatCompareShowsOnlyInDst(t *testing.T) {
	cr := &CompareResult{
		Same:       map[string]string{},
		OnlyInSrc:  map[string]string{},
		OnlyInDst:  map[string]string{"NEW": "42"},
		Conflicted: map[string][2]string{},
	}
	out := FormatCompare(cr)
	if !strings.Contains(out, "Only in destination") {
		t.Error("expected 'Only in destination' section")
	}
	if !strings.Contains(out, "- NEW=42") {
		t.Error("expected '- NEW=42' in output")
	}
}

func TestFormatCompareShowsConflicted(t *testing.T) {
	cr := &CompareResult{
		Same:       map[string]string{},
		OnlyInSrc:  map[string]string{},
		OnlyInDst:  map[string]string{},
		Conflicted: map[string][2]string{"PORT": {"8080", "9090"}},
	}
	out := FormatCompare(cr)
	if !strings.Contains(out, "Conflicted") {
		t.Error("expected 'Conflicted' section")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected PORT in conflicted section")
	}
	if !strings.Contains(out, "8080") || !strings.Contains(out, "9090") {
		t.Error("expected both conflict values in output")
	}
}

func TestFormatCompareEmptyResultIsBlank(t *testing.T) {
	cr := &CompareResult{
		Same:       map[string]string{},
		OnlyInSrc:  map[string]string{},
		OnlyInDst:  map[string]string{},
		Conflicted: map[string][2]string{},
	}
	out := FormatCompare(cr)
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected blank output for empty result, got %q", out)
	}
}
