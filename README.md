# terem-bom

A Go library and CLI tool for converting Bureau of Meteorology (BOM) weather data CSV files to structured JSON format.

## Quick Start

### Build and Convert

1. **Build the application:**
   ```bash
   make build
   ```

2. **Convert a CSV file to JSON:**
   ```bash
   ./bin/bom convert -i your_data.csv -o output.json
   ```

### Example

```bash
# Build the tool
make build

# Convert BOM weather data
./bin/bom convert -i test_Data/IDCJAC0009_066062_1800_Data.csv -o weather_output.json

# Validate a CSV file
./bin/bom validate -i test_Data/IDCJAC0009_066062_1800_Data.csv

# Show help
./bin/bom --help
```

## Background

The Bureau of Meteorology has recorded all of the historical rainfall data and it is downloadable here:
http://www.bom.gov.au/jsp/ncc/cdio/weatherData/av?p_nccObsCode=136&p_display_type=dailyDataFile&p_startYear=&p_c=&p_stn_num=066062

And click on the button "All years of data" on the top right of the page.

## Task

You are to create a library which reads any given BOM weather data CSV file and converts the data to a JSON of the following format:

Note: I tweaked the following format slightly to support multiple year CSV files

```json
{
  "WeatherData": {
    "WeatherDataForYear": {
      "Year": "2019",
      "FirstRecordedDate": "2019-01-01",
      "LastRecordedDate": "2019-04-19",
      "TotalRainfall": "374.2",
      "AverageDailyRainfall": "3.433027523",
      "DaysWithNoRainfall": "65",
      "DaysWithRainfall": "44",
      "LongestDaysRaining": "5",
      "MonthlyAggregates": {
        "WeatherDataForMonth": [
          {
            "Month": "January",
            "FirstRecordedDate": "2019-01-01",
            "LastRecordedDate": "2019-01-31",
            "TotalRainfall": "48.8",
            "AverageDailyRainfall": "1.574193548",
            "MedianDailyRainfall": "0.0",
            "DaysWithNoRainfall": "21",
            "DaysWithRainfall": "10"
          }
        ]
      }
    }
  }
}
```

## Acceptance Criteria

Create a CLI tool that you can point to CSV data.
- All data in the CSV will be converted to a corresponding JSON output, except:
- Dates where the "Rainfall amount (millimetres)" is empty / blank should not be counted / recorded when determining FirstRecordedDate / LastRecordedDate
- A year data should contain:
  - Year value
  - First and last recorded dates
  - Total rainfall
  - Average daily rainfall
  - Days with rainfall
  - Days with no rainfall
  - Longest number days raining
  - Monthly Aggregates
- A month data should contain:
  - Month name
  - First and last recorded dates
  - Total rainfall
  - Average daily rainfall
  - Median Daily rainfall
  - Days with rainfall
  - Days with no rainfall
- Months that have yet to occur should not be included in the output data (i.e. If it's currently January 2000, a MonthlyAggregate node should not exist for February 2000)

## Features

- **CSV Parsing**: Robust parsing of BOM weather data CSV files
- **Data Aggregation**: Yearly and monthly data aggregation with statistics
- **JSON Output**: Structured JSON output matching the specified format
- **CLI Interface**: Command-line tool with flexible options
- **Error Handling**: Comprehensive error handling and validation
- **Testing**: Extensive unit tests with high coverage
- **Build Automation**: Makefile for easy building and testing



