package env

import (
	"testing"
)

func makeAuditLog(t *testing.T) *AuditLog {
	t.Helper()
	l, err := NewAuditLog("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return l
}

func TestNewAuditLogEmptyNameReturnsError(t *testing.T) {
	_, err := NewAuditLog("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestNewAuditLogCreatesEmptyEntries(t *testing.T) {
	l := makeAuditLog(t)
	if len(l.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(l.Entries))
	}
}

func TestRecordEmptyActorReturnsError(t *testing.T) {
	l := makeAuditLog(t)
	err := l.Record("", AuditPut, "KEY", "", "val")
	if err == nil {
		t.Fatal("expected error for empty actor")
	}
}

func TestRecordEmptyActionReturnsError(t *testing.T) {
	l := makeAuditLog(t)
	err := l.Record("alice", "", "KEY", "", "val")
	if err == nil {
		t.Fatal("expected error for empty action")
	}
}

func TestRecordAppendsEntry(t *testing.T) {
	l := makeAuditLog(t)
	if err := l.Record("alice", AuditPut, "DB_HOST", "", "localhost"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries))
	}
	e := l.Entries[0]
	if e.Actor != "alice" || e.Action != AuditPut || e.Key != "DB_HOST" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestFilterByActorReturnsMatchingEntries(t *testing.T) {
	l := makeAuditLog(t)
	_ = l.Record("alice", AuditPut, "A", "", "1")
	_ = l.Record("bob", AuditDelete, "B", "2", "")
	_ = l.Record("alice", AuditImport, "C", "", "3")

	result := l.FilterByActor("alice")
	if len(result) != 2 {
		t.Fatalf("expected 2 entries for alice, got %d", len(result))
	}
}

func TestFilterByActionReturnsMatchingEntries(t *testing.T) {
	l := makeAuditLog(t)
	_ = l.Record("alice", AuditPut, "A", "", "1")
	_ = l.Record("bob", AuditDelete, "B", "2", "")
	_ = l.Record("carol", AuditPut, "C", "", "3")

	result := l.FilterByAction(AuditPut)
	if len(result) != 2 {
		t.Fatalf("expected 2 put entries, got %d", len(result))
	}
}

func TestSummaryContainsExpectedFields(t *testing.T) {
	l := makeAuditLog(t)
	_ = l.Record("alice", AuditPut, "DB_URL", "", "postgres://")
	s := Summary(l.Entries[0])
	for _, want := range []string{"put", "alice", "DB_URL"} {
		if !contains(s, want) {
			t.Errorf("summary %q missing %q", s, want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
