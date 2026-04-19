package env

import (
	"strings"
	"testing"
)

func makeValidSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	return s
}

func TestValidateNilSetReturnsError(t *testing.T) {
	if err := Validate(nil, false); err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestValidateCleanSetNoError(t *testing.T) {
	s := makeValidSet(t)
	_ = s.Put("APP_ENV", "production")
	_ = s.Put("PORT", "8080")
	if err := Validate(s, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateInvalidKeyReturnsError(t *testing.T) {
	s := makeValidSet(t)
	_ = s.Put("1INVALID", "value")
	err := Validate(s, false)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(err.Error(), "invalid key") {
		t.Errorf("unexpected message: %v", err)
	}
}

func TestValidateEmptyValueDisallowed(t *testing.T) {
	s := makeValidSet(t)
	_ = s.Put("EMPTY_VAR", "")
	err := Validate(s, false)
	if err == nil {
		t.Fatal("expected error for empty value")
	}
	if !strings.Contains(err.Error(), "empty value") {
		t.Errorf("unexpected message: %v", err)
	}
}

func TestValidateEmptyValueAllowed(t *testing.T) {
	s := makeValidSet(t)
	_ = s.Put("EMPTY_VAR", "")
	if err := Validate(s, true); err != nil {
		t.Fatalf("unexpected error when allowEmpty=true: %v", err)
	}
}

func TestValidateMultipleViolations(t *testing.T) {
	s := makeValidSet(t)
	_ = s.Put("BAD KEY", "")
	err := Validate(s, false)
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(ve.Violations))
	}
}
