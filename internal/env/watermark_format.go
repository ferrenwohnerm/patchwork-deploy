package env

import (
	"fmt"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

// WatermarkEntry holds display information for a single environment watermark.
type WatermarkEntry struct {
	Environment string
	Current     int
	Limit       int
	Exceeded    bool
}

// CollectWatermarks gathers watermark status for all environments that have a
// watermark configured.
func CollectWatermarks(st *state.State, envs []string) []WatermarkEntry {
	var out []WatermarkEntry
	for _, e := range envs {
		limit := GetWatermark(st, e)
		if limit == 0 {
			continue
		}
		exceeded, current, _ := CheckWatermark(st, e)
		out = append(out, WatermarkEntry{
			Environment: e,
			Current:     current,
			Limit:       limit,
			Exceeded:    exceeded,
		})
	}
	return out
}

// FormatWatermarks renders a human-readable table of watermark statuses.
func FormatWatermarks(entries []WatermarkEntry) string {
	if len(entries) == 0 {
		return "no watermarks configured"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %8s %8s %s\n", "ENVIRONMENT", "CURRENT", "LIMIT", "STATUS"))
	for _, e := range entries {
		status := "ok"
		if e.Exceeded {
			status = "WARNING"
		}
		sb.WriteString(fmt.Sprintf("%-20s %8d %8d %s\n", e.Environment, e.Current, e.Limit, status))
	}
	return strings.TrimRight(sb.String(), "\n")
}
