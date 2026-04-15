package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

const stateFileName = ".patchwork-state.json"

// Record represents a single applied patch record.
type Record struct {
	PatchID     string    `json:"patch_id"`
	Environment string    `json:"environment"`
	AppliedAt   time.Time `json:"applied_at"`
}

// State holds the full deployment state for a directory.
type State struct {
	Records []Record `json:"records"`
	path    string
}

// Load reads the state file from dir, or returns an empty State if not found.
func Load(dir string) (*State, error) {
	p := filepath.Join(dir, stateFileName)
	s := &State{path: p}

	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}
	return s, nil
}

// Save writes the current state to disk.
func (s *State) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

// Add appends a new record to the state.
func (s *State) Add(patchID, environment string) {
	s.Records = append(s.Records, Record{
		PatchID:     patchID,
		Environment: environment,
		AppliedAt:   time.Now().UTC(),
	})
}

// Has returns true if the given patchID has been applied to environment.
func (s *State) Has(patchID, environment string) bool {
	for _, r := range s.Records {
		if r.PatchID == patchID && r.Environment == environment {
			return true
		}
	}
	return false
}

// ForEnvironment returns all records for a given environment.
func (s *State) ForEnvironment(environment string) []Record {
	var out []Record
	for _, r := range s.Records {
		if r.Environment == environment {
			out = append(out, r)
		}
	}
	return out
}
