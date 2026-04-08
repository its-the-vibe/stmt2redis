package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/its-the-vibe/stmt2redis/internal/config"
	"github.com/its-the-vibe/stmt2redis/internal/parser"
	"github.com/its-the-vibe/stmt2redis/internal/redisclient"
	"github.com/spf13/cobra"
)

var (
	csvType    string
	csvFile    string
	stdoutOnly bool
	envFile    string
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Parse a CSV file and push transactions to Redis (or stdout)",
	Long: `push reads a bank statement CSV, converts each transaction row to JSON,
and either RPUSHes it to the configured Redis list or prints it to stdout.`,
	RunE: runPush,
}

func init() {
	pushCmd.Flags().StringVarP(&csvType, "type", "t", "", "CSV type: starling, amex, monzo, monzo-flex (required)")
	pushCmd.Flags().StringVarP(&csvFile, "file", "f", "", "path to the CSV file (required)")
	pushCmd.Flags().BoolVar(&stdoutOnly, "stdout", false, "print JSON to stdout instead of publishing to Redis")
	pushCmd.Flags().StringVar(&envFile, "env-file", ".env", ".env file path")

	_ = pushCmd.MarkFlagRequired("type")
	_ = pushCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(pushCmd)
}

// newParser returns the Parser implementation for the given csvType.
func newParser(csvType string) (parser.Parser, error) {
	switch csvType {
	case "starling":
		return parser.StarlingParser{}, nil
	case "amex":
		return parser.AmexParser{}, nil
	case "monzo":
		return parser.MonzoParser{}, nil
	case "monzo-flex":
		return parser.MonzoFlexParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported CSV type %q: must be one of starling, amex, monzo, monzo-flex", csvType)
	}
}

func runPush(cmd *cobra.Command, args []string) error {
	// Load .env for secrets (non-fatal if the file doesn't exist).
	if err := godotenv.Load(envFile); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "warning: could not load env file %q: %v\n", envFile, err)
	}

	// Load config (only needed when publishing to Redis).
	var cfg *config.Config
	var listKey string

	if !stdoutOnly {
		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		listKey, err = cfg.ListKey(csvType)
		if err != nil {
			return err
		}
		if listKey == "" {
			return fmt.Errorf("no Redis list configured for type %q in config file", csvType)
		}
	}

	// Open the CSV file.
	f, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("opening CSV file %q: %w", csvFile, err)
	}
	defer f.Close()

	// Select parser.
	p, err := newParser(csvType)
	if err != nil {
		return err
	}

	// Parse CSV into JSON records.
	records, err := p.Parse(f, filepath.Base(csvFile))
	if err != nil {
		return fmt.Errorf("parsing CSV: %w", err)
	}

	if len(records) == 0 {
		fmt.Fprintln(os.Stderr, "no records found in CSV file")
		return nil
	}

	// Output to stdout or Redis.
	if stdoutOnly {
		for _, rec := range records {
			fmt.Println(rec)
		}
		return nil
	}

	client, err := redisclient.New(cfg)
	if err != nil {
		return fmt.Errorf("creating Redis client: %w", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.RPush(ctx, listKey, records...); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "pushed %d records to Redis list %q\n", len(records), listKey)
	return nil
}
