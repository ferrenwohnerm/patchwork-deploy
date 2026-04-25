package env

import (
	"fmt"
	"sort"
	"strings"
)

// FormatChains renders a chain map as a human-readable table.
func FormatChains(chains map[string]string) string {
	if len(chains) == 0 {
		return "no chains defined"
	}

	patches := make([]string, 0, len(chains))
	for p := range chains {
		patches = append(patches, p)
	}
	sort.Strings(patches)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %s\n", "PATCH", "SUCCESSOR"))
	sb.WriteString(strings.Repeat("-", 55) + "\n")
	for _, p := range patches {
		sb.WriteString(fmt.Sprintf("%-30s %s\n", p, chains[p]))
	}
	return strings.TrimRight(sb.String(), "\n")
}
