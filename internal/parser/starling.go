package parser

import "io"

// StarlingParser parses Starling bank statement CSVs.
//
// Expected header:
//
//	Date,Counter Party,Reference,Type,Amount (GBP),Balance (GBP),Spending Category,Notes
//
// The Date column is exported by Starling in DD/MM/YYYY format and is
// normalised to YYYY-MM-DD in the output JSON.
type StarlingParser struct{}

// Parse implements Parser for Starling CSV files.
// Date values in DD/MM/YYYY format are converted to YYYY-MM-DD.
func (StarlingParser) Parse(r io.Reader) ([]string, error) {
	return parseCSVWithTransform(r, func(record map[string]string) {
		if v, ok := record["date"]; ok {
			record["date"] = reformatDateDDMMYYYY(v)
		}
	})
}
