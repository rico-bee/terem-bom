package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/terem/bom/internal/bom"
)

func NewValidateCmd(verbose *bool) *cobra.Command {
	var inputFile string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate BOM CSV file format",
		Long: `Validate a Bureau of Meteorology (BOM) CSV file format.

The validate command checks if a CSV file has the correct BOM format without
performing any conversion. It validates the header structure and data format.

Example:
  bom validate -i weather.csv`,
		RunE: func(cmd *cobra.Command, args []string) error {
			processor := bom.NewProcessorWithVerbose(*verbose)
			if *verbose {
				fmt.Fprintf(cmd.ErrOrStderr(), "Validating CSV file: %s\n", inputFile)
			}

			// Validate the CSV file
			err := processor.ValidateCSVFile(inputFile)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			if *verbose {
				fmt.Fprintf(cmd.ErrOrStderr(), "✓ CSV file is valid\n")
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "✓ CSV file is valid\n")
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input CSV file path (required)")
	cmd.MarkFlagRequired("input")

	return cmd
}
