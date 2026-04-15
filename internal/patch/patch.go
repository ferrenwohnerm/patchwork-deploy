package patch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Patch represents a single incremental deployment patch file.
type Patch struct {
	Name        string
	Environment string
	FilePath    string
	Applied     bool
}

// Loader handles discovery and ordering of patch files for an environment.
type Loader struct {
	PatchDir string
}

// NewLoader creates a Loader that reads patches from the given directory.
func NewLoader(patchDir string) *Loader {
	return &Loader{PatchDir: patchDir}
}

// Discover returns all patch files for the given environment, sorted by name.
// Patch files are expected to follow the pattern: <env>_<name>.yaml
func (l *Loader) Discover(env string) ([]Patch, error) {
	pattern := filepath.Join(l.PatchDir, env+"_*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob error for pattern %q: %w", pattern, err)
	}

	if len(matches) == 0 {
		return []Patch{}, nil
	}

	sort.Strings(matches)

	patches := make([]Patch, 0, len(matches))
	for _, fp := range matches {
		base := filepath.Base(fp)
		name := strings.TrimSuffix(base, ".yaml")
		patches = append(patches, Patch{
			Name:        name,
			Environment: env,
			FilePath:    fp,
		})
	}
	return patches, nil
}

// Exists reports whether a patch file exists on disk.
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
