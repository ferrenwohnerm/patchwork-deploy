package env

import (
	"strings"
	"testing"
)

func TestFormatQuotaResult_NoQuotas(t *testing.T) {
	out := FormatQuotaResult("staging", nil, nil)
	if !strings.Contains(out, "no quotas") {
		t.Fatalf("expected 'no quotas' message, got: %s", out)
	}
}

func TestFormatQuotaResult_WithEntries(t *testing.T) {
	entries := map[string]int{"001-init.sql": 3, "002-seed.sql": 1}
	counts := map[string]int{"001-init.sql": 1, "002-seed.sql": 1}
	out := FormatQuotaResult("prod", entries, counts)
	if !strings.Contains(out, "001-init.sql") {
		t.Errorf("expected patch name in output")
	}
	if !strings.Contains(out, "limit=3") {
		t.Errorf("expected limit in output")
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("expected ok status")
	}
}

func TestFormatQuotaResult_ExceededStatus(t *testing.T) {
	entries := map[string]int{"001-init.sql": 2}
	counts := map[string]int{"001-init.sql": 2}
	out := FormatQuotaResult("prod", entries, counts)
	if !strings.Contains(out, "EXCEEDED") {
		t.Errorf("expected EXCEEDED status, got: %s", out)
	}
}

func TestFormatQuotaResult_SortedOutput(t *testing.T) {
	entries := map[string]int{"zzz.sql": 1, "aaa.sql": 1}
	counts := map[string]int{}
	out := FormatQuotaResult("dev", entries, counts)
	idxA := strings.Index(out, "aaa.sql")
	idxZ := strings.Index(out, "zzz.sql")
	if idxA > idxZ {
		t.Errorf("expected sorted output: aaa before zzz")
	}
}
