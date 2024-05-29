package display

import (
	"fmt"
	"os"
	"time"

	openmeteo "meteo/internal/api/open_meteo"

	"github.com/olekukonko/tablewriter"
)

const maxRows = 12

var weatherCodes = map[int]string{
	0:  "Clear sky",
	1:  "Mainly clear",
	2:  "Partly cloudy",
	3:  "Overcast",
	45: "Fog",
	48: "Rime fog",
	51: "Light drizzle",
	53: "Moderate drizzle",
	55: "Dense drizzle",
	56: "Light freezing drizzle",
	57: "Dense freezing drizzle",
	61: "Slight rain",
	63: "Moderate rain",
	65: "Heavy rain",
	66: "Light freezing rain",
	67: "Heavy freezing rain",
	71: "Slight snow",
	73: "Moderate snow",
	75: "Heavy snow",
	77: "Snow grains",
	80: "Slight rain",
	81: "Moderate rain",
	82: "Heavy rain",
	85: "Slight snow",
	86: "Heavy snow",
	95: "Thunderstorm",
	96: "Thunderstorm with slight hail",
	99: "Thunderstorm with heavy hail",
}

func prepareWeatherData(weather openmeteo.WeatherData, timezone string) [][]string {
	currentTime := time.Now()

	var data [][]string

	i := 0
	rowsCnt := 0
	for i < len(weather.Hourly.Time) {
		datetime := weather.Hourly.Time[i]
		temperature := weather.Hourly.Temperature2m[i]
		weatherCode := weather.Hourly.WeatherCode[i]
		precipitationProbability := weather.Hourly.PrecipitationProbability[i]
		windSpeed := weather.Hourly.WindSpeed10m[i]

		weatherState, ok := weatherCodes[weatherCode]
		if !ok {
			weatherState = ""
		}

		location, err := time.LoadLocation(timezone)
		if err != nil {
			fmt.Println("Error loading timezone:", err)
			os.Exit(1)
		}
		datetimeInLocation := datetime.In(location)
		if err != nil {
			fmt.Println("Error parsing datetime:", err)
			os.Exit(1)
		}

		if datetimeInLocation.Before(currentTime) {
			i += 1
			continue
		}
		if rowsCnt > maxRows {
			break
		}

		hour := datetimeInLocation.Hour()
		formattedHour := fmt.Sprintf("%02d:00", hour)
		formattedTemperature := fmt.Sprintf("%.1f°C", temperature)
		formattedPrecipitationProbability := fmt.Sprintf("%.0f%%", precipitationProbability)
		formattedWindSpeed := fmt.Sprintf("%.1fkm/h", windSpeed)

		data = append(data, []string{
			formattedHour,
			formattedTemperature,
			formattedPrecipitationProbability,
			formattedWindSpeed,
			weatherState,
		})
		i += 1
		rowsCnt += 1
	}
	return data
}

func DisplayTable(weather openmeteo.WeatherData, timezone string) {
	data := prepareWeatherData(weather, timezone)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Time\n----",
		"Temp\n----",
		"Rain\n----",
		"Wind\n----",
		"Condition\n---------",
	})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("  ") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
