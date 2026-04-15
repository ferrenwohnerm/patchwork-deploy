package patch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ApplyRecord tracks a single applied patch for a given environment.
type ApplyRecord struct {
	PatchName   string
	Environment string
	AppliedAt   time.Time
}

// Applier handles applying patches and recording their state.
type Applier struct {
	loader  *Loader
	records []ApplyRecord
}

// NewApplier creates an Applier backed by the given Loader.
func NewApplier(loader *Loader) *Applier {
	return &Applier{loader: loader}
}

// Apply reads the patch file and writes its contents to the destination path,
// then records the application.
func (a *Applier) Apply(patchName, environment, destDir string) error {
	if !a.loader.Exists(patchName) {
		return fmt.Errorf("patch %q not found", patchName)
	}

	src := filepath.Join(a.loader.dir, patchName)
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("reading patch %q: %w", patchName, err)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return fmt.Errorf("creating dest dir: %w", err)
	}

	dest := filepath.Join(destDir, patchName)
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return fmt.Errorf("writing patch to %q: %w", dest, err)
	}

	a.records = append(a.records, ApplyRecord{
		PatchName:   patchName,
		Environment: environment,
		AppliedAt:   time.Now(),
	})
	return nil
}

// Applied returns all recorded apply operations.
func (a *Applier) Applied() []ApplyRecord {
	out := make([]ApplyRecord, len(a.records))
	copy(out, a.records)
	return out
}

// WasApplied reports whether patchName was applied to the given environment.
func (a *Applier) WasApplied(patchName, environment string) bool {
	for _, r := range a.records {
		if r.PatchName == patchName && r.Environment == environment {
			return true
		}
	}
	return false
}
