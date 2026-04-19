package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

type RolloutStrategy string

const (
	RolloutCanary     RolloutStrategy = "canary"
	RolloutBlueGreen  RolloutStrategy = "blue-green"
	RolloutImmediate  RolloutStrategy = "immediate"
	rolloutPrefix                     = "rollout:"
)

type RolloutConfig struct {
	Strategy  RolloutStrategy
	Patch     string
	CreatedAt time.Time
}

func SetRollout(st *state.State, env, patch string, strategy RolloutStrategy) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(records, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if strategy != RolloutCanary && strategy != RolloutBlueGreen && strategy != RolloutImmediate {
		return fmt.Errorf("unknown strategy %q: must be canary, blue-green, or immediate", strategy)
	}
	key := rolloutPrefix + env + ":" + patch
	st.SetMeta(key, string(strategy))
	return nil
}

func GetRollout(st *state.State, env, patch string) (RolloutConfig, bool) {
	key := rolloutPrefix + env + ":" + patch
	val, ok := st.GetMeta(key)
	if !ok {
		return RolloutConfig{}, false
	}
	return RolloutConfig{
		Strategy: RolloutStrategy(val),
		Patch:    patch,
	}, true
}

func ClearRollout(st *state.State, env, patch string) {
	key := rolloutPrefix + env + ":" + patch
	st.DeleteMeta(key)
}

func ListRollouts(st *state.State, env string) []RolloutConfig {
	var out []RolloutConfig
	prefix := rolloutPrefix + env + ":"
	for _, kv := range st.AllMeta() {
		if len(kv.Key) > len(prefix) && kv.Key[:len(prefix)] == prefix {
			patch := kv.Key[len(prefix):]
			out = append(out, RolloutConfig{Strategy: RolloutStrategy(kv.Value), Patch: patch})
		}
	}
	return out
}

func patchExistsInEnv(records []state.Record, patch string) bool {
	for _, r := range records {
		if r.Patch == patch {
			return true
		}
	}
	return false
}
