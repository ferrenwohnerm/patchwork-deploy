package env

import (
	"fmt"
	"github.com/patchwork-deploy/internal/state"
)

const quotaSentinelPrefix = "__quota__"

type QuotaConfig struct {
	Environment string
	MaxPatches  int
}

type QuotaResult struct {
	Environment string
	Applied     int
	Limit       int
	Exceeded    bool
}

func SetQuota(st *state.State, env string, max int) error {
	if env == "" {
		return fmt.Errorf("environment name required")
	}
	if max < 1 {
		return fmt.Errorf("quota must be at least 1")
	}
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	key := quotaSentinelPrefix + env
	st.SetMeta(key, fmt.Sprintf("%d", max))
	return nil
}

func GetQuota(st *state.State, env string) (int, bool) {
	key := quotaSentinelPrefix + env
	val, ok := st.GetMeta(key)
	if !ok {
		return 0, false
	}
	var n int
	fmt.Sscanf(val, "%d", &n)
	return n, true
}

func CheckQuota(st *state.State, env string) (QuotaResult, error) {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return QuotaResult{}, fmt.Errorf("environment %q not found", env)
	}
	limit, hasQuota := GetQuota(st, env)
	result := QuotaResult{
		Environment: env,
		Applied:     len(records),
		Limit:       limit,
		Exceeded:    hasQuota && len(records) > limit,
	}
	return result, nil
}

func RemoveQuota(st *state.State, env string) error {
	key := quotaSentinelPrefix + env
	st.DeleteMeta(key)
	return nil
}
