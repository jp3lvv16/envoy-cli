package env

import "testing"

func TestSnapshotRoundTrip(t *testing.T) {
	src, err := NewSet("rt-src")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	pairs := map[string]string{
		"DB_HOST": "db.local",
		"DB_PORT": "5432",
		"APP_ENV": "staging",
	}
	for k, v := range pairs {
		if err := src.Put(k, v); err != nil {
			t.Fatalf("Put(%q): %v", k, err)
		}
	}

	snap, err := TakeSnapshot(src)
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	dst, err := NewSet("rt-dst")
	if err != nil {
		t.Fatalf("NewSet dst: %v", err)
	}
	if err := RestoreSnapshot(snap, dst); err != nil {
		t.Fatalf("RestoreSnapshot: %v", err)
	}

	for k, want := range pairs {
		got, err := dst.Get(k)
		if err != nil {
			t.Errorf("Get(%q): %v", k, err)
			continue
		}
		if got != want {
			t.Errorf("key %q: want %q, got %q", k, want, got)
		}
	}
}
