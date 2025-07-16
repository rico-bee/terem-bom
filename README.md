# terem-bom

## Background

Terem Technologies is very curious about weather trends at Observatory Hill in Sydney. Luckily, the Bureau of Meteorology has recorded all of the historical rainfall data and it is downloadable here:
http://www.bom.gov.au/jsp/ncc/cdio/weatherData/av?p_nccObsCode=136&p_display_type=dailyDataFile&p_startYear=&p_c=&p_stn_num=066062

And click on the button "All years of data" on the top right of the page.

## Task

You are to create a library which reads any given BOM weather data CSV fi le and converts the data to a JSON of the following format:

```
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
      "MonthlyAggregates": {
        "WeatherDataForMonth": {
          "Month": "January",
          "FirstRecordedDate": "2019-01-01",
          "LastRecordedDate": "2019-01-31",
          "TotalRainfall": "48.8",
          "AverageDailyRainfall": "1.574193548",
          "DaysWithNoRainfall": "21",
          "DaysWithRainfall": "10"
        }
      }
    }
  }
}
```

## Acceptance Criteria

Create a CLI tool that you can point to CSV data.
- All data in the CSV will be converted to a corresponding JSON output, except:
- Dates where the “Rainfall amount (millimetres)” is empty / blank should not be counted / recorded when determining FirstRecordedDate / LastRecordedDate
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
- Months that have yet to occur should not be included in the output data (i.e. If it’s currently January 2000, a MonthlyAggregate node should not exist for February 2000)


## Evaluation Criteria

- Technology best practices
- Show us your work through your commit history
- We're looking for you to produce working code, with enough room to demonstrate how to structure components in a small program
- Completeness: did you complete the features?
- Correctness: does the functionality act in sensible, thought-out ways?
- Maintainability: is it written in a clean, maintainable way?
- Testing: is the system adequately tested?

