package env

import (
	"testing"
)

func makeSearchSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("search-test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("DB_HOST", "localhost")
	_ = s.Put("DB_PORT", "5432")
	_ = s.Put("APP_SECRET", "s3cr3t")
	_ = s.Put("APP_DEBUG", "true")
	return s
}

func TestSearchByKeyNilSetReturnsError(t *testing.T) {
	_, err := SearchByKey(nil, "DB")
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestSearchByKeyEmptySubstrReturnsError(t *testing.T) {
	s := makeSearchSet(t)
	_, err := SearchByKey(s, "")
	if err == nil {
		t.Fatal("expected error for empty substr")
	}
}

func TestSearchByKeyFindsMatches(t *testing.T) {
	s := makeSearchSet(t)
	results, err := SearchByKey(s, "db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key != "DB_HOST" && r.Key != "DB_PORT" {
			t.Errorf("unexpected key: %s", r.Key)
		}
	}
}

func TestSearchByKeyNoMatchReturnsEmpty(t *testing.T) {
	s := makeSearchSet(t)
	results, err := SearchByKey(s, "REDIS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchByValueNilSetReturnsError(t *testing.T) {
	_, err := SearchByValue(nil, "local")
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestSearchByValueFindsMatches(t *testing.T) {
	s := makeSearchSet(t)
	results, err := SearchByValue(s, "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "APP_DEBUG" {
		t.Fatalf("expected APP_DEBUG, got %+v", results)
	}
}

func TestSearchByValueCaseInsensitive(t *testing.T) {
	s := makeSearchSet(t)
	results, err := SearchByValue(s, "LOCALHOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "DB_HOST" {
		t.Fatalf("expected DB_HOST, got %+v", results)
	}
}
