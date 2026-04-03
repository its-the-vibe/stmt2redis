package parser

import "io"

// StarlingParser parses Starling bank statement CSVs.
//
// Expected header:
//
//	Date,Counter Party,Reference,Type,Amount (GBP),Balance (GBP),Spending Category,Notes
type StarlingParser struct{}

// Parse implements Parser for Starling CSV files.
func (StarlingParser) Parse(r io.Reader) ([]string, error) {
	return parseCSV(r)
}
