package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/zsefvlol/timezonemapper"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"
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

type WeatherData struct {
	Latitude    float64         `json:"latitude"`
	Longitude   float64         `json:"longitude"`
	HourlyUnits HourlyUnitsData `json:"hourly_units"`
	Hourly      HourlyData      `json:"hourly"`
}

type HourlyUnitsData struct {
	Time                     string `json:"time"`
	Temperature2m            string `json:"temperature_2m"`
	PrecipitationProbability string `json:"precipitation_probability"`
	Weathercode              string `json:"weathercode"`
	Windspeed10m             string `json:"windspeed_10m"`
}

type HourlyData struct {
	Time                     []string  `json:"time"` // TODO: convert directly to time.Time
	Temperature2m            []float64 `json:"temperature_2m"`
	PrecipitationProbability []float64 `json:"precipitation_probability"`
	Weathercode              []int     `json:"weathercode"`
	Windspeed10m             []float64 `json:"windspeed_10m"`
}

type Location struct {
	lat, lng float64
}

func createURL(apiURL string, loc Location) string {
	// Add coordinates
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f", apiURL, loc.lat, loc.lng)
	// Add timezone
	timezone := timezonemapper.LatLngToTimezoneString(loc.lat, loc.lng)
	url = fmt.Sprintf("%s&timezone=%s", url, timezone)
	// Add other constant filters
	url = url + "&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3"
	return url
}

func main() {
	// TODO: pass coords as arguments of CLI or get coords from device
	var lat, lng float64 = 55.7522, 37.6156
	timezone := timezonemapper.LatLngToTimezoneString(lat, lng) // TODO: refactor to avoid repeated code in createURL
	res, err := http.Get(createURL(apiURL, Location{lat, lng}))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var weather WeatherData
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	currentTime := time.Now()
	// fmt.Println(currentTime)

	i := 0
	rowsCnt := 0
	for i < len(weather.Hourly.Time) {
		datetime := weather.Hourly.Time[i]
		temperature := weather.Hourly.Temperature2m[i]
		weathercode := weather.Hourly.Weathercode[i]
		precipitation := weather.Hourly.PrecipitationProbability[i]
		windspeed := weather.Hourly.Windspeed10m[i]

		// Decode weather conditions
		weatherState, ok := weatherCodes[weathercode]
		if !ok {
			weatherState = ""
		}

		// Parse the datetime string
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

		// Skip historical data
		if datetimeParsed.Before(currentTime) {
			i += 1
			continue
		}
		// Limit output size
		if rowsCnt > maxRows {
			break
		}

		// Extract and format the hour
		hour := datetimeParsed.Hour()
		formattedHour := fmt.Sprintf("%02d:00", hour)

		// Define colors
		temperatureColor := color.New(color.FgCyan).Sprintf("%.1f°C", temperature)
		precipitationColor := color.New(color.FgGreen).Sprintf("%.0f%%", precipitation)
		windspeedColor := color.New(color.FgBlue).Sprintf("%.1fkm/h", windspeed)

		formattedRow := fmt.Sprintf(
			"%7v %15v %15v %12v %s",
			formattedHour,
			temperatureColor,
			windspeedColor,
			precipitationColor,
			weatherState,
		)

		fmt.Println(formattedRow)
		i += 1
		rowsCnt += 1
	}
}
