package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	rootCmd := &cobra.Command{Use: "bom", Version: "1.0.0"}
	cmd := NewVersionCmd()
	rootCmd.AddCommand(cmd)
	rootCmd.SetArgs([]string{"version"})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&errBuf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}

	output := buf.String() + errBuf.String()
	expectedStrings := []string{
		"BOM Weather Data Processor v1.0.0",
		"A tool for processing Bureau of Meteorology weather CSV files",
		"into structured JSON with detailed rainfall statistics.",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestVersionCommandHelp(t *testing.T) {
	rootCmd := &cobra.Command{Use: "bom"}
	cmd := NewVersionCmd()
	rootCmd.AddCommand(cmd)
	rootCmd.SetArgs([]string{"version", "--help"})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&errBuf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Version command help failed: %v", err)
	}

	output := buf.String() + errBuf.String()
	expectedStrings := []string{
		"Display version information",
		"Usage:",
		"version",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestVersionCommandWithArgs(t *testing.T) {
	rootCmd := &cobra.Command{Use: "bom", Version: "1.0.0"}
	cmd := NewVersionCmd()
	rootCmd.AddCommand(cmd)
	rootCmd.SetArgs([]string{"version", "extra", "arguments"})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&errBuf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Version command with extra args failed: %v", err)
	}

	output := buf.String() + errBuf.String()
	if !strings.Contains(output, "BOM Weather Data Processor v1.0.0") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

func TestVersionCommandDescription(t *testing.T) {
	cmd := NewVersionCmd()
	if cmd.Short == "" {
		t.Error("Version command should have a short description")
	}

	if cmd.Long == "" {
		t.Error("Version command should have a long description")
	}

	if !strings.Contains(cmd.Short, "version") {
		t.Errorf("Version command short description should mention 'version', got: %s", cmd.Short)
	}
}
