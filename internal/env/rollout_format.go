package env

import (
	"fmt"
	"strings"
)

func FormatRollouts(env string, rollouts []RolloutConfig) string {
	if len(rollouts) == 0 {
		return fmt.Sprintf("no rollout strategies set for environment %q", env)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rollout strategies for %q:\n", env))
	for _, r := range rollouts {
		sb.WriteString(fmt.Sprintf("  %-30s %s\n", r.Patch, r.Strategy))
	}
	return strings.TrimRight(sb.String(), "\n")
}
