package bom

import (
	"time"
)

// WeatherData represents the root structure of the JSON output
type WeatherData struct {
	WeatherDataForYear []WeatherDataForYear `json:"WeatherData"`
}

// WeatherDataForYear represents yearly weather data
type WeatherDataForYear struct {
	Year                 string            `json:"Year"`
	FirstRecordedDate    string            `json:"FirstRecordedDate"`
	LastRecordedDate     string            `json:"LastRecordedDate"`
	TotalRainfall        string            `json:"TotalRainfall"`
	AverageDailyRainfall string            `json:"AverageDailyRainfall"`
	DaysWithNoRainfall   string            `json:"DaysWithNoRainfall"`
	DaysWithRainfall     string            `json:"DaysWithRainfall"`
	LongestDaysRaining   string            `json:"LongestDaysRaining"`
	MonthlyAggregates    MonthlyAggregates `json:"MonthlyAggregates"`
}

// MonthlyAggregates contains monthly weather data
type MonthlyAggregates struct {
	WeatherDataForMonth []WeatherDataForMonth `json:"WeatherDataForMonth"`
}

// WeatherDataForMonth represents monthly weather data
type WeatherDataForMonth struct {
	Month                string `json:"Month"`
	FirstRecordedDate    string `json:"FirstRecordedDate"`
	LastRecordedDate     string `json:"LastRecordedDate"`
	TotalRainfall        string `json:"TotalRainfall"`
	AverageDailyRainfall string `json:"AverageDailyRainfall"`
	MedianDailyRainfall  string `json:"MedianDailyRainfall"`
	DaysWithNoRainfall   string `json:"DaysWithNoRainfall"`
	DaysWithRainfall     string `json:"DaysWithRainfall"`
}

// DailyRecord represents a single day's weather record
type DailyRecord struct {
	Date     time.Time
	Rainfall float64
	HasData  bool
}

// MonthData represents aggregated data for a month
type MonthData struct {
	Month              time.Time
	Records            []DailyRecord
	TotalRainfall      float64
	DaysWithRainfall   int
	DaysWithNoRainfall int
	FirstRecordedDate  time.Time
	LastRecordedDate   time.Time
	RainfallValues     []float64 // for median calculation
}

// YearData represents aggregated data for a year
type YearData struct {
	Year               int
	Records            []DailyRecord
	TotalRainfall      float64
	DaysWithRainfall   int
	DaysWithNoRainfall int
	LongestDaysRaining int
	FirstRecordedDate  time.Time
	LastRecordedDate   time.Time
	Months             map[time.Time]*MonthData
}
