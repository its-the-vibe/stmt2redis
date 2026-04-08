package parser_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/its-the-vibe/stmt2redis/internal/parser"
)

func TestStarlingParser(t *testing.T) {
	csv := `Date,Counter Party,Reference,Type,Amount (GBP),Balance (GBP),Spending Category,Notes
2024-01-15,TESCO STORES,REF12345,FASTER PAYMENT,-12.50,1500.00,GROCERIES,weekly shop
2024-01-16,AMAZON,REF67890,CARD PAYMENT,-29.99,1470.01,SHOPPING,
`
	p := parser.StarlingParser{}
	records, err := p.Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	for _, rec := range records {
		if !strings.Contains(rec, `"date"`) {
			t.Errorf("expected record to contain date key, got: %s", rec)
		}
		if !strings.Contains(rec, `"amount_gbp"`) {
			t.Errorf("expected record to contain amount_gbp key, got: %s", rec)
		}
	}
}

func TestStarlingParserDateReformat(t *testing.T) {
	// Starling exports dates as DD/MM/YYYY; the parser must convert them to
	// YYYY-MM-DD in the output JSON.
	csv := `Date,Counter Party,Reference,Type,Amount (GBP),Balance (GBP),Spending Category,Notes
21/03/2026,TESCO STORES,REF12345,FASTER PAYMENT,-12.50,1500.00,GROCERIES,weekly shop
05/01/2025,AMAZON,REF67890,CARD PAYMENT,-29.99,1470.01,SHOPPING,
`
	p := parser.StarlingParser{}
	records, err := p.Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	wantDates := []string{"2026-03-21", "2025-01-05"}
	for i, rec := range records {
		want := fmt.Sprintf(`"date":"%s"`, wantDates[i])
		if !strings.Contains(rec, want) {
			t.Errorf("record %d: expected %s in %s", i, want, rec)
		}
	}
}

func TestAmexParser(t *testing.T) {
	csv := `Date,Description,Amount,Extended Details,Appears On Your Statement As,Address,Town/City,Postcode,Country,Reference,Category
01/01/2024,COFFEE SHOP,4.50,Coffee,COFFEE SHOP,123 High St,London,EC1A 1BB,UK,REF001,Dining
02/01/2024,GYM MEMBERSHIP,40.00,Monthly,GYM,456 Fitness Rd,London,SW1A 2AA,UK,REF002,Health
`
	p := parser.AmexParser{}
	records, err := p.Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	for _, rec := range records {
		if !strings.Contains(rec, `"description"`) {
			t.Errorf("expected record to contain description key, got: %s", rec)
		}
		if !strings.Contains(rec, `"category"`) {
			t.Errorf("expected record to contain category key, got: %s", rec)
		}
	}
}

func TestMonzoParser(t *testing.T) {
	csv := `Transaction ID,Date,Time,Type,Name,Emoji,Category,Amount,Currency,Local amount,Local currency,Notes and #tags,Address,Receipt,Description,Category split,Money Out,Money In
tx_001,2024-01-15,10:30:00,card_payment,Tesco,,Groceries,-1250,GBP,-1250,GBP,weekly shop,,,,,-12.50,
tx_002,2024-01-15,14:00:00,faster_payment,Salary,,Income,200000,GBP,200000,GBP,,,,,,,2000.00
`
	p := parser.MonzoParser{}
	records, err := p.Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	for _, rec := range records {
		if !strings.Contains(rec, `"transaction_id"`) {
			t.Errorf("expected record to contain transaction_id key, got: %s", rec)
		}
	}
}

func TestMonzoFlexParser(t *testing.T) {
	csv := `Transaction ID,Date,Time,Type,Name,Emoji,Category,Amount,Currency,Local amount,Local currency,Notes and #tags,Address,Receipt,Description,Category split,Money Out,Money In
flex_001,2024-02-01,09:00:00,flex,Apple Store,,Shopping,-99900,GBP,-99900,GBP,,,,,,-999.00,
`
	p := parser.MonzoFlexParser{}
	records, err := p.Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if !strings.Contains(records[0], `"transaction_id"`) {
		t.Errorf("expected record to contain transaction_id key, got: %s", records[0])
	}
}

func TestParserEmptyCSV(t *testing.T) {
	// CSV with only a header row produces no records.
	csv := "Date,Counter Party,Reference,Type,Amount (GBP),Balance (GBP),Spending Category,Notes\n"
	p := parser.StarlingParser{}
	records, err := p.Parse(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected 0 records, got %d", len(records))
	}
}

func TestSanitiseKey(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Date", "date"},
		{"Counter Party", "counter_party"},
		{"Amount (GBP)", "amount_gbp"},
		{"Balance (GBP)", "balance_gbp"},
		{"Transaction ID", "transaction_id"},
		{"Notes and #tags", "notes_and_tags"},
		{"Town/City", "town_city"},
		{"Spending Category", "spending_category"},
		// multiple consecutive non-alpha chars collapse to a single underscore
		{"foo   bar", "foo_bar"},
		{"foo!!!bar", "foo_bar"},
		// leading and trailing non-alpha chars are stripped
		{"(amount)", "amount"},
		{"  leading", "leading"},
		{"trailing  ", "trailing"},
		// already sanitised input is unchanged
		{"already_clean", "already_clean"},
		// purely non-alpha input
		{"123", ""},
	}
	for _, tc := range cases {
		got := parser.SanitiseKey(tc.input)
		if got != tc.want {
			t.Errorf("SanitiseKey(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
