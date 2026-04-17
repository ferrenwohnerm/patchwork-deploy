package env

import (
	"fmt"
	"strings"
	"time"
)

func FormatSchedule(entries []ScheduleEntry) string {
	if len(entries) == 0 {
		return "no scheduled patches"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s  %s\n", "PATCH", "RUN AT"))
	sb.WriteString(strings.Repeat("-", 52) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-30s  %s\n", e.Patch, e.RunAt.UTC().Format(time.RFC3339)))
	}
	return strings.TrimRight(sb.String(), "\n")
}

func FormatScheduleEntry(e ScheduleEntry) string {
	return fmt.Sprintf("env=%s patch=%s run-at=%s",
		e.Environment, e.Patch, e.RunAt.UTC().Format(time.RFC3339))
}
