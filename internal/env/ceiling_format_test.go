package env

import (
	"strings"
	"testing"
)

func TestFormatCeilings_NoEntries(t *testing.T) {
	out := FormatCeilings(nil)
	if out != "no ceilings configured" {
		t.Fatalf("expected empty message, got %q", out)
	}
}

func TestFormatCeilings_WithEntries(t *testing.T) {
	results := []CeilingResult{
		{Patch: "002-add-index", Limit: 5, Current: 2, Exceeded: false},
		{Patch: "001-init", Limit: 3, Current: 3, Exceeded: false},
	}
	out := FormatCeilings(results)
	if !strings.Contains(out, "001-init") {
		t.Errorf("expected 001-init in output")
	}
	if !strings.Contains(out, "002-add-index") {
		t.Errorf("expected 002-add-index in output")
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// header + 2 entries
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestFormatCeilings_ExceededStatus(t *testing.T) {
	results := []CeilingResult{
		{Patch: "003-migrate", Limit: 2, Current: 5, Exceeded: true},
	}
	out := FormatCeilings(results)
	if !strings.Contains(out, "EXCEEDED") {
		t.Errorf("expected EXCEEDED in output, got: %s", out)
	}
}

func TestFormatCeilings_SortedOutput(t *testing.T) {
	results := []CeilingResult{
		{Patch: "zzz-last", Limit: 1, Current: 0, Exceeded: false},
		{Patch: "aaa-first", Limit: 1, Current: 0, Exceeded: false},
	}
	out := FormatCeilings(results)
	idxA := strings.Index(out, "aaa-first")
	idxZ := strings.Index(out, "zzz-last")
	if idxA > idxZ {
		t.Errorf("expected aaa-first before zzz-last in output")
	}
}
