package env

import (
	"errors"
	"strings"
)

// TransformFunc is applied to each value in a Set.
type TransformFunc func(key, value string) string

// Transform applies fn to every value in src, writing results into a new Set.
// The original Set is not modified.
func Transform(src *Set, fn TransformFunc) (*Set, error) {
	if src == nil {
		return nil, errors.New("transform: src set is nil")
	}
	if fn == nil {
		return nil, errors.New("transform: transform func is nil")
	}

	dst, err := NewSet(src.Name())
	if err != nil {
		return nil, err
	}

	for _, k := range src.Keys() {
		v, _ := src.Get(k)
		transformed := fn(k, v)
		if err := dst.Put(k, transformed); err != nil {
			return nil, err
		}
	}
	return dst, nil
}

// UppercaseValues returns a TransformFunc that converts all values to uppercase.
func UppercaseValues() TransformFunc {
	return func(_, value string) string {
		return strings.ToUpper(value)
	}
}

// PrefixValues returns a TransformFunc that prepends prefix to every value.
func PrefixValues(prefix string) TransformFunc {
	return func(_, value string) string {
		return prefix + value
	}
}
