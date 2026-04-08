package parser

import "io"

// MonzoParser parses Monzo bank statement CSVs.
//
// Expected header:
//
//	Transaction ID,Date,Time,Type,Name,Emoji,Category,Amount,Currency,Local amount,Local currency,Notes and #tags,Address,Receipt,Description,Category split,Money Out,Money In
type MonzoParser struct{}

// Parse implements Parser for Monzo CSV files.
func (MonzoParser) Parse(r io.Reader, filename string) ([]string, error) {
	return parseCSV(r, filename)
}

// MonzoFlexParser parses Monzo Flex bank statement CSVs.
//
// Monzo Flex exports share the same CSV structure as standard Monzo exports.
type MonzoFlexParser struct{}

// Parse implements Parser for Monzo Flex CSV files.
func (MonzoFlexParser) Parse(r io.Reader, filename string) ([]string, error) {
	return parseCSV(r, filename)
}
