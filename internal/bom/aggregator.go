package bom

import (
	"sort"
	"strconv"
	"time"
)

// Aggregator computes statistics from daily weather records
type Aggregator struct{}

// NewAggregator creates a new Aggregator
func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// Aggregate aggregates daily records into yearly and monthly statistics
func (a *Aggregator) Aggregate(records []DailyRecord) WeatherData {
	yearMap := make(map[int][]DailyRecord)
	for _, rec := range records {
		year := rec.Date.Year()
		yearMap[year] = append(yearMap[year], rec)
	}

	var years []int
	for y := range yearMap {
		years = append(years, y)
	}
	sort.Ints(years)

	var yearlyAggregates []WeatherDataForYear
	for _, year := range years {
		yearRecords := yearMap[year]
		yearlyAggregates = append(yearlyAggregates, a.aggregateYear(year, yearRecords))
	}

	return WeatherData{WeatherDataForYear: yearlyAggregates}
}

func (a *Aggregator) aggregateYear(year int, records []DailyRecord) WeatherDataForYear {
	if len(records) == 0 {
		return WeatherDataForYear{}
	}
	// Sort records by date
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})

	firstDate := records[0].Date.Format("2006-01-02")
	lastDate := records[len(records)-1].Date.Format("2006-01-02")

	var totalRainfall float64
	var daysWithRainfall, daysWithNoRainfall, longestStreak, currentStreak int
	var prevRained bool

	monthMap := make(map[time.Month][]DailyRecord)
	for _, rec := range records {
		if rec.HasData {
			totalRainfall += rec.Rainfall
			if rec.Rainfall > 0 {
				daysWithRainfall++
				if prevRained {
					currentStreak++
				} else {
					currentStreak = 1
				}
				if currentStreak > longestStreak {
					longestStreak = currentStreak
				}
				prevRained = true
			} else {
				daysWithNoRainfall++
				prevRained = false
			}
			monthMap[rec.Date.Month()] = append(monthMap[rec.Date.Month()], rec)
		}
	}

	totalDays := daysWithRainfall + daysWithNoRainfall
	avgRain := 0.0
	if totalDays > 0 {
		avgRain = totalRainfall / float64(totalDays)
	}

	// Monthly aggregates
	var months []time.Month
	for m := range monthMap {
		months = append(months, m)
	}
	sort.Slice(months, func(i, j int) bool { return months[i] < months[j] })

	var monthlyAggregates []WeatherDataForMonth
	for _, m := range months {
		monthlyAggregates = append(monthlyAggregates, a.aggregateMonth(m, monthMap[m]))
	}

	return WeatherDataForYear{
		Year:                 strconv.Itoa(year),
		FirstRecordedDate:    firstDate,
		LastRecordedDate:     lastDate,
		TotalRainfall:        formatFloat(totalRainfall, 1),
		AverageDailyRainfall: formatFloat(avgRain, 12),
		DaysWithNoRainfall:   strconv.Itoa(daysWithNoRainfall),
		DaysWithRainfall:     strconv.Itoa(daysWithRainfall),
		LongestDaysRaining:   strconv.Itoa(longestStreak),
		MonthlyAggregates:    MonthlyAggregates{WeatherDataForMonth: monthlyAggregates},
	}
}

func (a *Aggregator) aggregateMonth(month time.Month, records []DailyRecord) WeatherDataForMonth {
	if len(records) == 0 {
		return WeatherDataForMonth{}
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})
	firstDate := records[0].Date.Format("2006-01-02")
	lastDate := records[len(records)-1].Date.Format("2006-01-02")

	var totalRainfall float64
	var daysWithRainfall, daysWithNoRainfall int
	var rainfallDays []float64

	for _, rec := range records {
		if rec.HasData {
			totalRainfall += rec.Rainfall
			if rec.Rainfall > 0 {
				daysWithRainfall++
				rainfallDays = append(rainfallDays, rec.Rainfall)
			} else {
				daysWithNoRainfall++
			}
		}
	}

	totalDays := daysWithRainfall + daysWithNoRainfall
	avgRain := 0.0
	if totalDays > 0 {
		avgRain = totalRainfall / float64(totalDays)
	}

	medianRain := 0.0
	if len(rainfallDays) > 0 {
		sort.Float64s(rainfallDays)
		if len(rainfallDays)%2 == 0 {
			// Even number of values: average of two middle values
			mid := len(rainfallDays) / 2
			medianRain = (rainfallDays[mid-1] + rainfallDays[mid]) / 2.0
		} else {
			// Odd number of values: middle value
			medianRain = rainfallDays[len(rainfallDays)/2]
		}
	}

	return WeatherDataForMonth{
		Month:                month.String(),
		FirstRecordedDate:    firstDate,
		LastRecordedDate:     lastDate,
		TotalRainfall:        formatFloat(totalRainfall, 1),
		AverageDailyRainfall: formatFloat(avgRain, 12),
		MedianDailyRainfall:  formatFloat(medianRain, 12),
		DaysWithNoRainfall:   strconv.Itoa(daysWithNoRainfall),
		DaysWithRainfall:     strconv.Itoa(daysWithRainfall),
	}
}

// formatFloat formats a float64 to string with specified precision (defaults to 9)
func formatFloat(f float64, precision ...int) string {
	p := 9 // default precision
	if len(precision) > 0 {
		p = precision[0]
	}
	return strconv.FormatFloat(f, 'f', p, 64)
}
