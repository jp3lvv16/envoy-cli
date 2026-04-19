package env

import (
	"errors"
	"strings"
)

// SearchResult holds a single match found during a search.
type SearchResult struct {
	Key   string
	Value string
}

// SearchByKey returns all entries whose key contains the given substring (case-insensitive).
func SearchByKey(s *Set, substr string) ([]SearchResult, error) {
	if s == nil {
		return nil, errors.New("search: set must not be nil")
	}
	if substr == "" {
		return nil, errors.New("search: substr must not be empty")
	}
	lower := strings.ToLower(substr)
	var results []SearchResult
	for _, k := range s.Keys() {
		if strings.Contains(strings.ToLower(k), lower) {
			v, _ := s.Get(k)
			results = append(results, SearchResult{Key: k, Value: v})
		}
	}
	return results, nil
}

// SearchByValue returns all entries whose value contains the given substring (case-insensitive).
func SearchByValue(s *Set, substr string) ([]SearchResult, error) {
	if s == nil {
		return nil, errors.New("search: set must not be nil")
	}
	if substr == "" {
		return nil, errors.New("search: substr must not be empty")
	}
	lower := strings.ToLower(substr)
	var results []SearchResult
	for _, k := range s.Keys() {
		v, _ := s.Get(k)
		if strings.Contains(strings.ToLower(v), lower) {
			results = append(results, SearchResult{Key: k, Value: v})
		}
	}
	return results, nil
}
