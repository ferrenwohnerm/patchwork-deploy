package state

import (
	"fmt"
	"os"
	"path/filepath"
)

// Import reads patch records from a legacy flat file (one patch ID per line)
// and merges them into the current state for the given environment.
func Import(s *State, env, legacyPath string) (int, error) {
	data, err := os.ReadFile(legacyPath)
	if err != nil {
		return 0, fmt.Errorf("import: reading legacy file: %w", err)
	}

	lines := splitLines(string(data))
	imported := 0
	for _, id := range lines {
		if id == "" {
			continue
		}
		if !s.Has(env, id) {
			s.Records = append(s.Records, Record{
				Environment: env,
				PatchID:     id,
				AppliedAt:   "",
			})
			imported++
		}
	}
	return imported, nil
}

// Prune removes all state records for patch IDs that no longer exist
// in the given patch directory for the specified environment.
func Prune(s *State, env, patchDir string) (int, error) {
	entries, err := os.ReadDir(patchDir)
	if err != nil {
		return 0, fmt.Errorf("prune: reading patch dir: %w", err)
	}

	existing := make(map[string]bool)
	for _, e := range entries {
		if !e.IsDir() {
			existing[stripExt(e.Name())] = true
		}
	}

	var kept []Record
	pruned := 0
	for _, r := range s.Records {
		if r.Environment != env || existing[r.PatchID] {
			kept = append(kept, r)
		} else {
			pruned++
		}
	}
	s.Records = kept
	return pruned, nil
}

func stripExt(name string) string {
	ext := filepath.Ext(name)
	if ext == "" {
		return name
	}
	return name[:len(name)-len(ext)]
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
