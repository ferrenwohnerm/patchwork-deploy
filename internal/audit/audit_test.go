package audit_test

import (
	"os"
	"testing"

	"github.com/user/patchwork-deploy/internal/audit"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "audit-test-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestRecord_And_Read(t *testing.T) {
	dir := tempDir(t)
	l := audit.New(dir)

	if err := l.Record("prod", "001_init.sql", "apply", true, ""); err != nil {
		t.Fatalf("Record: %v", err)
	}
	if err := l.Record("prod", "002_users.sql", "apply", false, "exec error"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := l.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Patch != "001_init.sql" || !entries[0].Success {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Message != "exec error" || entries[1].Success {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestRead_EmptyWhenMissing(t *testing.T) {
	dir := tempDir(t)
	l := audit.New(dir)

	entries, err := l.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestRecord_AppendsPersists(t *testing.T) {
	dir := tempDir(t)

	for i := 0; i < 3; i++ {
		l := audit.New(dir)
		if err := l.Record("staging", "patch.sql", "apply", true, ""); err != nil {
			t.Fatalf("Record iteration %d: %v", i, err)
		}
	}

	l := audit.New(dir)
	entries, err := l.Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries after 3 appends, got %d", len(entries))
	}
}
