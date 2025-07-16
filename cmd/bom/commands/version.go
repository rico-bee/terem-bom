package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long: `Display version information for the BOM Weather Data Processor.

Shows the current version, build information, and other relevant details.`,
		Run: func(cmd *cobra.Command, args []string) {
			version := "1.0.0"
			if cmd.Root() != nil && cmd.Root().Version != "" {
				version = cmd.Root().Version
			}
			fmt.Fprintf(cmd.OutOrStdout(), "BOM Weather Data Processor v%s\n", version)
			fmt.Fprintf(cmd.OutOrStdout(), "A tool for processing Bureau of Meteorology weather CSV files\n")
			fmt.Fprintf(cmd.OutOrStdout(), "into structured JSON with detailed rainfall statistics.\n")
		},
	}
	return cmd
}
