package bom

import (
	"testing"
	"time"
)

func TestNewAggregator(t *testing.T) {
	agg := NewAggregator()
	if agg == nil {
		t.Fatal("NewAggregator returned nil")
	}
}

func TestAggregate_EmptyRecords(t *testing.T) {
	agg := NewAggregator()
	result := agg.Aggregate([]DailyRecord{})

	if len(result.WeatherDataForYear) != 0 {
		t.Errorf("Expected 0 yearly aggregates, got %d", len(result.WeatherDataForYear))
	}
}

func TestAggregate_SingleYear(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.5, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.3, HasData: true},
		{Date: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), Rainfall: 5.2, HasData: true},
	}

	result := agg.Aggregate(records)

	if len(result.WeatherDataForYear) != 1 {
		t.Fatalf("Expected 1 yearly aggregate, got %d", len(result.WeatherDataForYear))
	}

	yearData := result.WeatherDataForYear[0]
	if yearData.Year != "2020" {
		t.Errorf("Expected year 2020, got %s", yearData.Year)
	}

	if yearData.FirstRecordedDate != "2020-01-01" {
		t.Errorf("Expected first date 2020-01-01, got %s", yearData.FirstRecordedDate)
	}

	if yearData.LastRecordedDate != "2020-02-01" {
		t.Errorf("Expected last date 2020-02-01, got %s", yearData.LastRecordedDate)
	}

	// Total rainfall: 10.5 + 0 + 15.3 + 5.2 = 31.0
	if yearData.TotalRainfall != "31.000000000000" {
		t.Errorf("Expected total rainfall 31.000000000000, got %s", yearData.TotalRainfall)
	}

	// Average: 31.0 / 4 = 7.75
	if yearData.AverageDailyRainfall != "7.750000000000" {
		t.Errorf("Expected average rainfall 7.750000000000, got %s", yearData.AverageDailyRainfall)
	}

	if yearData.DaysWithRainfall != "3" {
		t.Errorf("Expected 3 days with rainfall, got %s", yearData.DaysWithRainfall)
	}

	if yearData.DaysWithNoRainfall != "1" {
		t.Errorf("Expected 1 day with no rainfall, got %s", yearData.DaysWithNoRainfall)
	}

	// Check monthly aggregates
	if len(yearData.MonthlyAggregates.WeatherDataForMonth) != 2 {
		t.Errorf("Expected 2 monthly aggregates, got %d", len(yearData.MonthlyAggregates.WeatherDataForMonth))
	}
}

func TestAggregate_MultipleYears(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 20.0, HasData: true},
		{Date: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 30.0, HasData: true},
	}

	result := agg.Aggregate(records)

	if len(result.WeatherDataForYear) != 3 {
		t.Fatalf("Expected 3 yearly aggregates, got %d", len(result.WeatherDataForYear))
	}

	// Check years are sorted
	expectedYears := []string{"2020", "2021", "2022"}
	for i, yearData := range result.WeatherDataForYear {
		if yearData.Year != expectedYears[i] {
			t.Errorf("Expected year %s at position %d, got %s", expectedYears[i], i, yearData.Year)
		}
	}
}

func TestAggregate_LongestStreak(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true}, // Day 1 of streak
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 5.0, HasData: true},  // Day 2 of streak
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true}, // Day 3 of streak
		{Date: time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},  // Break
		{Date: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC), Rainfall: 8.0, HasData: true},  // Day 1 of new streak
		{Date: time.Date(2020, 1, 6, 0, 0, 0, 0, time.UTC), Rainfall: 12.0, HasData: true}, // Day 2 of new streak
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	// Longest streak should be 3 days
	if yearData.LongestDaysRaining != "3" {
		t.Errorf("Expected longest streak of 3 days, got %s", yearData.LongestDaysRaining)
	}
}

func TestAggregate_MissingData(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: false}, // Missing data
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	// Should only count records with HasData = true
	if yearData.TotalRainfall != "25.000000000000" {
		t.Errorf("Expected total rainfall 25.000000000000, got %s", yearData.TotalRainfall)
	}

	if yearData.DaysWithRainfall != "2" {
		t.Errorf("Expected 2 days with rainfall, got %s", yearData.DaysWithRainfall)
	}

	if yearData.DaysWithNoRainfall != "0" {
		t.Errorf("Expected 0 days with no rainfall, got %s", yearData.DaysWithNoRainfall)
	}
}

func TestAggregate_MonthlyBreakdown(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 5.0, HasData: true},
		{Date: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
		{Date: time.Date(2020, 2, 2, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	if len(yearData.MonthlyAggregates.WeatherDataForMonth) != 2 {
		t.Fatalf("Expected 2 monthly aggregates, got %d", len(yearData.MonthlyAggregates.WeatherDataForMonth))
	}

	// Check January
	jan := yearData.MonthlyAggregates.WeatherDataForMonth[0]
	if jan.Month != "January" {
		t.Errorf("Expected January, got %s", jan.Month)
	}
	if jan.TotalRainfall != "15.000000000000" {
		t.Errorf("Expected January total 15.000000000000, got %s", jan.TotalRainfall)
	}
	if jan.DaysWithRainfall != "2" {
		t.Errorf("Expected January 2 days with rain, got %s", jan.DaysWithRainfall)
	}
	if jan.DaysWithNoRainfall != "0" {
		t.Errorf("Expected January 0 days with no rain, got %s", jan.DaysWithNoRainfall)
	}

	// Check February
	feb := yearData.MonthlyAggregates.WeatherDataForMonth[1]
	if feb.Month != "February" {
		t.Errorf("Expected February, got %s", feb.Month)
	}
	if feb.TotalRainfall != "15.000000000000" {
		t.Errorf("Expected February total 15.000000000000, got %s", feb.TotalRainfall)
	}
	if feb.DaysWithRainfall != "1" {
		t.Errorf("Expected February 1 day with rain, got %s", feb.DaysWithRainfall)
	}
	if feb.DaysWithNoRainfall != "1" {
		t.Errorf("Expected February 1 day with no rain, got %s", feb.DaysWithNoRainfall)
	}
}

func TestAggregate_MedianCalculation(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 5.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
		{Date: time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true}, // Not included in median
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]
	jan := yearData.MonthlyAggregates.WeatherDataForMonth[0]

	// Median of [5.0, 10.0, 15.0] should be 10.0
	if jan.MedianDailyRainfall != "10.000000000000" {
		t.Errorf("Expected median 10.000000000000, got %s", jan.MedianDailyRainfall)
	}
}

func TestAggregate_MedianWithEvenNumberOfValues(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 5.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
		{Date: time.Date(2020, 1, 4, 0, 0, 0, 0, time.UTC), Rainfall: 20.0, HasData: true},
		{Date: time.Date(2020, 1, 5, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true}, // Not included in median
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]
	jan := yearData.MonthlyAggregates.WeatherDataForMonth[0]

	// Median of [5.0, 10.0, 15.0, 20.0] should be 12.5 (average of 10 and 15)
	if jan.MedianDailyRainfall != "12.500000000000" {
		t.Errorf("Expected median 12.500000000000, got %s", jan.MedianDailyRainfall)
	}
}

func TestAggregate_MonthlyMissingData(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: false}, // Missing data
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
		{Date: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), Rainfall: 5.0, HasData: true},
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	// Check January (should only count records with HasData = true)
	jan := yearData.MonthlyAggregates.WeatherDataForMonth[0]
	if jan.TotalRainfall != "25.000000000000" {
		t.Errorf("Expected January total 25.000000000000, got %s", jan.TotalRainfall)
	}
	if jan.DaysWithRainfall != "2" {
		t.Errorf("Expected January 2 days with rain, got %s", jan.DaysWithRainfall)
	}
	if jan.DaysWithNoRainfall != "0" {
		t.Errorf("Expected January 0 days with no rain, got %s", jan.DaysWithNoRainfall)
	}
}

func TestAggregate_MonthlyZeroRainfall(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},
		{Date: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	// Check January (all zero rainfall)
	jan := yearData.MonthlyAggregates.WeatherDataForMonth[0]
	if jan.TotalRainfall != "0.000000000000" {
		t.Errorf("Expected January total 0.000000000000, got %s", jan.TotalRainfall)
	}
	if jan.AverageDailyRainfall != "0.000000000000" {
		t.Errorf("Expected January average 0.000000000000, got %s", jan.AverageDailyRainfall)
	}
	if jan.MedianDailyRainfall != "0.000000000000" {
		t.Errorf("Expected January median 0.000000000000, got %s", jan.MedianDailyRainfall)
	}
	if jan.DaysWithRainfall != "0" {
		t.Errorf("Expected January 0 days with rain, got %s", jan.DaysWithRainfall)
	}
	if jan.DaysWithNoRainfall != "2" {
		t.Errorf("Expected January 2 days with no rain, got %s", jan.DaysWithNoRainfall)
	}
}

func TestAggregate_MonthlyDateRanges(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 31, 0, 0, 0, 0, time.UTC), Rainfall: 20.0, HasData: true},
		{Date: time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
		{Date: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC), Rainfall: 25.0, HasData: true}, // Leap year
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	// Check January date range
	jan := yearData.MonthlyAggregates.WeatherDataForMonth[0]
	if jan.FirstRecordedDate != "2020-01-01" {
		t.Errorf("Expected January first date 2020-01-01, got %s", jan.FirstRecordedDate)
	}
	if jan.LastRecordedDate != "2020-01-31" {
		t.Errorf("Expected January last date 2020-01-31, got %s", jan.LastRecordedDate)
	}

	// Check February date range
	feb := yearData.MonthlyAggregates.WeatherDataForMonth[1]
	if feb.FirstRecordedDate != "2020-02-01" {
		t.Errorf("Expected February first date 2020-02-01, got %s", feb.FirstRecordedDate)
	}
	if feb.LastRecordedDate != "2020-02-29" {
		t.Errorf("Expected February last date 2020-02-29, got %s", feb.LastRecordedDate)
	}
}

func TestAggregate_ZeroRainfall(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 0.0, HasData: true},
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	if yearData.TotalRainfall != "0.000000000000" {
		t.Errorf("Expected total rainfall 0.000000000000, got %s", yearData.TotalRainfall)
	}

	if yearData.AverageDailyRainfall != "0.000000000000" {
		t.Errorf("Expected average rainfall 0.000000000000, got %s", yearData.AverageDailyRainfall)
	}

	if yearData.DaysWithRainfall != "0" {
		t.Errorf("Expected 0 days with rainfall, got %s", yearData.DaysWithRainfall)
	}

	if yearData.DaysWithNoRainfall != "2" {
		t.Errorf("Expected 2 days with no rainfall, got %s", yearData.DaysWithNoRainfall)
	}

	if yearData.LongestDaysRaining != "0" {
		t.Errorf("Expected longest streak 0, got %s", yearData.LongestDaysRaining)
	}
}

func TestAggregate_RecordsNotSorted(t *testing.T) {
	agg := NewAggregator()
	records := []DailyRecord{
		{Date: time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC), Rainfall: 15.0, HasData: true},
		{Date: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Rainfall: 10.0, HasData: true},
		{Date: time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC), Rainfall: 5.0, HasData: true},
	}

	result := agg.Aggregate(records)
	yearData := result.WeatherDataForYear[0]

	// Should be sorted by date
	if yearData.FirstRecordedDate != "2020-01-01" {
		t.Errorf("Expected first date 2020-01-01, got %s", yearData.FirstRecordedDate)
	}

	if yearData.LastRecordedDate != "2020-01-03" {
		t.Errorf("Expected last date 2020-01-03, got %s", yearData.LastRecordedDate)
	}
}
