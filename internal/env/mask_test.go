package env

import (
	"testing"
)

func makeMaskSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	pairs := map[string]string{
		"APP_HOST":     "localhost",
		"APP_PORT":     "8080",
		"DB_PASSWORD":  "supersecret",
		"API_KEY":      "abc123",
		"AUTH_TOKEN":   "tok_xyz",
		"DESCRIPTION":  "hello world",
	}
	for k, v := range pairs {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put(%s): %v", k, err)
		}
	}
	return s
}

func TestMaskNilSetReturnsError(t *testing.T) {
	_, err := Mask(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestMaskSensitiveKeysAreRedacted(t *testing.T) {
	s := makeMaskSet(t)
	res, err := Mask(s, nil)
	if err != nil {
		t.Fatalf("Mask: %v", err)
	}

	redacted := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN"}
	for _, k := range redacted {
		if res.Vars[k] != maskedValue {
			t.Errorf("expected %s to be masked, got %q", k, res.Vars[k])
		}
	}
}

func TestMaskNonSensitiveKeysArePreserved(t *testing.T) {
	s := makeMaskSet(t)
	res, err := Mask(s, nil)
	if err != nil {
		t.Fatalf("Mask: %v", err)
	}

	if res.Vars["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", res.Vars["APP_HOST"])
	}
	if res.Vars["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080, got %q", res.Vars["APP_PORT"])
	}
	if res.Vars["DESCRIPTION"] != "hello world" {
		t.Errorf("expected DESCRIPTION='hello world', got %q", res.Vars["DESCRIPTION"])
	}
}

func TestMaskOriginalSetUnchanged(t *testing.T) {
	s := makeMaskSet(t)
	_, err := Mask(s, nil)
	if err != nil {
		t.Fatalf("Mask: %v", err)
	}
	v, _ := s.Get("DB_PASSWORD")
	if v != "supersecret" {
		t.Errorf("original set was mutated: DB_PASSWORD=%q", v)
	}
}

func TestMaskCustomSubstrings(t *testing.T) {
	s := makeMaskSet(t)
	res, err := Mask(s, []string{"host"})
	if err != nil {
		t.Fatalf("Mask: %v", err)
	}
	if res.Vars["APP_HOST"] != maskedValue {
		t.Errorf("expected APP_HOST to be masked with custom substrings")
	}
	// DB_PASSWORD should NOT be masked since custom list only has "host"
	if res.Vars["DB_PASSWORD"] == maskedValue {
		t.Errorf("DB_PASSWORD should not be masked with custom substrings [host]")
	}
}
