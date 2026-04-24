package env

import (
	"fmt"
	"sort"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

// FenceEntry holds the display representation of a single fence record.
type FenceEntry struct {
	Patch   string
	Name    string
	Fenced  bool
}

// CollectFences reads all fence sentinels for the given environment from st
// and returns a sorted slice of FenceEntry values.
func CollectFences(st *state.State, env string) []FenceEntry {
	prefix := fenceKey(env, "")
	var entries []FenceEntry

	for _, rec := range st.Records {
		if !strings.HasPrefix(rec.Key, prefix) {
			continue
		}
		patch := strings.TrimPrefix(rec.Key, prefix)
		if patch == "" {
			continue
		}
		entries = append(entries, FenceEntry{
			Patch:  patch,
			Name:   rec.Value,
			Fenced: true,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Patch < entries[j].Patch
	})
	return entries
}

// FormatFences renders a human-readable table of fence entries.
// Returns a plain string suitable for CLI output.
func FormatFences(entries []FenceEntry) string {
	if len(entries) == 0 {
		return "no fences set"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s  %s\n", "PATCH", "FENCE NAME"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-30s  %s\n", e.Patch, e.Name))
	}

	return strings.TrimRight(sb.String(), "\n")
}
