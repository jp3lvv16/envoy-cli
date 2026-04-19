package env

import (
	"testing"
)

func TestAddEmptyTagReturnsError(t *testing.T) {
	idx := NewTagIndex()
	if err := idx.Add("", "mySet"); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestAddEmptySetReturnsError(t *testing.T) {
	idx := NewTagIndex()
	if err := idx.Add("prod", ""); err == nil {
		t.Fatal("expected error for empty set name")
	}
}

func TestAddAndRetrieveSets(t *testing.T) {
	idx := NewTagIndex()
	_ = idx.Add("prod", "api")
	_ = idx.Add("prod", "db")
	sets, err := idx.SetsForTag("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sets) != 2 {
		t.Fatalf("expected 2 sets, got %d", len(sets))
	}
}

func TestAddDuplicateSetIsIdempotent(t *testing.T) {
	idx := NewTagIndex()
	_ = idx.Add("prod", "api")
	_ = idx.Add("prod", "api")
	sets, _ := idx.SetsForTag("prod")
	if len(sets) != 1 {
		t.Fatalf("expected 1 set, got %d", len(sets))
	}
}

func TestRemoveSetDeletesEmptyTag(t *testing.T) {
	idx := NewTagIndex()
	_ = idx.Add("staging", "web")
	if err := idx.Remove("staging", "web"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := idx.SetsForTag("staging"); err == nil {
		t.Fatal("expected tag to be deleted")
	}
}

func TestRemoveNonExistentTagReturnsError(t *testing.T) {
	idx := NewTagIndex()
	if err := idx.Remove("ghost", "set"); err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestTagsReturnsSortedNames(t *testing.T) {
	idx := NewTagIndex()
	_ = idx.Add("zebra", "s1")
	_ = idx.Add("alpha", "s2")
	_ = idx.Add("mango", "s3")
	tags := idx.Tags()
	if tags[0] != "alpha" || tags[1] != "mango" || tags[2] != "zebra" {
		t.Fatalf("unexpected order: %v", tags)
	}
}

func TestSetsForMissingTagReturnsError(t *testing.T) {
	idx := NewTagIndex()
	if _, err := idx.SetsForTag("nope"); err == nil {
		t.Fatal("expected error")
	}
}
