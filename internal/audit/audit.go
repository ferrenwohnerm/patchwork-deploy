package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp   time.Time `json:"timestamp"`
	Environment string    `json:"environment"`
	Patch       string    `json:"patch"`
	Action      string    `json:"action"`
	Success     bool      `json:"success"`
	Message     string    `json:"message,omitempty"`
}

// Logger writes audit entries to a newline-delimited JSON file.
type Logger struct {
	path string
}

// New returns a Logger that appends to the audit log at dir/audit.log.
func New(dir string) *Logger {
	return &Logger{path: filepath.Join(dir, "audit.log")}
}

// Record appends an Entry to the audit log.
func (l *Logger) Record(env, patch, action string, success bool, msg string) error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp:   time.Now().UTC(),
		Environment: env,
		Patch:       patch,
		Action:      action,
		Success:     success,
		Message:     msg,
	}
	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// Read returns all entries from the audit log.
func (l *Logger) Read() ([]Entry, error) {
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
