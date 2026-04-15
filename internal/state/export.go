package state

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// ExportFormat defines the output format for state export.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatText ExportFormat = "text"
)

// ExportOptions controls what gets exported.
type ExportOptions struct {
	Environment string
	Format      ExportFormat
}

// Export writes state records to the given file path in the requested format.
// If opts.Environment is non-empty, only records for that environment are written.
func Export(s *State, path string, opts ExportOptions) error {
	records := s.Records
	if opts.Environment != "" {
		records = s.ForEnvironment(opts.Environment)
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Environment != records[j].Environment {
			return records[i].Environment < records[j].Environment
		}
		return records[i].Patch < records[j].Patch
	})

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("export: create file: %w", err)
	}
	defer f.Close()

	switch opts.Format {
	case FormatJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(records); err != nil {
			return fmt.Errorf("export: encode json: %w", err)
		}
	case FormatText, "":
		for _, r := range records {
			line := fmt.Sprintf("%s\t%s\t%s\n", r.Environment, r.Patch, r.AppliedAt.Format("2006-01-02T15:04:05Z"))
			if _, err := f.WriteString(line); err != nil {
				return fmt.Errorf("export: write text: %w", err)
			}
		}
	default:
		return fmt.Errorf("export: unknown format %q", opts.Format)
	}

	return nil
}
