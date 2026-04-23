package env

import (
	"fmt"
	"sort"
	"strings"
)

// FormatRelays produces a human-readable table of relay mappings.
func FormatRelays(env string, relays map[string]string) string {
	if len(relays) == 0 {
		return fmt.Sprintf("no relays configured for environment %q\n", env)
	}

	patches := make([]string, 0, len(relays))
	for p := range relays {
		patches = append(patches, p)
	}
	sort.Strings(patches)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("relays for %q:\n", env))
	for _, p := range patches {
		sb.WriteString(fmt.Sprintf("  %-30s -> %s\n", p, relays[p]))
	}
	return sb.String()
}
