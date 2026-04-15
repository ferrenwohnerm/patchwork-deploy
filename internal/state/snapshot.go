package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time copy of the state for a given environment.
type Snapshot struct {
	Environment string    `json:"environment"`
	TakenAt     time.Time `json:"taken_at"`
	Records     []Record  `json:"records"`
}

// TakeSnapshot writes a JSON snapshot of the current state for the given
// environment into snapshotDir. The file is named
// <env>-<unix-timestamp>.snapshot.json so multiple snapshots can coexist.
func TakeSnapshot(s *State, environment, snapshotDir string) (string, error) {
	records := s.ForEnvironment(environment)

	snap := Snapshot{
		Environment: environment,
		TakenAt:     time.Now().UTC(),
		Records:     records,
	}

	if err := os.MkdirAll(snapshotDir, 0o755); err != nil {
		return "", fmt.Errorf("create snapshot dir: %w", err)
	}

	filename := fmt.Sprintf("%s-%d.snapshot.json", environment, snap.TakenAt.Unix())
	path := filepath.Join(snapshotDir, filename)

	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create snapshot file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return "", fmt.Errorf("encode snapshot: %w", err)
	}

	return path, nil
}

// LoadSnapshot reads a snapshot file from disk and returns the Snapshot.
func LoadSnapshot(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open snapshot: %w", err)
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, fmt.Errorf("decode snapshot: %w", err)
	}
	return &snap, nil
}
