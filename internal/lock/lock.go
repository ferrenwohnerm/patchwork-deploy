package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const lockFileName = ".patchwork.lock"

// ErrLocked is returned when a lock file already exists.
var ErrLocked = errors.New("another patchwork process is running")

// Locker manages a filesystem-based lock for a given directory.
type Locker struct {
	path string
}

// New creates a Locker rooted at dir.
func New(dir string) *Locker {
	return &Locker{path: filepath.Join(dir, lockFileName)}
}

// Acquire creates the lock file. Returns ErrLocked if it already exists.
func (l *Locker) Acquire() error {
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("%w: lock file %s", ErrLocked, l.path)
		}
		return fmt.Errorf("acquire lock: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "pid=%d ts=%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
	return err
}

// Release removes the lock file.
func (l *Locker) Release() error {
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("release lock: %w", err)
	}
	return nil
}

// IsHeld reports whether the lock file currently exists.
func (l *Locker) IsHeld() bool {
	_, err := os.Stat(l.path)
	return err == nil
}
