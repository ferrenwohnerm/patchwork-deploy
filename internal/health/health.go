package health

import (
	"fmt"
	"os"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

// Status represents the overall health of a deployment environment.
type Status struct {
	Environment string
	PatchCount  int
	LastApplied time.Time
	StateFile   string
	StateOK     bool
	LockFile    string
	Locked      bool
}

// String returns a human-readable summary of the health status.
func (s Status) String() string {
	lockStatus := "unlocked"
	if s.Locked {
		lockStatus = "LOCKED"
	}
	stateStatus := "ok"
	if !s.StateOK {
		stateStatus = "ERROR"
	}
	last := "never"
	if !s.LastApplied.IsZero() {
		last = s.LastApplied.Format(time.RFC3339)
	}
	return fmt.Sprintf(
		"env=%-12s patches=%d last_applied=%s state=%s lock=%s",
		s.Environment, s.PatchCount, last, stateStatus, lockStatus,
	)
}

// Check inspects the state and lock files for the given environment and
// working directory, returning a populated Status.
func Check(env, stateFile, lockFile string) Status {
	s := Status{
		Environment: env,
		StateFile:   stateFile,
		LockFile:    lockFile,
	}

	st, err := state.Load(stateFile)
	if err != nil {
		s.StateOK = false
	} else {
		s.StateOK = true
		records := st.ForEnvironment(env)
		s.PatchCount = len(records)
		for _, r := range records {
			if r.AppliedAt.After(s.LastApplied) {
				s.LastApplied = r.AppliedAt
			}
		}
	}

	if _, err := os.Stat(lockFile); err == nil {
		s.Locked = true
	}

	return s
}
