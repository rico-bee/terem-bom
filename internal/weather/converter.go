package weather

import (
	"encoding/json"
	"fmt"
)

// Converter handles conversion of WeatherData to JSON
// Only responsible for conversion and error wrapping.
type Converter struct{}

// NewConverter creates a new Converter
func NewConverter() *Converter {
	return &Converter{}
}

// ToJSON serializes WeatherData to pretty-printed JSON.
// Returns error with context if conversion fails.
func (c *Converter) ToJSON(data WeatherData) ([]byte, error) {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		yearCount := len(data.WeatherDataForYear)
		return nil, fmt.Errorf("failed to convert weather data to JSON (years: %d): %w", yearCount, err)
	}
	return out, nil
}
