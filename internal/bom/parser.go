package bom

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Parser handles CSV parsing for BOM weather data
type Parser struct {
	verbose bool
}

// NewParser creates a new parser instance
func NewParser(verbose bool) *Parser {
	return &Parser{
		verbose: verbose,
	}
}

// ParseCSV reads and parses a BOM weather data CSV file
func (p *Parser) ParseCSV(reader io.Reader) ([]DailyRecord, error) {
	csvReader := csv.NewReader(reader)

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	if p.verbose {
		fmt.Printf("CSV Header: %v\n", header)
	}

	// Validate header structure
	if err := p.validateHeader(header); err != nil {
		return nil, fmt.Errorf("invalid CSV header: %w", err)
	}

	var records []DailyRecord

	// Read data rows
	lineNum := 1 // Start counting from line 1 (after header)
	for {
		lineNum++ // Increment to get the current data row line number
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading row %d: %w", lineNum, err)
		}

		// Skip empty rows
		if p.isEmptyRow(row) {
			continue
		}

		record, err := p.parseRow(row)
		if err != nil {
			// Silently skip invalid rows (matching original behavior)
			// Only log in verbose mode for debugging
			if p.verbose {
				fmt.Printf("Warning: skipping row %d: %v\n", lineNum, err)
			}
			continue
		}

		records = append(records, record)
	}

	if p.verbose {
		fmt.Printf("Parsed %d valid records\n", len(records))
	}

	return records, nil
}

// validateHeader checks if the CSV header has the expected structure
func (p *Parser) validateHeader(header []string) error {
	if len(header) < 8 {
		return fmt.Errorf("expected at least 8 columns, got %d", len(header))
	}

	expectedColumns := []string{
		"Product code",
		"Bureau of Meteorology station number",
		"Year",
		"Month",
		"Day",
		"Rainfall amount (millimetres)",
		"Period over which rainfall was measured (days)",
		"Quality",
	}

	for i, expected := range expectedColumns {
		if i >= len(header) {
			return fmt.Errorf("missing column: %s", expected)
		}
		if !strings.EqualFold(strings.TrimSpace(header[i]), expected) {
			return fmt.Errorf("expected column %d to be '%s', got '%s'", i+1, expected, header[i])
		}
	}

	return nil
}

// isEmptyRow checks if a row is empty or contains only whitespace
func (p *Parser) isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// parseRow parses a single CSV row into a DailyRecord
func (p *Parser) parseRow(row []string) (DailyRecord, error) {
	if len(row) < 8 {
		return DailyRecord{}, fmt.Errorf("insufficient columns")
	}

	// Parse year (column 2, index 2)
	yearStr := strings.TrimSpace(row[2])
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return DailyRecord{}, fmt.Errorf("invalid year '%s': %w", yearStr, err)
	}

	// Parse month (column 3, index 3)
	monthStr := strings.TrimSpace(row[3])
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return DailyRecord{}, fmt.Errorf("invalid month '%s': %w", monthStr, err)
	}
	if month < 1 || month > 12 {
		return DailyRecord{}, fmt.Errorf("month out of range: %d", month)
	}

	// Parse day (column 4, index 4)
	dayStr := strings.TrimSpace(row[4])
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return DailyRecord{}, fmt.Errorf("invalid day '%s': %w", dayStr, err)
	}
	if day < 1 || day > 31 {
		return DailyRecord{}, fmt.Errorf("day out of range: %d", day)
	}

	// Parse rainfall (column 5, index 5)
	rainfallStr := strings.TrimSpace(row[5])
	rainfall, hasData, err := p.parseRainfall(rainfallStr)
	if err != nil {
		return DailyRecord{}, fmt.Errorf("invalid rainfall '%s': %w", rainfallStr, err)
	}

	// Create date
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate date (check for invalid dates like February 30th)
	if date.Year() != year || date.Month() != time.Month(month) || date.Day() != day {
		return DailyRecord{}, fmt.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return DailyRecord{
		Date:     date,
		Rainfall: rainfall,
		HasData:  hasData,
	}, nil
}

// parseRainfall parses rainfall value, handling missing data indicators
func (p *Parser) parseRainfall(value string) (float64, bool, error) {
	value = strings.TrimSpace(value)

	// Handle missing data indicators (matching original logic)
	if value == "" || value == "NA" || value == "N/A" || value == "-" {
		return 0.0, false, nil
	}

	// Parse numeric value
	rainfall, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, false, fmt.Errorf("cannot parse as float: %w", err)
	}

	// Remove range validation to match original behavior
	// Original silently accepted any numeric value

	return rainfall, true, nil
}
