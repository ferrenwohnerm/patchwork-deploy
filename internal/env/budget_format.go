package env

import (
	"fmt"
	"sort"
	"strings"
)

// BudgetEntry holds display data for a single budget entry.
type BudgetEntry struct {
	Patch   string
	Limit   int
	Used    int
	Exceeded bool
}

// CollectBudgetEntries builds a sorted slice of BudgetEntry for display.
func CollectBudgetEntries(budgets map[string]int, used map[string]int) []BudgetEntry {
	entries := make([]BudgetEntry, 0, len(budgets))
	for patch, limit := range budgets {
		u := used[patch]
		entries = append(entries, BudgetEntry{
			Patch:    patch,
			Limit:    limit,
			Used:     u,
			Exceeded: u >= limit,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Patch < entries[j].Patch
	})
	return entries
}

// FormatBudgets returns a human-readable table of budget entries.
func FormatBudgets(entries []BudgetEntry) string {
	if len(entries) == 0 {
		return "no budgets set"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %6s %6s %s\n", "PATCH", "USED", "LIMIT", "STATUS"))
	for _, e := range entries {
		status := "ok"
		if e.Exceeded {
			status = "EXCEEDED"
		}
		sb.WriteString(fmt.Sprintf("%-30s %6d %6d %s\n", e.Patch, e.Used, e.Limit, status))
	}
	return strings.TrimRight(sb.String(), "\n")
}
