package env

import (
	"fmt"
	"strings"
)

// ExportFormat defines the output format for exported environment variables.
type ExportFormat string

const (
	FormatExport ExportFormat = "export"
	FormatDotenv ExportFormat = "dotenv"
	FormatJSON   ExportFormat = "json"
)

// Export renders the variables in a Set to a string in the given format.
func Export(s *Set, format ExportFormat) (string, error) {
	if s == nil {
		return "", fmt.Errorf("set must not be nil")
	}

	switch format {
	case FormatExport:
		return exportShell(s), nil
	case FormatDotenv:
		return exportDotenv(s), nil
	case FormatJSON:
		return exportJSON(s), nil
	default:
		return "", fmt.Errorf("unknown format: %q", format)
	}
}

func exportShell(s *Set) string {
	var sb strings.Builder
	for k, v := range s.Vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}

func exportDotenv(s *Set) string {
	var sb strings.Builder
	for k, v := range s.Vars {
		fmt.Fprintf(&sb, "%s=%s\n", k, v)
	}
	return sb.String()
}

func exportJSON(s *Set) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	i := 0
	for k, v := range s.Vars {
		comma := ","
		if i == len(s.Vars)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, v, comma)
		i++
	}
	sb.WriteString("}\n")
	return sb.String()
}
