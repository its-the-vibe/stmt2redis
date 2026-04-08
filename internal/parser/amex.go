package parser

import "io"

// AmexParser parses American Express bank statement CSVs.
//
// Expected header:
//
//	Date,Description,Amount,Extended Details,Appears On Your Statement As,Address,Town/City,Postcode,Country,Reference,Category
type AmexParser struct{}

// Parse implements Parser for Amex CSV files.
func (AmexParser) Parse(r io.Reader, filename string) ([]string, error) {
	return parseCSV(r, filename)
}
