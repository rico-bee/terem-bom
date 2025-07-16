package weather

import (
	"time"
)

// DailyWeatherRecord represents a single day's weather data
type DailyWeatherRecord struct {
	Date     time.Time `json:"Date"`
	Year     int       `json:"Year"`
	Month    int       `json:"Month"`
	Day      int       `json:"Day"`
	Rainfall float64   `json:"Rainfall"`
}

// MonthlyAggregate represents aggregated weather data for a month
type MonthlyAggregate struct {
	Month                string `json:"Month"`
	FirstRecordedDate    string `json:"FirstRecordedDate"`
	LastRecordedDate     string `json:"LastRecordedDate"`
	TotalRainfall        string `json:"TotalRainfall"`
	AverageDailyRainfall string `json:"AverageDailyRainfall"`
	MedianDailyRainfall  string `json:"MedianDailyRainfall"`
	DaysWithNoRainfall   int    `json:"DaysWithNoRainfall"`
	DaysWithRainfall     int    `json:"DaysWithRainfall"`
}

// MonthlyAggregatesList holds the monthly aggregates
type MonthlyAggregatesList struct {
	WeatherDataForMonth []MonthlyAggregate `json:"WeatherDataForMonth"`
}

// YearlyAggregate represents aggregated weather data for a year
type YearlyAggregate struct {
	Year                 string                `json:"Year"`
	FirstRecordedDate    string                `json:"FirstRecordedDate"`
	LastRecordedDate     string                `json:"LastRecordedDate"`
	TotalRainfall        string                `json:"TotalRainfall"`
	AverageDailyRainfall string                `json:"AverageDailyRainfall"`
	DaysWithNoRainfall   int                   `json:"DaysWithNoRainfall"`
	DaysWithRainfall     int                   `json:"DaysWithRainfall"`
	LongestDaysRaining   int                   `json:"LongestDaysRaining"`
	MonthlyAggregates    MonthlyAggregatesList `json:"MonthlyAggregates"`
}

// WeatherDataList holds the complete weather data structure
type WeatherDataList struct {
	WeatherData []YearlyAggregate `json:"WeatherData"`
}

// ProcessingOptions holds configuration options for data processing
type ProcessingOptions struct {
	PrettyPrint bool
	Verbose     bool
}

// ProcessingResult holds the result of data processing
type ProcessingResult struct {
	Data    *WeatherDataList
	Error   error
	Records int
}
