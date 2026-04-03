// Package cmd contains the Cobra CLI commands for stmt2redis.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "stmt2redis",
	Short: "Parse bank statement CSVs and stream transactions into Redis",
	Long: `stmt2redis is a CLI tool that reads bank statement CSV files
and pushes each transaction as a JSON object into a Redis list.

Supported CSV types: starling, amex, monzo, monzo-flex`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.yaml", "config file path")
}
