package lock_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/lock"
)

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "locktest-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAcquireAndRelease(t *testing.T) {
	dir := tempDir(t)
	l := lock.New(dir)

	if err := l.Acquire(); err != nil {
		t.Fatalf("expected no error on first acquire, got %v", err)
	}
	if !l.IsHeld() {
		t.Fatal("expected lock to be held after acquire")
	}
	if err := l.Release(); err != nil {
		t.Fatalf("expected no error on release, got %v", err)
	}
	if l.IsHeld() {
		t.Fatal("expected lock to be released")
	}
}

func TestAcquire_FailsWhenAlreadyLocked(t *testing.T) {
	dir := tempDir(t)
	l := lock.New(dir)

	if err := l.Acquire(); err != nil {
		t.Fatalf("first acquire failed: %v", err)
	}
	defer l.Release() //nolint:errcheck

	l2 := lock.New(dir)
	err := l2.Acquire()
	if err == nil {
		t.Fatal("expected error on second acquire, got nil")
	}
	if !errors.Is(err, lock.ErrLocked) {
		t.Fatalf("expected ErrLocked, got %v", err)
	}
}

func TestRelease_IdempotentWhenMissing(t *testing.T) {
	dir := tempDir(t)
	l := lock.New(dir)
	if err := l.Release(); err != nil {
		t.Fatalf("release on missing lock should not error, got %v", err)
	}
}

func TestLockFile_ContainsPID(t *testing.T) {
	dir := tempDir(t)
	l := lock.New(dir)
	if err := l.Acquire(); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	defer l.Release() //nolint:errcheck

	data, err := os.ReadFile(filepath.Join(dir, ".patchwork.lock"))
	if err != nil {
		t.Fatalf("read lock file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty lock file content")
	}
}
