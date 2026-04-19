package env

import "fmt"

// Stats holds summary statistics for a Set.
type Stats struct {
	Name     string
	Count    int
	EmptyValues int
	UniqueValues int
}

// Describe returns a human-readable summary string.
func (s Stats) Describe() string {
	return fmt.Sprintf("set=%s keys=%d empty=%d unique_values=%d",
		s.Name, s.Count, s.EmptyValues, s.UniqueValues)
}

// Stat computes statistics for the given Set.
func Stat(s *Set) (Stats, error) {
	if s == nil {
		return Stats{}, fmt.Errorf("stat: set must not be nil")
	}

	pairs, err := SortedPairs(s, true)
	if err != nil {
		return Stats{}, fmt.Errorf("stat: %w", err)
	}

	seen := make(map[string]struct{})
	empty := 0

	for _, p := range pairs {
		if p.Value == "" {
			empty++
		}
		seen[p.Value] = struct{}{}
	}

	return Stats{
		Name:         s.Name(),
		Count:        len(pairs),
		EmptyValues:  empty,
		UniqueValues: len(seen),
	}, nil
}
