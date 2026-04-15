package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"patchwork-deploy/internal/notify"
)

func executeNotify(args []string) (string, error) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{
		Use:  "notify",
		RunE: runNotify,
	}
	cmd.Flags().StringVarP(&notifyEnv, "env", "e", "", "")
	cmd.Flags().StringVarP(&notifyLevel, "level", "l", "info", "")
	cmd.Flags().StringVarP(&notifyMsg, "message", "m", "", "")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestNotifyCmd_InfoLevel(t *testing.T) {
	out, err := executeNotify([]string{"--env", "staging", "--message", "deploy done"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "event recorded") {
		t.Errorf("expected confirmation in output, got: %q", out)
	}
}

func TestNotifyCmd_WarnLevel(t *testing.T) {
	out, err := executeNotify([]string{"--env", "prod", "--level", "warn", "--message", "slow patch"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "event recorded") {
		t.Errorf("expected confirmation in output, got: %q", out)
	}
}

func TestNotifyCmd_ErrorLevel(t *testing.T) {
	out, err := executeNotify([]string{"--env", "dev", "--level", "error", "--message", "failed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "event recorded") {
		t.Errorf("expected confirmation in output, got: %q", out)
	}
}

func TestNotify_LevelMapping(t *testing.T) {
	cases := []struct {
		input    string
		expected notify.Level
	}{
		{"info", notify.LevelInfo},
		{"warn", notify.LevelWarn},
		{"error", notify.LevelError},
		{"unknown", notify.LevelInfo},
	}
	for _, tc := range cases {
		var lvl notify.Level
		switch tc.input {
		case "warn":
			lvl = notify.LevelWarn
		case "error":
			lvl = notify.LevelError
		default:
			lvl = notify.LevelInfo
		}
		if lvl != tc.expected {
			t.Errorf("input %q: expected %s, got %s", tc.input, tc.expected, lvl)
		}
	}
}
