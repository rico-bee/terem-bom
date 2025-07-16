package commands

import (
	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the CLI
func NewRootCmd() *cobra.Command {
	var verbose bool

	rootCmd := &cobra.Command{
		Use:   "bom",
		Short: "BOM Weather Data Processor",
		Long: `BOM Weather Data Processor

A command-line tool for processing Bureau of Meteorology (BOM) weather CSV files
into structured JSON with detailed yearly and monthly rainfall statistics.

The tool reads BOM CSV files and outputs aggregated weather data including:
- Total rainfall by year and month
- Average daily rainfall
- Days with/without rainfall
- Longest consecutive days of rainfall
- Median daily rainfall (monthly)`,
		Version: "1.0.0",
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(NewConvertCmd(&verbose))
	rootCmd.AddCommand(NewValidateCmd(&verbose))
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd
}

// Execute runs the root command
func Execute() error {
	return NewRootCmd().Execute()
}
