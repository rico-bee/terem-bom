package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/terem/bom/internal/bom"
)

func NewConvertCmd(verbose *bool) *cobra.Command {
	var inputFile string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "Convert BOM CSV file to JSON",
		Long: `Convert a Bureau of Meteorology (BOM) CSV file to structured JSON format.

The convert command reads a BOM weather CSV file and outputs aggregated weather data
in JSON format with detailed yearly and monthly rainfall statistics.

Example:
  bom convert -i weather.csv -o output.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			processor := bom.NewProcessorWithVerbose(*verbose)

			// Open input file
			inFile, err := os.Open(inputFile)
			if err != nil {
				return fmt.Errorf("failed to open input file %s: %w", inputFile, err)
			}
			defer inFile.Close()

			// Determine output destination
			var output io.Writer
			var outputName string
			if outputFile == "" {
				output = cmd.OutOrStdout()
				outputName = "stdout"
			} else {
				outFile, err := os.Create(outputFile)
				if err != nil {
					return fmt.Errorf("failed to create output file %s: %w", outputFile, err)
				}
				defer outFile.Close()
				output = outFile
				outputName = outputFile
			}

			// Process the data
			if err := processor.ProcessWeatherData(inFile, output); err != nil {
				return fmt.Errorf("conversion failed: %w", err)
			}

			if *verbose {
				fmt.Fprintf(cmd.ErrOrStderr(), "Successfully converted %s to %s\n", inputFile, outputName)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input CSV file path (required)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output JSON file path (required)")
	cmd.MarkFlagRequired("input")

	return cmd
}
