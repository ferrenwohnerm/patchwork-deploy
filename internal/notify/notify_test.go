package notify

import (
	"strings"
	"testing"
)

func TestSend_FormatsOutput(t *testing.T) {
	var buf strings.Builder
	n := New(&buf)

	e := n.Send("staging", LevelInfo, "deployment started")

	if e.Environment != "staging" {
		t.Errorf("expected environment staging, got %s", e.Environment)
	}
	if e.Level != LevelInfo {
		t.Errorf("expected level INFO, got %s", e.Level)
	}
	if e.Message != "deployment started" {
		t.Errorf("unexpected message: %s", e.Message)
	}

	out := buf.String()
	if !strings.Contains(out, "INFO") {
		t.Errorf("output missing INFO: %q", out)
	}
	if !strings.Contains(out, "staging") {
		t.Errorf("output missing environment: %q", out)
	}
	if !strings.Contains(out, "deployment started") {
		t.Errorf("output missing message: %q", out)
	}
}

func TestInfo_UsesInfoLevel(t *testing.T) {
	var buf strings.Builder
	n := New(&buf)
	e := n.Info("prod", "patch applied")
	if e.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", e.Level)
	}
}

func TestWarn_UsesWarnLevel(t *testing.T) {
	var buf strings.Builder
	n := New(&buf)
	e := n.Warn("dev", "patch skipped")
	if e.Level != LevelWarn {
		t.Errorf("expected WARN, got %s", e.Level)
	}
	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("output missing WARN label")
	}
}

func TestError_UsesErrorLevel(t *testing.T) {
	var buf strings.Builder
	n := New(&buf)
	e := n.Error("prod", "apply failed")
	if e.Level != LevelError {
		t.Errorf("expected ERROR, got %s", e.Level)
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Should not panic when nil is passed.
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer")
	}
}
