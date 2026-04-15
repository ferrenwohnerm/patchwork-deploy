package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Record represents a single patch entry for diffing.
type Record struct {
	Patch       string
	Environment string
	AppliedAt   string
}

// Result holds the outcome of comparing two sets of patch records.
type Result struct {
	OnlyInA []Record
	OnlyInB []Record
	InBoth  []Record
}

// Compare returns patches present in a but not b, b but not a, and both.
func Compare(a, b []Record) Result {
	aMap := index(a)
	bMap := index(b)

	var result Result

	for key, rec := range aMap {
		if _, found := bMap[key]; found {
			result.InBoth = append(result.InBoth, rec)
		} else {
			result.OnlyInA = append(result.OnlyInA, rec)
		}
	}

	for key, rec := range bMap {
		if _, found := aMap[key]; !found {
			result.OnlyInB = append(result.OnlyInB, rec)
		}
	}

	sortRecords(result.OnlyInA)
	sortRecords(result.OnlyInB)
	sortRecords(result.InBoth)

	return result
}

// Format renders a diff result as a human-readable string.
func Format(r Result, labelA, labelB string) string {
	var sb strings.Builder

	for _, rec := range r.OnlyInA {
		fmt.Fprintf(&sb, "- [%s] %s\n", rec.Environment, rec.Patch)
	}
	for _, rec := range r.OnlyInB {
		fmt.Fprintf(&sb, "+ [%s] %s\n", rec.Environment, rec.Patch)
	}
	for _, rec := range r.InBoth {
		fmt.Fprintf(&sb, "  [%s] %s\n", rec.Environment, rec.Patch)
	}

	if sb.Len() == 0 {
		return fmt.Sprintf("no differences between %s and %s\n", labelA, labelB)
	}
	return sb.String()
}

func index(records []Record) map[string]Record {
	m := make(map[string]Record, len(records))
	for _, r := range records {
		key := r.Environment + "|" + r.Patch
		m[key] = r
	}
	return m
}

func sortRecords(recs []Record) {
	sort.Slice(recs, func(i, j int) bool {
		if recs[i].Environment != recs[j].Environment {
			return recs[i].Environment < recs[j].Environment
		}
		return recs[i].Patch < recs[j].Patch
	})
}
