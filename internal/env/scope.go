package env

import "fmt"

// Scope represents a named environment scope (e.g. dev, staging, prod)
// that wraps a Set and carries a priority level for resolution.
type Scope struct {
	Name     string
	Priority int
	Set      *Set
}

// ScopeResolver resolves a key across multiple scopes in priority order.
type ScopeResolver struct {
	scopes []*Scope
}

// NewScopeResolver creates a ScopeResolver from the given scopes.
// Scopes are sorted by descending priority at resolution time.
func NewScopeResolver(scopes []*Scope) (*ScopeResolver, error) {
	if len(scopes) == 0 {
		return nil, fmt.Errorf("scope: at least one scope is required")
	}
	for _, s := range scopes {
		if s == nil || s.Set == nil {
			return nil, fmt.Errorf("scope: nil scope or set provided")
		}
		if s.Name == "" {
			return nil, fmt.Errorf("scope: scope name must not be empty")
		}
	}
	return &ScopeResolver{scopes: scopes}, nil
}

// Resolve returns the value for key from the highest-priority scope that
// contains it, along with the name of the winning scope.
func (r *ScopeResolver) Resolve(key string) (string, string, error) {
	if key == "" {
		return "", "", fmt.Errorf("scope: key must not be empty")
	}
	best := (*Scope)(nil)
	for _, s := range r.scopes {
		val, err := s.Set.Get(key)
		if err != nil {
			continue
		}
		if val == "" {
			continue
		}
		if best == nil || s.Priority > best.Priority {
			best = s
		}
	}
	if best == nil {
		return "", "", fmt.Errorf("scope: key %q not found in any scope", key)
	}
	val, _ := best.Set.Get(key)
	return val, best.Name, nil
}

// ResolveAll returns a flat Set containing the winning value for every key
// found across all scopes, using the highest-priority scope per key.
func (r *ScopeResolver) ResolveAll(name string) (*Set, error) {
	if name == "" {
		return nil, fmt.Errorf("scope: result set name must not be empty")
	}
	out, err := NewSet(name)
	if err != nil {
		return nil, err
	}
	seen := map[string]int{} // key -> best priority so far
	for _, s := range r.scopes {
		keys, err := SortedKeys(s.Set, true)
		if err != nil {
			continue
		}
		for _, k := range keys {
			v, _ := s.Set.Get(k)
			prev, exists := seen[k]
			if !exists || s.Priority > prev {
				_ = out.Put(k, v)
				seen[k] = s.Priority
			}
		}
	}
	return out, nil
}
