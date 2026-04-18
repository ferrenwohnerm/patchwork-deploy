package env

import (
	"fmt"
	"time"

	"github.com/patchwork-deploy/internal/state"
)

const checkpointKeyPrefix = "checkpoint:"

// CheckpointEntry records a named point-in-time marker for an environment.
type CheckpointEntry struct {
	Name      string
	Patch     string
	CreatedAt time.Time
}

// CreateCheckpoint saves a named checkpoint for the latest applied patch in env.
func CreateCheckpoint(st *state.State, env, name string) error {
	if env == "" {
		return fmt.Errorf("environment name must not be empty")
	}
	if name == "" {
		return fmt.Errorf("checkpoint name must not be empty")
	}

	records := st.ForEnvironment(env)
	if len(records) == 0 {
		return fmt.Errorf("environment %q has no applied patches", env)
	}

	// Use the most recently applied patch as the checkpoint target.
	latest := records[len(records)-1]

	key := checkpointKeyPrefix + env
	entry := CheckpointEntry{
		Name:      name,
		Patch:     latest.Patch,
		CreatedAt: time.Now().UTC(),
	}

	existing := listCheckpoints(st, env)
	for _, cp := range existing {
		if cp.Name == name {
			return fmt.Errorf("checkpoint %q already exists for environment %q", name, env)
		}
	}

	st.Add(state.Record{
		Environment: key,
		Patch:       entry.Name + ":" + entry.Patch,
		AppliedAt:   entry.CreatedAt,
	})
	return nil
}

// ListCheckpoints returns all named checkpoints for the given environment.
func ListCheckpoints(st *state.State, env string) ([]CheckpointEntry, error) {
	if env == "" {
		return nil, fmt.Errorf("environment name must not be empty")
	}
	return listCheckpoints(st, env), nil
}

// RemoveCheckpoint deletes a named checkpoint from an environment.
func RemoveCheckpoint(st *state.State, env, name string) error {
	if env == "" {
		return fmt.Errorf("environment name must not be empty")
	}
	if name == "" {
		return fmt.Errorf("checkpoint name must not be empty")
	}

	key := checkpointKeyPrefix + env
	records := st.ForEnvironment(key)
	prefix := name + ":"
	removed := false

	st.DeleteEnvironment(key)
	for _, r := range records {
		if len(r.Patch) >= len(prefix) && r.Patch[:len(prefix)] == prefix {
			removed = true
			continue
		}
		st.Add(r)
	}

	if !removed {
		return fmt.Errorf("checkpoint %q not found for environment %q", name, env)
	}
	return nil
}

// listCheckpoints is the internal helper that parses checkpoint records.
func listCheckpoints(st *state.State, env string) []CheckpointEntry {
	key := checkpointKeyPrefix + env
	records := st.ForEnvironment(key)
	out := make([]CheckpointEntry, 0, len(records))
	for _, r := range records {
		name, patch := splitCheckpointPatch(r.Patch)
		out = append(out, CheckpointEntry{
			Name:      name,
			Patch:     patch,
			CreatedAt: r.AppliedAt,
		})
	}
	return out
}

// splitCheckpointPatch splits "name:patch" into its two components.
func splitCheckpointPatch(s string) (string, string) {
	for i := 0; i < len(s); i++ {
		if s[i] == ':' {
			return s[:i], s[i+1:]
		}
	}
	return s, ""
}
