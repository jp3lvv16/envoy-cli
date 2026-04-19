package env

import (
	"errors"
	"sort"
)

// SortOrder defines the direction of sorting.
type SortOrder int

const (
	Ascending  SortOrder = iota
	Descending
)

// SortedKeys returns the keys of the set sorted by key name.
func SortedKeys(s *Set, order SortOrder) ([]string, error) {
	if s == nil {
		return nil, errors.New("sort: set must not be nil")
	}

	keys := make([]string, 0, len(s.Vars))
	for k := range s.Vars {
		keys = append(keys, k)
	}

	if order == Ascending {
		sort.Strings(keys)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	}

	return keys, nil
}

// SortedPairs returns key-value pairs sorted by key.
type Pair struct {
	Key   string
	Value string
}

func SortedPairs(s *Set, order SortOrder) ([]Pair, error) {
	keys, err := SortedKeys(s, order)
	if err != nil {
		return nil, err
	}

	pairs := make([]Pair, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, Pair{Key: k, Value: s.Vars[k]})
	}
	return pairs, nil
}
