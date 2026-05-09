package parser

import "io"

// SantanderParser parses Santander bank statement TSV files.
//
// Expected header:
//
//	Date|Description|Money In|Money Out|Balance
type SantanderParser struct{}

// Parse implements Parser for Santander TSV files.
func (SantanderParser) Parse(r io.Reader, filename string) ([]string, error) {
	return parseDelimitedWithTransform(r, filename, '|', nil)
}
