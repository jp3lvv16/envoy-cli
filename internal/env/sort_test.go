package env

import (
	"testing"
)

func makeSortSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("sort-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = s.Put("ZEBRA", "z")
	_ = s.Put("ALPHA", "a")
	_ = s.Put("MANGO", "m")
	return s
}

func TestSortedKeysNilSetReturnsError(t *testing.T) {
	_, err := SortedKeys(nil, Ascending)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestSortedKeysAscending(t *testing.T) {
	s := makeSortSet(t)
	keys, err := SortedKeys(s, Ascending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ALPHA", "MANGO", "ZEBRA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSortedKeysDescending(t *testing.T) {
	s := makeSortSet(t)
	keys, err := SortedKeys(s, Descending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ZEBRA", "MANGO", "ALPHA"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: got %q, want %q", i, k, expected[i])
		}
	}
}

func TestSortedPairsReturnsCorrectValues(t *testing.T) {
	s := makeSortSet(t)
	pairs, err := SortedPairs(s, Ascending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	if pairs[0].Key != "ALPHA" || pairs[0].Value != "a" {
		t.Errorf("unexpected first pair: %+v", pairs[0])
	}
}

func TestSortedPairsNilSetReturnsError(t *testing.T) {
	_, err := SortedPairs(nil, Ascending)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}
