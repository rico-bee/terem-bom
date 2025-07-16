package bom

import (
	"fmt"
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

// ProcessWeatherDataFile reads a CSV file and outputs to JSON file
func (p *Processor) ProcessWeatherDataFile(inputPath, outputPath string) error {
	// Open the CSV file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file %s: %w", inputPath, err)
	}
	defer file.Close()

	// Parse CSV file
	records, err := p.parser.ParseCSV(file)
	if err != nil {
		return fmt.Errorf("failed to parse CSV file %s: %w", inputPath, err)
	}

	// Aggregate the records
	weatherData := p.aggregator.Aggregate(records)

	// Convert to JSON
	jsonData, err := p.converter.ToJSON(weatherData)
	if err != nil {
		return fmt.Errorf("failed to convert weather data to JSON: %w", err)
	}

	// Write to JSON file
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write output to %s: %w", outputPath, err)
	}

	return nil
}
