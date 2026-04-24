package env

import (
	"testing"
)

func makePatchSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("patch-test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("HOST", "localhost")
	_ = s.Put("PORT", "8080")
	_ = s.Put("DEBUG", "true")
	return s
}

func TestPatchNilSetReturnsError(t *testing.T) {
	err := Patch(nil, []PatchOp{{Op: "set", Key: "K", Value: "v"}})
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestPatchSetOp(t *testing.T) {
	s := makePatchSet(t)
	err := Patch(s, []PatchOp{{Op: "set", Key: "HOST", Value: "prod.example.com"}})
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}
	val, _ := s.Get("HOST")
	if val != "prod.example.com" {
		t.Errorf("expected prod.example.com, got %q", val)
	}
}

func TestPatchDeleteOp(t *testing.T) {
	s := makePatchSet(t)
	err := Patch(s, []PatchOp{{Op: "delete", Key: "DEBUG"}})
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}
	_, err = s.Get("DEBUG")
	if err == nil {
		t.Fatal("expected error getting deleted key")
	}
}

func TestPatchRenameOp(t *testing.T) {
	s := makePatchSet(t)
	err := Patch(s, []PatchOp{{Op: "rename", Key: "PORT", NewKey: "APP_PORT"}})
	if err != nil {
		t.Fatalf("Patch: %v", err)
	}
	val, err := s.Get("APP_PORT")
	if err != nil {
		t.Fatalf("Get APP_PORT: %v", err)
	}
	if val != "8080" {
		t.Errorf("expected 8080, got %q", val)
	}
	_, err = s.Get("PORT")
	if err == nil {
		t.Fatal("expected old key PORT to be gone")
	}
}

func TestPatchUnknownOpReturnsError(t *testing.T) {
	s := makePatchSet(t)
	err := Patch(s, []PatchOp{{Op: "upsert", Key: "X", Value: "y"}})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestPatchMultipleOpsAppliedInOrder(t *testing.T) {
	s := makePatchSet(t)
	ops := []PatchOp{
		{Op: "set", Key: "ENV", Value: "staging"},
		{Op: "delete", Key: "DEBUG"},
		{Op: "rename", Key: "HOST", NewKey: "APP_HOST"},
	}
	if err := Patch(s, ops); err != nil {
		t.Fatalf("Patch: %v", err)
	}
	if v, _ := s.Get("ENV"); v != "staging" {
		t.Errorf("ENV: expected staging, got %q", v)
	}
	if _, err := s.Get("DEBUG"); err == nil {
		t.Error("DEBUG should have been deleted")
	}
	if v, _ := s.Get("APP_HOST"); v != "localhost" {
		t.Errorf("APP_HOST: expected localhost, got %q", v)
	}
}
