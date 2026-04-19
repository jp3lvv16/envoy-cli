package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// FormatDiff writes a human-readable diff summary to w.
func FormatDiff(w io.Writer, d *DiffResult) {
	if d.IsEmpty() {
		fmt.Fprintln(w, "No differences.")
		return
	}

	if len(d.Added) > 0 {
		keys := sortedKeys(d.Added)
		for _, k := range keys {
			fmt.Fprintf(w, "+ %s=%s\n", k, d.Added[k])
		}
	}

	if len(d.Removed) > 0 {
		keys := sortedKeys(d.Removed)
		for _, k := range keys {
			fmt.Fprintf(w, "- %s=%s\n", k, d.Removed[k])
		}
	}

	if len(d.Changed) > 0 {
		changedKeys := make([]string, 0, len(d.Changed))
		for k := range d.Changed {
			changedKeys = append(changedKeys, k)
		}
		sort.Strings(changedKeys)
		for _, k := range changedKeys {
			pair := d.Changed[k]
			fmt.Fprintf(w, "~ %s: %s -> %s\n", k, pair[0], pair[1])
		}
	}
}

// Summary returns a one-line summary of the diff.
func (d *DiffResult) Summary() string {
	parts := []string{}
	if n := len(d.Added); n > 0 {
		parts = append(parts, fmt.Sprintf("%d added", n))
	}
	if n := len(d.Removed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d removed", n))
	}
	if n := len(d.Changed); n > 0 {
		parts = append(parts, fmt.Sprintf("%d changed", n))
	}
	if len(parts) == 0 {
		return "no differences"
	}
	return strings.Join(parts, ", ")
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
