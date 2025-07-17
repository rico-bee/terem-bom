package bom

import (
	"fmt"
	"io"
	"os"
)

// Processor orchestrates the parsing, aggregation, and conversion of weather data
type Processor struct {
	parser     *Parser
	aggregator *Aggregator
	converter  *Converter
}

// NewProcessor creates a new Processor with all required components
func NewProcessor() *Processor {
	return &Processor{
		parser:     NewParser(false), // Default to non-verbose
		aggregator: NewAggregator(),
		converter:  NewConverter(),
	}
}

// NewProcessorWithVerbose creates a new Processor with verbose parsing
func NewProcessorWithVerbose(verbose bool) *Processor {
	return &Processor{
		parser:     NewParser(verbose),
		aggregator: NewAggregator(),
		converter:  NewConverter(),
	}
}

// ProcessWeatherData processes weather data from a CSV reader and writes JSON to the writer
func (p *Processor) ProcessWeatherData(input io.Reader, output io.Writer) error {
	// Parse CSV file
	records, err := p.parser.ParseCSV(input)
	if err != nil {
		return fmt.Errorf("failed to parse CSV: %w", err)
	}

	// Aggregate the records
	weatherData := p.aggregator.Aggregate(records)

	// Convert to JSON
	jsonData, err := p.converter.ToJSON(weatherData)
	if err != nil {
		return fmt.Errorf("failed to convert weather data to JSON: %w", err)
	}

	// Write to output
	_, err = output.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

// ValidateCSVFile validates a CSV file by attempting to parse it
// Returns error if the file is invalid, nil if valid
func (p *Processor) ValidateCSVFile(inputPath string) error {
	// Open the CSV file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file %s: %w", inputPath, err)
	}
	defer file.Close()

	// Parse CSV file - if this succeeds, the file is valid
	_, err = p.parser.ParseCSV(file)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}

// ProcessWeatherDataFile reads a CSV file and outputs to JSON file
func (p *Processor) ProcessWeatherDataFile(inputPath, outputPath string) error {
	// Open the CSV file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file %s: %w", inputPath, err)
	}
	defer file.Close()

	// Open output file for writing
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
	}
	defer outFile.Close()

	return p.ProcessWeatherData(file, outFile)
}
