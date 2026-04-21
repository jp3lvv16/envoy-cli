package env

import (
	"fmt"
	"sort"
	"strings"
)

// FormatCompare returns a human-readable summary of a CompareResult.
// Each section is prefixed with a label and lists key=value pairs in sorted order.
func FormatCompare(cr *CompareResult) string {
	if cr == nil {
		return ""
	}
	var sb strings.Builder

	if len(cr.Same) > 0 {
		sb.WriteString("=== Same ===\n")
		for _, k := range sortedMapKeys(cr.Same) {
			fmt.Fprintf(&sb, "  %s=%s\n", k, cr.Same[k])
		}
	}

	if len(cr.OnlyInSrc) > 0 {
		sb.WriteString("=== Only in source ===\n")
		for _, k := range sortedMapKeys(cr.OnlyInSrc) {
			fmt.Fprintf(&sb, "  + %s=%s\n", k, cr.OnlyInSrc[k])
		}
	}

	if len(cr.OnlyInDst) > 0 {
		sb.WriteString("=== Only in destination ===\n")
		for _, k := range sortedMapKeys(cr.OnlyInDst) {
			fmt.Fprintf(&sb, "  - %s=%s\n", k, cr.OnlyInDst[k])
		}
	}

	if len(cr.Conflicted) > 0 {
		sb.WriteString("=== Conflicted ===\n")
		keys := make([]string, 0, len(cr.Conflicted))
		for k := range cr.Conflicted {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			pair := cr.Conflicted[k]
			fmt.Fprintf(&sb, "  ~ %s: %q -> %q\n", k, pair[0], pair[1])
		}
	}

	return sb.String()
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
