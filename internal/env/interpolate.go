package env

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR_NAME} and $VAR_NAME style references.
var interpolatePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateResult holds the output of an interpolation pass.
type InterpolateResult struct {
	// Resolved contains keys whose values were fully interpolated.
	Resolved []string
	// Unresolved contains keys that still reference missing variables after interpolation.
	Unresolved []string
}

// Interpolate expands variable references within the values of src using the
// variables defined in src itself (self-referential expansion). References
// take the form ${VAR} or $VAR. Missing references are left in place.
//
// A maximum of 10 expansion passes are performed to handle chained references.
// If src is nil an error is returned.
func Interpolate(src *Set) (*InterpolateResult, error) {
	if src == nil {
		return nil, fmt.Errorf("interpolate: source set must not be nil")
	}

	const maxPasses = 10

	for pass := 0; pass < maxPasses; pass++ {
		changed := false
		src.mu.Lock()
		for k, v := range src.vars {
			expanded := expand(v, src.vars)
			if expanded != v {
				src.vars[k] = expanded
				changed = true
			}
		}
		src.mu.Unlock()
		if !changed {
			break
		}
	}

	result := &InterpolateResult{}
	src.mu.RLock()
	defer src.mu.RUnlock()
	for k, v := range src.vars {
		if interpolatePattern.MatchString(v) {
			result.Unresolved = append(result.Unresolved, k)
		} else {
			result.Resolved = append(result.Resolved, k)
		}
	}
	return result, nil
}

// InterpolateWith expands variable references in src using an external lookup
// set ext. Values in src that reference keys present in ext are replaced;
// values referencing keys absent from both sets are left unchanged.
//
// Neither src nor ext may be nil.
func InterpolateWith(src *Set, ext *Set) (*InterpolateResult, error) {
	if src == nil {
		return nil, fmt.Errorf("interpolate: source set must not be nil")
	}
	if ext == nil {
		return nil, fmt.Errorf("interpolate: external set must not be nil")
	}

	ext.mu.RLock()
	lookup := make(map[string]string, len(ext.vars))
	for k, v := range ext.vars {
		lookup[k] = v
	}
	ext.mu.RUnlock()

	// Merge src vars into lookup so self-references also resolve.
	src.mu.Lock()
	for k, v := range src.vars {
		if _, exists := lookup[k]; !exists {
			lookup[k] = v
		}
	}
	for k, v := range src.vars {
		src.vars[k] = expand(v, lookup)
	}
	src.mu.Unlock()

	result := &InterpolateResult{}
	src.mu.RLock()
	defer src.mu.RUnlock()
	for k, v := range src.vars {
		if interpolatePattern.MatchString(v) {
			result.Unresolved = append(result.Unresolved, k)
		} else {
			result.Resolved = append(result.Resolved, k)
		}
	}
	return result, nil
}

// expand performs a single-pass replacement of variable references in s using
// the provided lookup map.
func expand(s string, lookup map[string]string) string {
	return interpolatePattern.ReplaceAllStringFunc(s, func(match string) string {
		// Extract the variable name from either ${VAR} or $VAR form.
		name := strings.TrimPrefix(match, "$")
		name = strings.TrimPrefix(name, "{")
		name = strings.TrimSuffix(name, "}")
		if val, ok := lookup[name]; ok {
			return val
		}
		return match
	})
}
