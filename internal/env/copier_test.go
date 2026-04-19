package env

import (
	"testing"
)

func TestCopyNilSrcReturnsError(t *testing.T) {
	dst, _ := NewSet("dst")
	_, err := Copy(nil, dst, CopyOptions{})
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestCopyNilDstReturnsError(t *testing.T) {
	src, _ := NewSet("src")
	_, err := Copy(src, nil, CopyOptions{})
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestCopyNoOverwrite(t *testing.T) {
	src, _ := NewSet("src")
	_ = src.Put("KEY1", "val1")
	_ = src.Put("KEY2", "val2")

	dst, _ := NewSet("dst")
	_ = dst.Put("KEY1", "original")

	n, err := Copy(src, dst, CopyOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 copied, got %d", n)
	}

	v, _ := dst.Get("KEY1")
	if v != "original" {
		t.Errorf("expected KEY1 to remain 'original', got %q", v)
	}
	v2, _ := dst.Get("KEY2")
	if v2 != "val2" {
		t.Errorf("expected KEY2 to be 'val2', got %q", v2)
	}
}

func TestCopyWithOverwrite(t *testing.T) {
	src, _ := NewSet("src")
	_ = src.Put("KEY1", "new")

	dst, _ := NewSet("dst")
	_ = dst.Put("KEY1", "old")

	n, err := Copy(src, dst, CopyOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1 copied, got %d", n)
	}
	v, _ := dst.Get("KEY1")
	if v != "new" {
		t.Errorf("expected KEY1='new', got %q", v)
	}
}

func TestMergeOverwritesExisting(t *testing.T) {
	src, _ := NewSet("src")
	_ = src.Put("A", "1")
	_ = src.Put("B", "2")

	dst, _ := NewSet("dst")
	_ = dst.Put("A", "old")

	n, err := Merge(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 merged, got %d", n)
	}
	v, _ := dst.Get("A")
	if v != "1" {
		t.Errorf("expected A='1', got %q", v)
	}
}
