package display

import (
	"fmt"
	"os"
	"time"

	meteo "meteo/internal/api/open_meteo"

	"github.com/fatih/color"
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

func PrintWeatherForecast(weather meteo.WeatherData, timezone string) {
	currentTime := time.Now()

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
		datetimeParsed, err := time.ParseInLocation("2006-01-02T15:04", datetime, location)
		if err != nil {
			fmt.Println("Error parsing datetime:", err)
			os.Exit(1)
		}

		if datetimeParsed.Before(currentTime) {
			i += 1
			continue
		}
		if rowsCnt > maxRows {
			break
		}

		hour := datetimeParsed.Hour()
		formattedHour := fmt.Sprintf("%02d:00", hour)

		temperatureColor := color.New(color.FgCyan).Sprintf("%.1f°C", temperature)
		precipitationProbabilityColor := color.New(color.FgGreen).Sprintf("%.0f%%", precipitationProbability)
		windSpeedColor := color.New(color.FgBlue).Sprintf("%.1fkm/h", windSpeed)

		timeColumnSize := 6
		temperatureColumnSize := 10
		windSpeedColumnSize := 10
		precipitationProbablityColumnSize := 13
		weatherColumnSize := 10

		formatString := fmt.Sprintf(
			"%%-%ds %%-%ds %%-%ds %%-%ds  %%-%ds\n",
			timeColumnSize,
			temperatureColumnSize,
			windSpeedColumnSize,
			precipitationProbablityColumnSize,
			weatherColumnSize,
		)

		fmt.Printf(
			formatString,
			formattedHour,
			temperatureColor,
			windSpeedColor,
			precipitationProbabilityColor,
			weatherState,
		)
		i += 1
		rowsCnt += 1
	}
}
