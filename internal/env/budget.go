package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/patchwork-deploy/internal/state"
)

func budgetKey(env, patch string) string {
	return fmt.Sprintf("budget:%s:%s", env, patch)
}

// SetBudget assigns a maximum allowed apply count for a patch in an environment.
func SetBudget(st *state.State, env, patch string, limit int) error {
	if !st.HasEnvironment(env) {
		return fmt.Errorf("environment %q not found", env)
	}
	if !patchExistsInEnv(st, env, patch) {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	if limit <= 0 {
		return fmt.Errorf("budget limit must be greater than zero")
	}
	st.SetMeta(budgetKey(env, patch), strconv.Itoa(limit))
	return nil
}

// GetBudget returns the budget limit for a patch, and whether one is set.
func GetBudget(st *state.State, env, patch string) (int, bool) {
	v, ok := st.GetMeta(budgetKey(env, patch))
	if !ok {
		return 0, false
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return n, true
}

// ClearBudget removes the budget limit for a patch.
func ClearBudget(st *state.State, env, patch string) {
	st.DeleteMeta(budgetKey(env, patch))
}

// CheckBudget returns an error if applying the patch would exceed its budget.
func CheckBudget(st *state.State, env, patch string) error {
	limit, ok := GetBudget(st, env, patch)
	if !ok {
		return nil
	}
	records := st.ForEnvironment(env)
	count := 0
	for _, r := range records {
		if r.Patch == patch {
			count++
		}
	}
	if count >= limit {
		return fmt.Errorf("patch %q has reached its budget of %d in environment %q", patch, limit, env)
	}
	return nil
}

// ListBudgets returns all budget entries for an environment as patch->limit map.
func ListBudgets(st *state.State, env string) map[string]int {
	prefix := fmt.Sprintf("budget:%s:", env)
	result := map[string]int{}
	for k, v := range st.AllMeta() {
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		patch := strings.TrimPrefix(k, prefix)
		if n, err := strconv.Atoi(v); err == nil {
			result[patch] = n
		}
	}
	return result
}
