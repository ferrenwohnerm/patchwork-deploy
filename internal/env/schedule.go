package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

type ScheduleEntry struct {
	Environment string
	Patch       string
	RunAt       time.Time
	CreatedAt   time.Time
}

const schedulePrefix = "__schedule__"

func Schedule(st *state.State, env, patch string, runAt time.Time) error {
	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q not found", env)
	}
	found := false
	for _, r := range records {
		if r.Patch == patch {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("patch %q not found in environment %q", patch, env)
	}
	key := schedulePrefix + env
	st.Add(key, state.Record{
		Environment: key,
		Patch:       patch,
		AppliedAt:   runAt,
	})
	return nil
}

func ListScheduled(st *state.State, env string) []ScheduleEntry {
	key := schedulePrefix + env
	records := st.ForEnvironment(key)
	entries := make([]ScheduleEntry, 0, len(records))
	for _, r := range records {
		entries = append(entries, ScheduleEntry{
			Environment: env,
			Patch:       r.Patch,
			RunAt:       r.AppliedAt,
		})
	}
	return entries
}

func CancelScheduled(st *state.State, env, patch string) error {
	key := schedulePrefix + env
	records := st.ForEnvironment(key)
	var kept []state.Record
	found := false
	for _, r := range records {
		if r.Patch == patch {
			found = true
			continue
		}
		kept = append(kept, r)
	}
	if !found {
		return fmt.Errorf("no scheduled entry for patch %q in environment %q", patch, env)
	}
	st.DeleteEnvironment(key)
	for _, r := range kept {
		st.Add(key, r)
	}
	return nil
}
