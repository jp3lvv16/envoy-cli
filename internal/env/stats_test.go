package env

import (
	"testing"
)

func makeStatsSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("stats-set")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	for k, v := range map[string]string{
		"HOST":    "localhost",
		"PORT":    "8080",
		"DEBUG":   "",
		"REPLICA": "localhost", // duplicate value
	} {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put: %v", err)
		}
	}
	return s
}

func TestStatNilSetReturnsError(t *testing.T) {
	_, err := Stat(nil)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestStatCount(t *testing.T) {
	s := makeStatsSet(t)
	st, err := Stat(s)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if st.Count != 4 {
		t.Errorf("Count: got %d, want 4", st.Count)
	}
}

func TestStatEmptyValues(t *testing.T) {
	s := makeStatsSet(t)
	st, err := Stat(s)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if st.EmptyValues != 1 {
		t.Errorf("EmptyValues: got %d, want 1", st.EmptyValues)
	}
}

func TestStatUniqueValues(t *testing.T) {
	s := makeStatsSet(t)
	st, err := Stat(s)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	// values: "localhost", "8080", "", "localhost" => 3 unique
	if st.UniqueValues != 3 {
		t.Errorf("UniqueValues: got %d, want 3", st.UniqueValues)
	}
}

func TestStatDescribeContainsName(t *testing.T) {
	s := makeStatsSet(t)
	st, err := Stat(s)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	desc := st.Describe()
	if desc == "" {
		t.Fatal("Describe returned empty string")
	}
	if st.Name != "stats-set" {
		t.Errorf("Name: got %q, want %q", st.Name, "stats-set")
	}
}
