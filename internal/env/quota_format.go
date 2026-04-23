package env

import (
	"fmt"
	"sort"
	"strings"
)

// FormatQuotaResult formats quota information for display.
func FormatQuotaResult(env string, entries map[string]int, counts map[string]int) string {
	if len(entries) == 0 {
		return fmt.Sprintf("no quotas set for environment %q", env)
	}

	keys := make([]string, 0, len(entries))
	for k := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("quotas for %s:\n", env))
	for _, patch := range keys {
		limit := entries[patch]
		used := counts[patch]
		status := "ok"
		if used >= limit {
			status = "EXCEEDED"
		}
		sb.WriteString(fmt.Sprintf("  %-30s limit=%-4d used=%-4d [%s]\n", patch, limit, used, status))
	}
	return strings.TrimRight(sb.String(), "\n")
}
