package diff_test

import (
	"strings"
	"testing"

	"github.com/patchwork-deploy/internal/diff"
)

func recs(env string, patches ...string) []diff.Record {
	var out []diff.Record
	for _, p := range patches {
		out = append(out, diff.Record{Patch: p, Environment: env})
	}
	return out
}

func TestCompare_OnlyInA(t *testing.T) {
	a := recs("prod", "001_init.sql", "002_add_users.sql")
	b := recs("prod", "001_init.sql")

	r := diff.Compare(a, b)

	if len(r.OnlyInA) != 1 || r.OnlyInA[0].Patch != "002_add_users.sql" {
		t.Errorf("expected 002_add_users.sql only in A, got %+v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 0 {
		t.Errorf("expected nothing only in B, got %+v", r.OnlyInB)
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := recs("staging", "001_init.sql")
	b := recs("staging", "001_init.sql", "003_indexes.sql")

	r := diff.Compare(a, b)

	if len(r.OnlyInB) != 1 || r.OnlyInB[0].Patch != "003_indexes.sql" {
		t.Errorf("expected 003_indexes.sql only in B, got %+v", r.OnlyInB)
	}
}

func TestCompare_InBoth(t *testing.T) {
	a := recs("dev", "001_init.sql", "002_add_users.sql")
	b := recs("dev", "001_init.sql", "002_add_users.sql")

	r := diff.Compare(a, b)

	if len(r.InBoth) != 2 {
		t.Errorf("expected 2 in both, got %d", len(r.InBoth))
	}
	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Error("expected no exclusive records")
	}
}

func TestFormat_NoDifferences(t *testing.T) {
	a := recs("prod", "001_init.sql")
	r := diff.Compare(a, a)
	out := diff.Format(r, "snapshotA", "snapshotB")
	if !strings.Contains(out, "no differences") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestFormat_ShowsPrefixes(t *testing.T) {
	a := recs("prod", "001_init.sql")
	b := recs("prod", "002_new.sql")
	r := diff.Compare(a, b)
	out := diff.Format(r, "A", "B")
	if !strings.Contains(out, "- ") {
		t.Errorf("expected removal line, got: %s", out)
	}
	if !strings.Contains(out, "+ ") {
		t.Errorf("expected addition line, got: %s", out)
	}
}
