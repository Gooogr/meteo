package display

import (
	"fmt"
	"os"
	"time"

	"meteo/internal/domain"

	"github.com/olekukonko/tablewriter"
)

const maxRows = 12

func prepareWeatherData(weather *domain.WeatherData, timezone string) [][]string {
	currentTime := time.Now()

	var data [][]string

	i := 0
	rowsCnt := 0
	for i < len(weather.Time) {
		datetime := time.Unix(weather.Time[i], 0)
		temperature := weather.Temperature[i]
		weatherState := weather.WeatherState[i]
		precipitationProbability := weather.PrecipitationProbability[i]
		windSpeed := weather.WindSpeed[i]

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

func DisplayTable(weather *domain.WeatherData, timezone string) {
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
	table.SetTablePadding("  ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
