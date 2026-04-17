package env

import (
	"strings"
	"testing"
)

func TestFormatPromoteResult_AllApplied(t *testing.T) {
	r := PromoteResult{
		Source:  "staging",
		Target:  "production",
		Applied: []string{"001-init.sql", "002-add-index.sql"},
		Skipped: nil,
	}
	out := FormatPromoteResult(r)
	if !strings.Contains(out, "staging → production") {
		t.Errorf("expected source/target header, got: %s", out)
	}
	if !strings.Contains(out, "+ 001-init.sql") {
		t.Errorf("expected applied patch listed, got: %s", out)
	}
	if !strings.Contains(out, "2 applied, 0 skipped") {
		t.Errorf("expected summary counts, got: %s", out)
	}
}

func TestFormatPromoteResult_WithSkipped(t *testing.T) {
	r := PromoteResult{
		Source:  "dev",
		Target:  "staging",
		Applied: []string{"003-new.sql"},
		Skipped: []string{"001-init.sql"},
	}
	out := FormatPromoteResult(r)
	if !strings.Contains(out, "~ 001-init.sql (already applied)") {
		t.Errorf("expected skipped patch, got: %s", out)
	}
	if !strings.Contains(out, "1 applied, 1 skipped") {
		t.Errorf("expected summary counts, got: %s", out)
	}
}

func TestFormatPromoteResult_Empty(t *testing.T) {
	r := PromoteResult{
		Source: "dev",
		Target: "staging",
	}
	out := FormatPromoteResult(r)
	if !strings.Contains(out, "no patches to promote") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
