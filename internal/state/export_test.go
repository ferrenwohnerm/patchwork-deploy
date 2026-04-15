package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExport_TextFormat(t *testing.T) {
	s := &State{}
	now := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)
	s.Records = []Record{
		{Environment: "prod", Patch: "002-add-index", AppliedAt: now},
		{Environment: "staging", Patch: "001-init", AppliedAt: now},
	}

	out := filepath.Join(t.TempDir(), "export.txt")
	err := Export(s, out, ExportOptions{Format: FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	content := string(data)
	if !strings.Contains(content, "prod\t002-add-index") {
		t.Errorf("expected prod record in output, got:\n%s", content)
	}
	if !strings.Contains(content, "staging\t001-init") {
		t.Errorf("expected staging record in output, got:\n%s", content)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	s := &State{}
	now := time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC)
	s.Records = []Record{
		{Environment: "prod", Patch: "001-init", AppliedAt: now},
	}

	out := filepath.Join(t.TempDir(), "export.json")
	err := Export(s, out, ExportOptions{Format: FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(records) != 1 || records[0].Patch != "001-init" {
		t.Errorf("unexpected records: %+v", records)
	}
}

func TestExport_FilterByEnvironment(t *testing.T) {
	s := &State{}
	now := time.Now()
	s.Records = []Record{
		{Environment: "prod", Patch: "001-init", AppliedAt: now},
		{Environment: "staging", Patch: "001-init", AppliedAt: now},
	}

	out := filepath.Join(t.TempDir(), "export.txt")
	err := Export(s, out, ExportOptions{Environment: "prod", Format: FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(out)
	if strings.Contains(string(data), "staging") {
		t.Errorf("expected only prod records, got staging in output")
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	s := &State{}
	out := filepath.Join(t.TempDir(), "export.out")
	err := Export(s, out, ExportOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}
