// Package parser provides CSV parsers for supported bank statement formats.
package parser

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
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
	// filename is the base name of the source file and is included in each record.
	Parse(r io.Reader, filename string) ([]string, error)
}


// reformatDateDDMMYYYY converts a date string in DD/MM/YYYY format to
// YYYY-MM-DD. If the input is not in DD/MM/YYYY format the original value is
// returned unchanged.
func reformatDateDDMMYYYY(s string) string {
	t, err := time.Parse("02/01/2006", s)
	if err != nil {
		return s
	}
	return t.Format("2006-01-02")
}

// parseCSV is a generic CSV parser that converts every data row (after the
// header) into a JSON string using the column headers as keys.
func parseCSV(r io.Reader, filename string) ([]string, error) {
	return parseCSVWithTransform(r, filename, nil)
}

// parseCSVWithTransform is like parseCSV but applies an optional transform
// function to each record map before marshalling it to JSON.  A nil transform
// is a no-op.
func parseCSVWithTransform(r io.Reader, filename string, transform func(map[string]string)) ([]string, error) {
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
		j, err := rowToJSONWithTransform(headers, row, filename, transform)
		if err != nil {
			return nil, err
		}
		results = append(results, j)
	}
	return results, nil
}

// rowToJSONWithTransform maps CSV headers to row values, applies an optional
// transform to the resulting map, then marshals it to JSON. filename is added
// to the record as the "filename" field.
func rowToJSONWithTransform(headers, row []string, filename string, transform func(map[string]string)) (string, error) {
	if len(headers) != len(row) {
		return "", fmt.Errorf("header count %d does not match row count %d", len(headers), len(row))
	}
	record := make(map[string]string, len(headers))
	for i, h := range headers {
		record[SanitiseKey(h)] = strings.TrimSpace(row[i])
	}
	if transform != nil {
		transform(record)
	}
	record["filename"] = filename
	b, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("marshalling record: %w", err)
	}
	return string(b), nil
}
