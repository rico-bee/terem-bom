package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestConvertCommandHelp(t *testing.T) {
	verbose := false
	cmd := NewConvertCmd(&verbose)
	cmd.SetArgs([]string{"--help"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Convert command help failed: %v", err)
	}

	output := buf.String()
	expectedStrings := []string{
		"Convert a Bureau of Meteorology (BOM) CSV file to structured JSON format",
		"Usage:",
		"Flags:",
		"--input",
		"--output",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected output to contain '%s', got: %s", expected, output)
		}
	}
}

func TestConvertCommandMissingInput(t *testing.T) {
	verbose := false
	cmd := NewConvertCmd(&verbose)
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

func TestConvertCommandInvalidFile(t *testing.T) {
	verbose := false
	cmd := NewConvertCmd(&verbose)
	cmd.SetArgs([]string{"--input", "nonexistent_file.csv"})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestConvertCommandValidCSV(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,5.2,1,Y
IDCJAC0009,066062,2020,1,2,0.0,1,Y
IDCJAC0009,066062,2020,1,3,12.5,1,Y`

	tmpFile, err := os.CreateTemp("", "test_convert_*.csv")
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

	tmpOutput, err := os.CreateTemp("", "test_output_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	verbose := false
	cmd := NewConvertCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name(), "--output", tmpOutput.Name()})

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Convert command failed: %v", err)
	}

	outputData, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputData)
	if !strings.Contains(outputStr, "WeatherData") {
		t.Errorf("Expected JSON output to contain 'WeatherData', got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "2020") {
		t.Errorf("Expected JSON output to contain year '2020', got: %s", outputStr)
	}
}

func TestConvertCommandToStdout(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,5.2,1,Y
IDCJAC0009,066062,2020,1,2,0.0,1,Y`

	tmpFile, err := os.CreateTemp("", "test_convert_stdout_*.csv")
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
	cmd := NewConvertCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name()})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&errBuf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Convert command to stdout failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "WeatherData") {
		t.Errorf("Expected stdout output to contain 'WeatherData', got: %s", output)
	}
}

func TestConvertCommandVerbose(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,5.2,1,Y`

	tmpFile, err := os.CreateTemp("", "test_convert_verbose_*.csv")
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
	cmd := NewConvertCmd(&verbose)
	cmd.SetArgs([]string{"--input", tmpFile.Name()})

	var buf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&errBuf)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("Convert command with verbose failed: %v", err)
	}

	errOutput := errBuf.String()
	if !strings.Contains(errOutput, "Successfully converted") {
		t.Errorf("Expected verbose output about successful conversion, got: %s", errOutput)
	}
}

func TestConvertCommandFlags(t *testing.T) {
	verbose := false
	cmd := NewConvertCmd(&verbose)

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

	outputFlag := cmd.Flags().Lookup("output")
	if outputFlag == nil {
		t.Fatal("Output flag not found")
	}

	if outputFlag.Name != "output" {
		t.Errorf("Expected flag name 'output', got: %s", outputFlag.Name)
	}

	if outputFlag.Shorthand != "o" {
		t.Errorf("Expected flag shorthand 'o', got: %s", outputFlag.Shorthand)
	}
}
