package env

import (
	"testing"
)

func makeMergeSet(t *testing.T, name string, pairs map[string]string) *Set {
	t.Helper()
	s, err := NewSet(name)
	if err != nil {
		t.Fatalf("NewSet(%q): %v", name, err)
	}
	for k, v := range pairs {
		if err := s.Put(k, v); err != nil {
			t.Fatalf("Put(%q, %q): %v", k, v, err)
		}
	}
	return s
}

func TestMergeWithStrategyNilSrcReturnsError(t *testing.T) {
	dst := makeMergeSet(t, "dst", nil)
	if err := MergeWithStrategy(nil, dst, MergeStrategySkip); err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestMergeWithStrategyNilDstReturnsError(t *testing.T) {
	src := makeMergeSet(t, "src", nil)
	if err := MergeWithStrategy(src, nil, MergeStrategySkip); err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestMergeStrategySkipKeepsDstOnConflict(t *testing.T) {
	src := makeMergeSet(t, "src", map[string]string{"KEY": "src-val", "ONLY_SRC": "yes"})
	dst := makeMergeSet(t, "dst", map[string]string{"KEY": "dst-val"})

	if err := MergeWithStrategy(src, dst, MergeStrategySkip); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := dst.Get("KEY"); v != "dst-val" {
		t.Errorf("expected dst-val, got %q", v)
	}
	if v, _ := dst.Get("ONLY_SRC"); v != "yes" {
		t.Errorf("expected ONLY_SRC to be merged, got %q", v)
	}
}

func TestMergeStrategyOverwriteUpdatesDst(t *testing.T) {
	src := makeMergeSet(t, "src", map[string]string{"KEY": "new-val"})
	dst := makeMergeSet(t, "dst", map[string]string{"KEY": "old-val"})

	if err := MergeWithStrategy(src, dst, MergeStrategyOverwrite); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, _ := dst.Get("KEY"); v != "new-val" {
		t.Errorf("expected new-val, got %q", v)
	}
}

func TestMergeStrategyErrorOnConflict(t *testing.T) {
	src := makeMergeSet(t, "src", map[string]string{"KEY": "val"})
	dst := makeMergeSet(t, "dst", map[string]string{"KEY": "other"})

	if err := MergeWithStrategy(src, dst, MergeStrategyError); err == nil {
		t.Fatal("expected error on conflict")
	}
}

func TestMergeStrategyNoConflictAlwaysSucceeds(t *testing.T) {
	src := makeMergeSet(t, "src", map[string]string{"A": "1", "B": "2"})
	dst := makeMergeSet(t, "dst", map[string]string{"C": "3"})

	if err := MergeWithStrategy(src, dst, MergeStrategyError); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, k := range []string{"A", "B", "C"} {
		if _, err := dst.Get(k); err != nil {
			t.Errorf("expected key %q in dst", k)
		}
	}
}
