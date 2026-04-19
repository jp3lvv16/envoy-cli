package env

import (
	"errors"
	"sort"
	"strings"
)

// Tag associates a label with a Set for grouping/filtering.
type Tag struct {
	Name string
	Sets []string
}

// TagIndex maps tag names to their Tag structs.
type TagIndex map[string]*Tag

// NewTagIndex creates an empty TagIndex.
func NewTagIndex() TagIndex {
	return make(TagIndex)
}

// Add adds a set name to the given tag, creating the tag if needed.
func (idx TagIndex) Add(tag, setName string) error {
	if strings.TrimSpace(tag) == "" {
		return errors.New("tag name must not be empty")
	}
	if strings.TrimSpace(setName) == "" {
		return errors.New("set name must not be empty")
	}
	t, ok := idx[tag]
	if !ok {
		t = &Tag{Name: tag}
		idx[tag] = t
	}
	for _, s := range t.Sets {
		if s == setName {
			return nil
		}
	}
	t.Sets = append(t.Sets, setName)
	sort.Strings(t.Sets)
	return nil
}

// Remove removes a set name from a tag. If the tag becomes empty it is deleted.
func (idx TagIndex) Remove(tag, setName string) error {
	t, ok := idx[tag]
	if !ok {
		return errors.New("tag not found: " + tag)
	}
	newSets := t.Sets[:0]
	for _, s := range t.Sets {
		if s != setName {
			newSets = append(newSets, s)
		}
	}
	if len(newSets) == 0 {
		delete(idx, tag)
	} else {
		t.Sets = newSets
	}
	return nil
}

// SetsForTag returns all set names associated with a tag.
func (idx TagIndex) SetsForTag(tag string) ([]string, error) {
	t, ok := idx[tag]
	if !ok {
		return nil, errors.New("tag not found: " + tag)
	}
	out := make([]string, len(t.Sets))
	copy(out, t.Sets)
	return out, nil
}

// Tags returns all tag names in sorted order.
func (idx TagIndex) Tags() []string {
	keys := make([]string, 0, len(idx))
	for k := range idx {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
