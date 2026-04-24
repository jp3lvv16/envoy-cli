package env

import (
	"testing"
)

func makeSchemaSet(t *testing.T) *Set {
	t.Helper()
	s, err := NewSet("schema-test")
	if err != nil {
		t.Fatalf("NewSet: %v", err)
	}
	_ = s.Put("HOST", "localhost")
	_ = s.Put("PORT", "8080")
	_ = s.Put("ENV", "production")
	return s
}

func TestValidateSchemaНilSetReturnsError(t *testing.T) {
	schema := &Schema{Name: "test", Fields: []FieldSchema{}}
	if err := ValidateSchema(nil, schema); err == nil {
		t.Fatal("expected error for nil set")
	}
}

func TestValidateSchemaNilSchemaReturnsError(t *testing.T) {
	s := makeSchemaSet(t)
	if err := ValidateSchema(s, nil); err == nil {
		t.Fatal("expected error for nil schema")
	}
}

func TestValidateSchemaAllRequiredPresent(t *testing.T) {
	s := makeSchemaSet(t)
	schema := &Schema{
		Name: "test",
		Fields: []FieldSchema{
			{Key: "HOST", Required: true},
			{Key: "PORT", Required: true},
		},
	}
	if err := ValidateSchema(s, schema); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSchemaMissingRequiredKeyReturnsError(t *testing.T) {
	s := makeSchemaSet(t)
	schema := &Schema{
		Name: "test",
		Fields: []FieldSchema{
			{Key: "MISSING_KEY", Required: true},
		},
	}
	err := ValidateSchema(s, schema)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	se, ok := err.(*SchemaError)
	if !ok {
		t.Fatalf("expected *SchemaError, got %T", err)
	}
	if se.Key != "MISSING_KEY" {
		t.Errorf("expected key MISSING_KEY, got %s", se.Key)
	}
}

func TestValidateSchemaOptionalMissingKeyIsOk(t *testing.T) {
	s := makeSchemaSet(t)
	schema := &Schema{
		Name: "test",
		Fields: []FieldSchema{
			{Key: "OPTIONAL_KEY", Required: false},
		},
	}
	if err := ValidateSchema(s, schema); err != nil {
		t.Fatalf("unexpected error for optional missing key: %v", err)
	}
}

func TestValidateSchemaPatternMatch(t *testing.T) {
	s := makeSchemaSet(t)
	schema := &Schema{
		Name: "test",
		Fields: []FieldSchema{
			{Key: "PORT", Required: true, Pattern: `^\d+$`},
		},
	}
	if err := ValidateSchema(s, schema); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSchemaPatternMismatchReturnsError(t *testing.T) {
	s := makeSchemaSet(t)
	schema := &Schema{
		Name: "test",
		Fields: []FieldSchema{
			{Key: "HOST", Required: true, Pattern: `^\d+$`},
		},
	}
	err := ValidateSchema(s, schema)
	if err == nil {
		t.Fatal("expected pattern mismatch error")
	}
	se, ok := err.(*SchemaError)
	if !ok {
		t.Fatalf("expected *SchemaError, got %T", err)
	}
	if se.Key != "HOST" {
		t.Errorf("expected key HOST, got %s", se.Key)
	}
}
