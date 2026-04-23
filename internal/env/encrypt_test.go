package env

import (
	"strings"
	"testing"
)

func makeEncryptSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("enc-test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("DB_PASS", "secret123")
	_ = s.Put("API_KEY", "myapikey")
	return s
}

var testKey = []byte("thisis32byteslongkeyforAES256!!!")

func TestEncryptNilSetReturnsError(t *testing.T) {
	if err := Encrypt(nil, testKey); err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestDecryptNilSetReturnsError(t *testing.T) {
	if err := Decrypt(nil, testKey); err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestEncryptBadKeyLengthReturnsError(t *testing.T) {
	s := makeEncryptSet(t)
	if err := Encrypt(s, []byte("shortkey")); err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestDecryptBadKeyLengthReturnsError(t *testing.T) {
	s := makeEncryptSet(t)
	if err := Decrypt(s, []byte("shortkey")); err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestEncryptValuesArePrefixed(t *testing.T) {
	s := makeEncryptSet(t)
	if err := Encrypt(s, testKey); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	for k, v := range s.Vars {
		if !strings.HasPrefix(v, "enc:") {
			t.Errorf("key %q value not encrypted: %q", k, v)
		}
	}
}

func TestEncryptThenDecryptRestoresValues(t *testing.T) {
	s := makeEncryptSet(t)
	original := map[string]string{}
	for k, v := range s.Vars {
		original[k] = v
	}
	if err := Encrypt(s, testKey); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if err := Decrypt(s, testKey); err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	for k, want := range original {
		got, err := s.Get(k)
		if err != nil {
			t.Errorf("Get(%q): %v", k, err)
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestEncryptAlreadyEncryptedIsSkipped(t *testing.T) {
	s := makeEncryptSet(t)
	if err := Encrypt(s, testKey); err != nil {
		t.Fatalf("first Encrypt: %v", err)
	}
	snapshot := map[string]string{}
	for k, v := range s.Vars {
		snapshot[k] = v
	}
	if err := Encrypt(s, testKey); err != nil {
		t.Fatalf("second Encrypt: %v", err)
	}
	for k, want := range snapshot {
		got := s.Vars[k]
		if got != want {
			t.Errorf("key %q changed on re-encrypt: got %q, want %q", k, got, want)
		}
	}
}
