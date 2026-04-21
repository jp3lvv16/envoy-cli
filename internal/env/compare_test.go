package env

import (
	"testing"
)

func makeCompareSet(t *testing.T, name string, pairs map[string]string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	for k, v := range pairs {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put(%q): %v", k, err)
		}
	}
	return s
}

func TestCompareNilSrcReturnsError(t *testing.T) {
	dst := makeCompareSet(t, "dst", nil)
	_, err := Compare(nil, dst)
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestCompareNilDstReturnsError(t *testing.T) {
	src := makeCompareSet(t, "src", nil)
	_, err := Compare(src, nil)
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestCompareIdenticalSets(t *testing.T) {
	pairs := map[string]string{"A": "1", "B": "2"}
	src := makeCompareSet(t, "src", pairs)
	dst := makeCompareSet(t, "dst", pairs)

	cr, err := Compare(src, dst)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(cr.Same) != 2 {
		t.Errorf("expected 2 same, got %d", len(cr.Same))
	}
	if len(cr.OnlyInSrc)+len(cr.OnlyInDst)+len(cr.Conflicted) != 0 {
		t.Error("expected no differences")
	}
}

func TestCompareDetectsOnlyInSrc(t *testing.T) {
	src := makeCompareSet(t, "src", map[string]string{"X": "10", "Y": "20"})
	dst := makeCompareSet(t, "dst", map[string]string{"X": "10"})

	cr, err := Compare(src, dst)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if _, ok := cr.OnlyInSrc["Y"]; !ok {
		t.Error("expected Y to be only in src")
	}
}

func TestCompareDetectsOnlyInDst(t *testing.T) {
	src := makeCompareSet(t, "src", map[string]string{"X": "10"})
	dst := makeCompareSet(t, "dst", map[string]string{"X": "10", "Z": "30"})

	cr, err := Compare(src, dst)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if _, ok := cr.OnlyInDst["Z"]; !ok {
		t.Error("expected Z to be only in dst")
	}
}

func TestCompareDetectsConflicts(t *testing.T) {
	src := makeCompareSet(t, "src", map[string]string{"K": "old"})
	dst := makeCompareSet(t, "dst", map[string]string{"K": "new"})

	cr, err := Compare(src, dst)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	pair, ok := cr.Conflicted["K"]
	if !ok {
		t.Fatal("expected K to be conflicted")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected conflict values: %v", pair)
	}
}

func TestIsEqualTrueForIdenticalSets(t *testing.T) {
	pairs := map[string]string{"A": "1"}
	src := makeCompareSet(t, "src", pairs)
	dst := makeCompareSet(t, "dst", pairs)

	eq, err := IsEqual(src, dst)
	if err != nil {
		t.Fatalf("IsEqual: %v", err)
	}
	if !eq {
		t.Error("expected sets to be equal")
	}
}

func TestIsEqualFalseForDifferentSets(t *testing.T) {
	src := makeCompareSet(t, "src", map[string]string{"A": "1"})
	dst := makeCompareSet(t, "dst", map[string]string{"A": "2"})

	eq, err := IsEqual(src, dst)
	if err != nil {
		t.Fatalf("IsEqual: %v", err)
	}
	if eq {
		t.Error("expected sets to be unequal")
	}
}
