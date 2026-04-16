package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/patchwork-deploy/internal/audit"
)

func tempAuditDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "audit-cmd-*")
	if err != nil {
		t.Fatalf("tempAuditDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAuditLog_WrittenAndReadBack(t *testing.T) {
	dir := tempAuditDir(t)
	l := audit.New(dir)

	_ = l.Record("prod", "001_schema.sql", "apply", true, "")
	_ = l.Record("prod", "002_data.sql", "apply", false, "timeout")

	entries, err := l.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2, got %d", len(entries))
	}
}

func TestAuditLog_FileCreatedOnDisk(t *testing.T) {
	dir := tempAuditDir(t)
	l := audit.New(dir)

	_ = l.Record("dev", "patch.sql", "apply", true, "")

	if _, err := os.Stat(filepath.Join(dir, "audit.log")); os.IsNotExist(err) {
		t.Error("expected audit.log to exist on disk")
	}
}

func TestAuditLog_FilterByEnv(t *testing.T) {
	dir := tempAuditDir(t)
	l := audit.New(dir)

	_ = l.Record("prod", "001.sql", "apply", true, "")
	_ = l.Record("staging", "001.sql", "apply", true, "")
	_ = l.Record("prod", "002.sql", "apply", true, "")

	entries, _ := l.Read()
	var prod []string
	for _, e := range entries {
		if e.Environment == "prod" {
			prod = append(prod, e.Patch)
		}
	}
	if len(prod) != 2 {
		t.Errorf("expected 2 prod entries, got %d", len(prod))
	}
	for _, p := range prod {
		if !strings.HasSuffix(p, ".sql") {
			t.Errorf("unexpected patch name: %s", p)
		}
	}
}

func TestAuditLog_FailedEntryHasReason(t *testing.T) {
	dir := tempAuditDir(t)
	l := audit.New(dir)

	_ = l.Record("prod", "003_index.sql", "apply", false, "connection refused")

	entries, err := l.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Success {
		t.Error("expected entry to be marked as failed")
	}
	if e.Reason != "connection refused" {
		t.Errorf("expected reason %q, got %q", "connection refused", e.Reason)
	}
}
