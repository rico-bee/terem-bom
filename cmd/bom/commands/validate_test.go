package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestValidateCommandHelp(t *testing.T) {
	verbose := false
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Validate command help failed: %v", err)
	}

	output := buf.String()
	expectedStrings := []string{
		"Validate a Bureau of Meteorology (BOM) CSV file format",
		"Usage:",
		"Flags:",
		"--input",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestValidateCommandMissingInput(t *testing.T) {
	verbose := false
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error for missing input file")
	}

	output := buf.String()
	if !strings.Contains(output, "required") {
		t.Errorf("Expected error message about required input, got: %s", output)
	}
}

func TestValidateCommandInvalidFile(t *testing.T) {
	verbose := false
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{"--input", "nonexistent_file.csv"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestValidateCommandValidCSV(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,5.2,1,Y
IDCJAC0009,066062,2020,1,2,0.0,1,Y`

	tmpFile, err := os.CreateTemp("", "test_validate_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	verbose := false
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name()})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&errBuf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Validate command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "✓ CSV file is valid") {
		t.Errorf("Expected validation success message, got: %s", output)
	}
}

func TestValidateCommandInvalidCSV(t *testing.T) {
	csvContent := `Invalid,Header
invalid,data`

	tmpFile, err := os.CreateTemp("", "test_validate_invalid_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	verbose := false
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name()})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.Execute()
	if err == nil {
		t.Fatal("Expected error for invalid CSV")
	}
}

func TestValidateCommandVerbose(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,5.2,1,Y`

	tmpFile, err := os.CreateTemp("", "test_validate_verbose_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	verbose := true
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name()})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&errBuf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Validate command with verbose failed: %v", err)
	}

	errOutput := errBuf.String()
	if !strings.Contains(errOutput, "Validating CSV file:") {
		t.Errorf("Expected verbose output about validation, got: %s", errOutput)
	}

	if !strings.Contains(errOutput, "✓ CSV file is valid") {
		t.Errorf("Expected verbose success message, got: %s", errOutput)
	}
}

func TestValidateCommandFlags(t *testing.T) {
	verbose := false
	cmd := NewValidateCmd(&verbose)

	inputFlag := cmd.Flags().Lookup("input")
	if inputFlag == nil {
		t.Fatal("Input flag not found")
	}

	if inputFlag.Name != "input" {
		t.Errorf("Expected flag name 'input', got: %s", inputFlag.Name)
	}

	if inputFlag.Shorthand != "i" {
		t.Errorf("Expected flag shorthand 'i', got: %s", inputFlag.Shorthand)
	}
}

func TestValidateCommandEmptyCSV(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality`

	tmpFile, err := os.CreateTemp("", "test_validate_empty_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = tmpFile.WriteString(csvContent)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	verbose := false
	cmd := NewValidateCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name()})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&errBuf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Validate command with empty CSV failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "✓ CSV file is valid") {
		t.Errorf("Expected validation success message for empty CSV, got: %s", output)
	}
}
