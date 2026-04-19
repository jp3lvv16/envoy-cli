package env

import (
	"testing"
)

func makeSnapshotSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("snap-src")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("HOST", "localhost")
	_ = s.Put("PORT", "8080")
	return s
}

func TestTakeSnapshotNilReturnsError(t *testing.T) {
	_, err := TakeSnapshot(nil)
	if err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestTakeSnapshotCapturesVars(t *testing.T) {
	s := makeSnapshotSet(t)
	snap, err := TakeSnapshot(s)
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}
	if snap.SetName != "snap-src" {
		t.Errorf("expected SetName %q, got %q", "snap-src", snap.SetName)
	}
	if snap.Vars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", snap.Vars["HOST"])
	}
	if snap.TakenAt.IsZero() {
		t.Error("expected non-zero TakenAt")
	}
}

func TestTakeSnapshotIsIndependent(t *testing.T) {
	s := makeSnapshotSet(t)
	snap, _ := TakeSnapshot(s)
	_ = s.Put("HOST", "changed")
	if snap.Vars["HOST"] != "localhost" {
		t.Error("snapshot should not reflect mutations to original set")
	}
}

func TestRestoreSnapshotNilSnapReturnsError(t *testing.T) {
	dst, _ := NewSet("dst")
	if err := RestoreSnapshot(nil, dst); err == nil {
		t.Fatal("expected error for nil snapshot")
	}
}

func TestRestoreSnapshotNilDstReturnsError(t *testing.T) {
	snap := &Snapshot{Vars: map[string]string{}}
	if err := RestoreSnapshot(snap, nil); err == nil {
		t.Fatal("expected error for nil destination")
	}
}

func TestRestoreSnapshotAppliesVars(t *testing.T) {
	src := makeSnapshotSet(t)
	snap, _ := TakeSnapshot(src)

	dst, _ := NewSet("dst")
	if err := RestoreSnapshot(snap, dst); err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}
	v, err := dst.Get("PORT")
	if err != nil || v != "8080" {
		t.Errorf("expected PORT=8080, got %q (err=%v)", v, err)
	}
}
