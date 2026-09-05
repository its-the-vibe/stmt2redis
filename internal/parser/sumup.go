package parser

import "io"

// SumUpParser parses SumUp bank statement CSVs.
//
// Expected header:
//
//	Date|Reference|Type|Amount|Description
type SumUpParser struct{}

// Parse implements Parser for SumUp pipe-delimited CSV files.
func (SumUpParser) Parse(r io.Reader, filename string) ([]string, error) {
	return parseDelimitedWithTransform(r, filename, '|', nil)
}
