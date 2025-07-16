package weather

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewConverter(t *testing.T) {
	conv := NewConverter()
	if conv == nil {
		t.Fatal("NewConverter returned nil")
	}
}

// Output is always pretty-printed (indented)
func TestToJSON_OutputIsPrettyPrinted(t *testing.T) {
	conv := NewConverter()
	data := WeatherData{
		WeatherDataForYear: []WeatherDataForYear{
			{Year: "2021", FirstRecordedDate: "2021-01-01", LastRecordedDate: "2021-12-31"},
		},
	}
	jsonBytes, err := conv.ToJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if !strings.Contains(string(jsonBytes), "\n  ") {
		t.Errorf("Expected pretty-printed output, got: %s", string(jsonBytes))
	}
	if !strings.Contains(string(jsonBytes), "2021-12-31") {
		t.Errorf("Expected output to contain date, got: %s", string(jsonBytes))
	}
}

func TestToJSON_Empty(t *testing.T) {
	conv := NewConverter()
	data := WeatherData{}
	jsonBytes, err := conv.ToJSON(data)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	var out WeatherData
	if err := json.Unmarshal(jsonBytes, &out); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
}

func TestToJSON_ErrorHandling(t *testing.T) {
	conv := NewConverter()
	data := WeatherData{
		WeatherDataForYear: []WeatherDataForYear{
			{Year: "2020", FirstRecordedDate: "2020-01-01", LastRecordedDate: "2020-12-31"},
		},
	}
	_, err := conv.ToJSON(data)
	if err != nil {
		t.Fatalf("Expected no error with valid data, got: %v", err)
	}
}
