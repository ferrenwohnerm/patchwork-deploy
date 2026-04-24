package env

import (
	"fmt"
	"sort"
	"strings"
)

// FormatWindows renders a deployment window map as a human-readable table.
func FormatWindows(windows map[string][2]string) string {
	if len(windows) == 0 {
		return "no deployment windows set"
	}

	patches := make([]string, 0, len(windows))
	for p := range windows {
		patches = append(patches, p)
	}
	sort.Strings(patches)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s  %s\n", "PATCH", "WINDOW"))
	sb.WriteString(strings.Repeat("-", 45) + "\n")
	for _, p := range patches {
		w := windows[p]
		sb.WriteString(fmt.Sprintf("%-30s  %s - %s\n", p, w[0], w[1]))
	}
	return strings.TrimRight(sb.String(), "\n")
}
