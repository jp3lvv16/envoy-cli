package env

import (
	"fmt"
	"regexp"
)

// FieldSchema describes the expected shape of a single environment variable.
type FieldSchema struct {
	Key      string
	Required bool
	Pattern  string // optional regex pattern the value must match
}

// Schema holds a collection of field definitions for a named set.
type Schema struct {
	Name   string
	Fields []FieldSchema
}

// SchemaError is returned when validation against a schema fails.
type SchemaError struct {
	Key     string
	Reason  string
}

func (e *SchemaError) Error() string {
	return fmt.Sprintf("schema violation for key %q: %s", e.Key, e.Reason)
}

// ValidateSchema checks that the given Set satisfies the Schema.
// It returns the first violation encountered, or nil on success.
func ValidateSchema(s *Set, schema *Schema) error {
	if s == nil {
		return fmt.Errorf("ValidateSchema: set must not be nil")
	}
	if schema == nil {
		return fmt.Errorf("ValidateSchema: schema must not be nil")
	}

	for _, field := range schema.Fields {
		val, err := s.Get(field.Key)
		if err != nil {
			// key is missing
			if field.Required {
				return &SchemaError{Key: field.Key, Reason: "required key is missing"}
			}
			continue
		}

		if field.Pattern != "" {
			re, compErr := regexp.Compile(field.Pattern)
			if compErr != nil {
				return fmt.Errorf("ValidateSchema: invalid pattern for key %q: %w", field.Key, compErr)
			}
			if !re.MatchString(val) {
				return &SchemaError{
					Key:    field.Key,
					Reason: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				}
			}
		}
	}

	return nil
}
