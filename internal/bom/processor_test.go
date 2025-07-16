package bom

import (
	"os"
	"strings"
	"testing"
)

func TestProcessorValidCSV(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,5.2,1,Y
IDCJAC0009,066062,2020,1,2,0.0,1,Y`

	tmpFile, err := os.CreateTemp("", "test_processor_*.csv")
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

	outputFile, err := os.CreateTemp("", "test_output_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	processor := NewProcessor()
	err = processor.ProcessWeatherDataFile(tmpFile.Name(), outputFile.Name())
	if err != nil {
		t.Fatalf("Processor failed: %v", err)
	}

	outputData, err := os.ReadFile(outputFile.Name())
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

func TestProcessorInvalidFile(t *testing.T) {
	processor := NewProcessor()
	err := processor.ProcessWeatherDataFile("nonexistent_file.csv", "output.json")
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestProcessorInvalidCSV(t *testing.T) {
	csvContent := `Invalid,Header
invalid,data`

	tmpFile, err := os.CreateTemp("", "test_processor_invalid_*.csv")
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

	outputFile, err := os.CreateTemp("", "test_output_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	processor := NewProcessor()
	err = processor.ProcessWeatherDataFile(tmpFile.Name(), outputFile.Name())
	if err == nil {
		t.Fatal("Expected error for invalid CSV")
	}
}

func TestProcessorEmptyCSV(t *testing.T) {
	csvContent := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality`

	tmpFile, err := os.CreateTemp("", "test_processor_empty_*.csv")
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

	outputFile, err := os.CreateTemp("", "test_output_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp output file: %v", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	processor := NewProcessor()
	err = processor.ProcessWeatherDataFile(tmpFile.Name(), outputFile.Name())
	if err != nil {
		t.Fatalf("Processor with empty CSV failed: %v", err)
	}

	outputData, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(outputData)
	if !strings.Contains(outputStr, "WeatherData") {
		t.Errorf("Expected JSON output to contain 'WeatherData', got: %s", outputStr)
	}
}
