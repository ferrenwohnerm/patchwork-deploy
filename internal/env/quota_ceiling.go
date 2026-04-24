package env

import (
	"fmt"

	"github.com/patchwork-deploy/internal/state"
)

// Watermark represents a soft warning threshold for patch counts.
const watermarkKeyPrefix = "watermark:"

func watermarkKey(env string) string {
	return watermarkKeyPrefix + env
}

// SetWatermark sets a soft warning threshold for the number of applied patches
// in an environment. When the count reaches or exceeds the watermark, callers
// should emit a warning but are not blocked.
func SetWatermark(st *state.State, env string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if limit <= 0 {
		return fmt.Errorf("watermark limit must be greater than zero")
	}
	st.SetMeta(watermarkKey(env), fmt.Sprintf("%d", limit))
	return nil
}

// GetWatermark returns the watermark limit for an environment, or 0 if not set.
func GetWatermark(st *state.State, env string) int {
	v := st.GetMeta(watermarkKey(env))
	if v == "" {
		return 0
	}
	var n int
	fmt.Sscanf(v, "%d", &n)
	return n
}

// ClearWatermark removes the watermark for an environment.
func ClearWatermark(st *state.State, env string) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	st.DeleteMeta(watermarkKey(env))
	return nil
}

// CheckWatermark returns true if the current patch count meets or exceeds the
// watermark. Returns false if no watermark is set.
func CheckWatermark(st *state.State, env string) (exceeded bool, current int, limit int) {
	limit = GetWatermark(st, env)
	if limit == 0 {
		return false, 0, 0
	}
	current = len(st.ForEnvironment(env))
	return current >= limit, current, limit
}
