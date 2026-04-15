package runner_test

import (
	"errors"
	"testing"

	"github.com/patchwork-deploy/internal/runner"
)

func TestSummarise_Counts(t *testing.T) {
	results := []runner.Result{
		{PatchID: "001", Applied: true},
		{PatchID: "002", Skipped: true},
		{PatchID: "003", Err: errors.New("boom")},
		{PatchID: "004", Applied: true},
	}

	s := runner.Summarise(results)

	if s.Total != 4 {
		t.Errorf("Total: want 4, got %d", s.Total)
	}
	if s.Applied != 2 {
		t.Errorf("Applied: want 2, got %d", s.Applied)
	}
	if s.Skipped != 1 {
		t.Errorf("Skipped: want 1, got %d", s.Skipped)
	}
	if s.Failed != 1 {
		t.Errorf("Failed: want 1, got %d", s.Failed)
	}
}

func TestSummarise_Empty(t *testing.T) {
	s := runner.Summarise(nil)
	if s.Total != 0 || s.Applied != 0 || s.Skipped != 0 || s.Failed != 0 {
		t.Errorf("expected all-zero summary, got %+v", s)
	}
}

func TestSummary_String(t *testing.T) {
	s := runner.Summary{Total: 3, Applied: 1, Skipped: 1, Failed: 1}
	want := "total=3 applied=1 skipped=1 failed=1"
	if got := s.String(); got != want {
		t.Errorf("String(): want %q, got %q", want, got)
	}
}
