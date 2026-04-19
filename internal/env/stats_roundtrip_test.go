package env

import (
	"strings"
	"testing"
)

// TestStatAfterImport ensures Stat works correctly on an imported set.
func TestStatAfterImport(t *testing.T) {
	raw := "APP=myapp\nPORT=9000\nSECRET=\n"

	dst, err := NewSet("imported")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}

	if err := Import(dst, "dotenv", strings.NewReader(raw)); err != nil {
		t.Fatalf("Import: %v", err)
	}

	st, err := Stat(dst)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	if st.Count != 3 {
		t.Errorf("Count: got %d, want 3", st.Count)
	}
	if st.EmptyValues != 1 {
		t.Errorf("EmptyValues: got %d, want 1", st.EmptyValues)
	}
	if st.UniqueValues != 3 {
		t.Errorf("UniqueValues: got %d, want 3", st.UniqueValues)
	}
}
