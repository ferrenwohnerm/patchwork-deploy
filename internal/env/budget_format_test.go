package env

import (
	"strings"
	"testing"
)

func TestFormatBudgets_NoEntries(t *testing.T) {
	out := FormatBudgets([]BudgetEntry{})
	if out != "no budgets set" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatBudgets_WithEntries(t *testing.T) {
	entries := []BudgetEntry{
		{Patch: "001-init.sql", Limit: 3, Used: 1, Exceeded: false},
		{Patch: "002-index.sql", Limit: 2, Used: 2, Exceeded: true},
	}
	out := FormatBudgets(entries)
	if !strings.Contains(out, "001-init.sql") {
		t.Error("expected 001-init.sql in output")
	}
	if !strings.Contains(out, "EXCEEDED") {
		t.Error("expected EXCEEDED status in output")
	}
	if !strings.Contains(out, "ok") {
		t.Error("expected ok status in output")
	}
}

func TestFormatBudgets_SortedOutput(t *testing.T) {
	entries := CollectBudgetEntries(
		map[string]int{"002-index.sql": 5, "001-init.sql": 3},
		map[string]int{"001-init.sql": 1},
	)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Patch != "001-init.sql" {
		t.Errorf("expected sorted first entry to be 001-init.sql, got %s", entries[0].Patch)
	}
}

func TestFormatBudgets_ExceededStatus(t *testing.T) {
	entries := []BudgetEntry{
		{Patch: "001-init.sql", Limit: 1, Used: 1, Exceeded: true},
	}
	out := FormatBudgets(entries)
	if !strings.Contains(out, "EXCEEDED") {
		t.Error("expected EXCEEDED in output")
	}
}
