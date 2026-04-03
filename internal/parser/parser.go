// Package parser provides CSV parsers for supported bank statement formats.
package parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var nonAlpha = regexp.MustCompile(`[^a-z]+`)

// SanitiseKey normalises a CSV column heading or JSON key to lowercase with
// all non-alphabetic characters (including spaces) replaced by underscores.
// Consecutive non-alphabetic characters are collapsed into a single underscore,
// and leading/trailing underscores are removed.
func SanitiseKey(s string) string {
	return strings.Trim(nonAlpha.ReplaceAllString(strings.ToLower(s), "_"), "_")
}

// Parser converts CSV rows into JSON-encoded transaction records.
type Parser interface {
	// Parse reads all records from r and returns a slice of JSON-encoded strings.
	Parse(r io.Reader) ([]string, error)
}

// rowToJSON maps CSV headers to row values and marshals the result to JSON.
func rowToJSON(headers, row []string) (string, error) {
	if len(headers) != len(row) {
		return "", fmt.Errorf("header count %d does not match row count %d", len(headers), len(row))
	}
	record := make(map[string]string, len(headers))
	for i, h := range headers {
		record[SanitiseKey(h)] = strings.TrimSpace(row[i])
	}
	b, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("marshalling record: %w", err)
	}
	return string(b), nil
}

// parseCSV is a generic CSV parser that converts every data row (after the
// header) into a JSON string using the column headers as keys.
func parseCSV(r io.Reader) ([]string, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true

	headers, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("reading CSV header: %w", err)
	}

	var results []string
	for {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV row: %w", err)
		}
		j, err := rowToJSON(headers, row)
		if err != nil {
			return nil, err
		}
		results = append(results, j)
	}
	return results, nil
}
