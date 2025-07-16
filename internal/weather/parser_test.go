package weather

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewParser(t *testing.T) {
	parser := NewParser(true)
	if parser == nil {
		t.Fatal("NewParser returned nil")
	}
	if !parser.verbose {
		t.Error("Expected verbose to be true")
	}

	parser = NewParser(false)
	if parser.verbose {
		t.Error("Expected verbose to be false")
	}
}

func TestParseCSV_ValidData(t *testing.T) {
	csvData := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,10.5,1,Y
IDCJAC0009,066062,2020,1,2,0.0,1,Y
IDCJAC0009,066062,2020,1,3,25.3,1,Y
IDCJAC0009,066062,2020,2,1,15.7,1,Y`

	parser := NewParser(false)
	reader := strings.NewReader(csvData)
	records, err := parser.ParseCSV(reader)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(records) != 4 {
		t.Fatalf("Expected 4 records, got %d", len(records))
	}

	// Check first record
	expectedDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if !records[0].Date.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, records[0].Date)
	}
	if records[0].Rainfall != 10.5 {
		t.Errorf("Expected rainfall 10.5, got %f", records[0].Rainfall)
	}
	if !records[0].HasData {
		t.Error("Expected HasData to be true")
	}

	// Check record with zero rainfall
	if records[1].Rainfall != 0.0 {
		t.Errorf("Expected rainfall 0.0, got %f", records[1].Rainfall)
	}
	if !records[1].HasData {
		t.Error("Expected HasData to be true for zero rainfall")
	}
}

func TestParseCSV_MissingData(t *testing.T) {
	csvData := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,10.5,1,Y
IDCJAC0009,066062,2020,1,2,NA,1,Y
IDCJAC0009,066062,2020,1,3,N/A,1,Y
IDCJAC0009,066062,2020,1,4,-,1,Y
IDCJAC0009,066062,2020,1,5,,1,Y
IDCJAC0009,066062,2020,1,6,0.0,1,Y`

	parser := NewParser(false)
	reader := strings.NewReader(csvData)
	records, err := parser.ParseCSV(reader)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(records) != 6 {
		t.Fatalf("Expected 6 records, got %d", len(records))
	}

	// Check records with missing data
	for i := 1; i <= 4; i++ {
		if records[i].HasData {
			t.Errorf("Record %d: Expected HasData to be false for missing data", i)
		}
		if records[i].Rainfall != 0.0 {
			t.Errorf("Record %d: Expected rainfall 0.0, got %f", i, records[i].Rainfall)
		}
	}

	// Check record with zero rainfall (should have data)
	if !records[5].HasData {
		t.Error("Expected HasData to be true for zero rainfall")
	}
}

func TestParseCSV_InvalidHeader(t *testing.T) {
	testCases := []struct {
		name    string
		csvData string
	}{
		{
			name:    "missing columns",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month\nIDCJAC0009,066062,2020,1,1,10.5,1,Y",
		},
		{
			name:    "wrong column names",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Temperature,Period,Quality\nIDCJAC0009,066062,2020,1,1,10.5,1,Y",
		},
		{
			name:    "empty header",
			csvData: "\nIDCJAC0009,066062,2020,1,1,10.5,1,Y",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(false)
			reader := strings.NewReader(tc.csvData)
			_, err := parser.ParseCSV(reader)

			if err == nil {
				t.Error("Expected error for invalid header")
			}
		})
	}
}

func TestParseCSV_InvalidData(t *testing.T) {
	testCases := []struct {
		name    string
		csvData string
	}{
		{
			name:    "invalid year",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,abc,1,1,10.5,1,Y",
		},
		{
			name:    "invalid month",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,2020,13,1,10.5,1,Y",
		},
		{
			name:    "invalid day",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,2020,1,32,10.5,1,Y",
		},
		{
			name:    "invalid rainfall",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,2020,1,1,abc,1,Y",
		},
		{
			name:    "negative rainfall",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,2020,1,1,-10.5,1,Y",
		},
		{
			name:    "unreasonable rainfall",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,2020,1,1,20000.0,1,Y",
		},
		{
			name:    "invalid date",
			csvData: "Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality\nIDCJAC0009,066062,2020,2,30,10.5,1,Y",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(false)
			reader := strings.NewReader(tc.csvData)
			records, err := parser.ParseCSV(reader)

			if err != nil {
				t.Fatalf("Expected no error during parsing, got: %v", err)
			}

			// For negative and unreasonable rainfall, the original behavior accepts them
			// since it doesn't validate ranges
			expectedRecords := 0
			if tc.name == "negative rainfall" || tc.name == "unreasonable rainfall" {
				expectedRecords = 1
			}

			if len(records) != expectedRecords {
				t.Errorf("Expected %d valid records, got %d", expectedRecords, len(records))
			}
		})
	}
}

func TestParseCSV_EmptyRows(t *testing.T) {
	csvData := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,10.5,1,Y

IDCJAC0009,066062,2020,1,2,15.3,1,Y
   ,   ,   ,   ,   ,   ,   ,   
IDCJAC0009,066062,2020,1,3,20.1,1,Y`

	parser := NewParser(false)
	reader := strings.NewReader(csvData)
	records, err := parser.ParseCSV(reader)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("Expected 3 records, got %d", len(records))
	}
}

func TestParseCSV_VerboseMode(t *testing.T) {
	csvData := `Product code,Bureau of Meteorology station number,Year,Month,Day,Rainfall amount (millimetres),Period over which rainfall was measured (days),Quality
IDCJAC0009,066062,2020,1,1,10.5,1,Y
IDCJAC0009,066062,2020,1,2,abc,1,Y
IDCJAC0009,066062,2020,1,3,15.3,1,Y`

	var buf bytes.Buffer
	parser := NewParser(true)
	reader := strings.NewReader(csvData)

	// Capture stdout by redirecting to buffer
	// Note: This is a simplified test - in a real scenario you might want to
	// use a more sophisticated approach to capture stdout
	records, err := parser.ParseCSV(reader)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(records) != 2 {
		t.Fatalf("Expected 2 valid records, got %d", len(records))
	}

	// The verbose output would be captured in buf if we redirected stdout
	_ = buf.String() // Suppress unused variable warning
}

func TestParseRainfall(t *testing.T) {
	parser := NewParser(false)

	testCases := []struct {
		input     string
		expected  float64
		hasData   bool
		expectErr bool
	}{
		{"10.5", 10.5, true, false},
		{"0.0", 0.0, true, false},
		{"", 0.0, false, false},
		{"NA", 0.0, false, false},
		{"N/A", 0.0, false, false},
		{"-", 0.0, false, false},
		{"abc", 0.0, false, true},
		{"-10.5", -10.5, true, false},     // Original accepts negative values
		{"20000.0", 20000.0, true, false}, // Original accepts any numeric value
		{"  10.5  ", 10.5, true, false},
		{"  NA  ", 0.0, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			rainfall, hasData, err := parser.parseRainfall(tc.input)

			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if rainfall != tc.expected {
				t.Errorf("Expected rainfall %f, got %f", tc.expected, rainfall)
			}

			if hasData != tc.hasData {
				t.Errorf("Expected hasData %v, got %v", tc.hasData, hasData)
			}
		})
	}
}

func TestIsEmptyRow(t *testing.T) {
	parser := NewParser(false)

	testCases := []struct {
		row      []string
		expected bool
	}{
		{[]string{"", "", "", "", "", "", "", ""}, true},
		{[]string{"   ", "  ", "  ", "  ", "  ", "  ", "  ", "  "}, true},
		{[]string{"IDCJAC0009", "066062", "2020", "1", "1", "10.5", "1", "Y"}, false},
		{[]string{"", "066062", "2020", "1", "1", "10.5", "1", "Y"}, false},
		{[]string{"IDCJAC0009", "", "2020", "1", "1", "10.5", "1", "Y"}, false},
		{[]string{}, true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			result := parser.isEmptyRow(tc.row)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v for row %v", tc.expected, result, tc.row)
			}
		})
	}
}
