package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/lock"
)

func tempLockDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "cmdlock-*")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestLockAcquireReleaseCycle(t *testing.T) {
	dir := tempLockDir(t)
	l := lock.New(dir)

	if err := l.Acquire(); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if !l.IsHeld() {
		t.Fatal("should be held after acquire")
	}
	if err := l.Release(); err != nil {
		t.Fatalf("release: %v", err)
	}
	if l.IsHeld() {
		t.Fatal("should not be held after release")
	}
}

func TestLockFile_ExistsOnDisk(t *testing.T) {
	dir := tempLockDir(t)
	l := lock.New(dir)
	if err := l.Acquire(); err != nil {
		t.Fatalf("acquire: %v", err)
	}
	defer l.Release() //nolint:errcheck

	if _, err := os.Stat(filepath.Join(dir, ".patchwork.lock")); os.IsNotExist(err) {
		t.Fatal("expected lock file on disk")
	}
}

func TestLockStatus_Unlocked(t *testing.T) {
	dir := tempLockDir(t)
	l := lock.New(dir)
	if l.IsHeld() {
		t.Fatal("new locker should not be held")
	}
}
