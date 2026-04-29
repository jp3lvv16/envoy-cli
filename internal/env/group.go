package env

import (
	"errors"
	"sort"
)

// GroupIndex maps group names to sets of env-set names.
type GroupIndex struct {
	groups map[string]map[string]struct{}
}

// NewGroupIndex creates an empty GroupIndex.
func NewGroupIndex() *GroupIndex {
	return &GroupIndex{groups: make(map[string]map[string]struct{})}
}

// Add associates setName with groupName. Both must be non-empty.
func (g *GroupIndex) Add(groupName, setName string) error {
	if groupName == "" {
		return errors.New("group name must not be empty")
	}
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	if _, ok := g.groups[groupName]; !ok {
		g.groups[groupName] = make(map[string]struct{})
	}
	g.groups[groupName][setName] = struct{}{}
	return nil
}

// Remove removes setName from groupName. Deletes the group if empty.
func (g *GroupIndex) Remove(groupName, setName string) error {
	if groupName == "" {
		return errors.New("group name must not be empty")
	}
	if setName == "" {
		return errors.New("set name must not be empty")
	}
	members, ok := g.groups[groupName]
	if !ok {
		return errors.New("group not found: " + groupName)
	}
	if _, exists := members[setName]; !exists {
		return errors.New("set not in group: " + setName)
	}
	delete(members, setName)
	if len(members) == 0 {
		delete(g.groups, groupName)
	}
	return nil
}

// Members returns a sorted list of set names belonging to groupName.
func (g *GroupIndex) Members(groupName string) ([]string, error) {
	if groupName == "" {
		return nil, errors.New("group name must not be empty")
	}
	members, ok := g.groups[groupName]
	if !ok {
		return []string{}, nil
	}
	out := make([]string, 0, len(members))
	for s := range members {
		out = append(out, s)
	}
	sort.Strings(out)
	return out, nil
}

// Groups returns a sorted list of all group names.
func (g *GroupIndex) Groups() []string {
	out := make([]string, 0, len(g.groups))
	for name := range g.groups {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
