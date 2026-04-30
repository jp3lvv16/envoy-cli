package env

import (
	"testing"
)

func makeTrimSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("trim-test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	return s
}

func TestTrimNilSetReturnsError(t *testing.T) {
	_, err := Trim(nil)
	if err == nil {
		t.Fatal("expected error for nil set, got nil")
	}
}

func TestTrimCleanSetNoChanges(t *testing.T) {
	s := makeTrimSet(t)
	_ = s.Put("KEY", "value")
	_ = s.Put("OTHER", "clean")

	results, err := Trim(s)
	if err != nil {
		t.Fatalf("Trim: %v", err)
	}
	for _, r := range results {
		if r.Changed {
			t.Errorf("expected no change for key %q, but got changed=true", r.Key)
		}
	}
}

func TestTrimRemovesLeadingAndTrailingWhitespace(t *testing.T) {
	s := makeTrimSet(t)
	_ = s.Put("KEY", "  hello world  ")

	_, err := Trim(s)
	if err != nil {
		t.Fatalf("Trim: %v", err)
	}

	val, _ := s.Get("KEY")
	if val != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", val)
	}
}

func TestTrimResultReportsChangedFlag(t *testing.T) {
	s := makeTrimSet(t)
	_ = s.Put("DIRTY", "\t spaced \n")
	_ = s.Put("CLEAN", "fine")

	results, err := Trim(s)
	if err != nil {
		t.Fatalf("Trim: %v", err)
	}

	changedCount := 0
	for _, r := range results {
		if r.Key == "DIRTY" && !r.Changed {
			t.Error("expected DIRTY to be marked as changed")
		}
		if r.Key == "CLEAN" && r.Changed {
			t.Error("expected CLEAN to not be marked as changed")
		}
		if r.Changed {
			changedCount++
		}
	}
	if changedCount != 1 {
		t.Errorf("expected 1 changed key, got %d", changedCount)
	}
}

func TestTrimKeysNilSetReturnsError(t *testing.T) {
	_, err := TrimKeys(nil)
	if err == nil {
		t.Fatal("expected error for nil set, got nil")
	}
}

func TestTrimKeysReturnsOnlyChangedEntries(t *testing.T) {
	s := makeTrimSet(t)
	_ = s.Put("PADDED", " value ")
	_ = s.Put("NORMAL", "ok")

	changed, err := TrimKeys(s)
	if err != nil {
		t.Fatalf("TrimKeys: %v", err)
	}
	if len(changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(changed))
	}
	if changed[0].Key != "PADDED" {
		t.Errorf("expected PADDED, got %q", changed[0].Key)
	}
	if changed[0].OldValue != " value " {
		t.Errorf("unexpected OldValue: %q", changed[0].OldValue)
	}
	if changed[0].NewValue != "value" {
		t.Errorf("unexpected NewValue: %q", changed[0].NewValue)
	}
}
