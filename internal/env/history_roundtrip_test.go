package env

import (
	"testing"
)

// TestHistoryRoundTripMultipleRollbacks records several states and rolls back
// through each one, verifying correctness at every step.
func TestHistoryRoundTripMultipleRollbacks(t *testing.T) {
	s, _ := NewSet("staging")
	h, _ := NewHistory("staging")

	states := []struct {
		key string
		val string
	}{
		{"DB_HOST", "localhost"},
		{"DB_HOST", "db.staging.internal"},
		{"DB_HOST", "db.prod.internal"},
	}

	for _, st := range states {
		_ = s.Put(st.key, st.val)
		if err := h.Record(s, st.val); err != nil {
			t.Fatalf("Record failed: %v", err)
		}
	}

	for i := len(states) - 1; i >= 0; i-- {
		if err := h.Rollback(s, i); err != nil {
			t.Fatalf("Rollback(%d) failed: %v", i, err)
		}
		v, err := s.Get("DB_HOST")
		if err != nil {
			t.Fatalf("Get after rollback %d: %v", i, err)
		}
		if v != states[i].val {
			t.Errorf("rollback %d: expected %q, got %q", i, states[i].val, v)
		}
	}
}
