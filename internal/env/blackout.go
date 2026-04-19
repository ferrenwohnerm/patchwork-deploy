package env

import (
	"fmt"
	"strings"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const blackoutPrefix = "blackout:"

// BlackoutWindow represents a time range during which patches are blocked.
type BlackoutWindow struct {
	Start time.Time
	End   time.Time
	Note  string
}

// SetBlackout records a blackout window for an environment.
func SetBlackout(st *state.State, env string, start, end time.Time, note string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !end.After(start) {
		return fmt.Errorf("end time must be after start time")
	}
	key := fmt.Sprintf("%s%s", blackoutPrefix, env)
	value := fmt.Sprintf("%s|%s|%s", start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339), note)
	st.SetMeta(key, value)
	return nil
}

// GetBlackout returns the active blackout window for an environment, if any.
func GetBlackout(st *state.State, env string) (*BlackoutWindow, error) {
	if !st.HasEnvironment(env) {
		return nil, fmt.Errorf("environment %q not found", env)
	}
	key := fmt.Sprintf("%s%s", blackoutPrefix, env)
	val, ok := st.GetMeta(key)
	if !ok {
		return nil, nil
	}
	parts := strings.SplitN(val, "|", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("malformed blackout entry for %q", env)
	}
	start, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid start time: %w", err)
	}
	end, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid end time: %w", err)
	}
	return &BlackoutWindow{Start: start, End: end, Note: parts[2]}, nil
}

// ClearBlackout removes the blackout window for an environment.
func ClearBlackout(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	key := fmt.Sprintf("%s%s", blackoutPrefix, env)
	st.DeleteMeta(key)
	return nil
}

// IsBlackedOut returns true if the environment is currently in a blackout window.
func IsBlackedOut(st *state.State, env string, now time.Time) (bool, error) {
	bw, err := GetBlackout(st, env)
	if err != nil || bw == nil {
		return false, err
	}
	return !now.Before(bw.Start) && now.Before(bw.End), nil
}
