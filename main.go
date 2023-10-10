package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/zsefvlol/timezonemapper"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"

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
	Weathercode              []float64 `json:"weathercode"`
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

// "https://api.open-meteo.com/v1/forecast?latitude=55.7522&longitude=37.6156&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&timezone=Europe%2FMoscow&forecast_days=3"

func main() {
	// TODO: pass coords as arguments of CLI or get coords from device
	var lat, lng float64 = 55.7522, 37.6156
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
	// fmt.Println(weather)

	for i, time := range weather.Hourly.Time {
		temperature := weather.Hourly.Temperature2m[i]
		weathercode := weather.Hourly.Weathercode[i]
		precipitation_probability := weather.Hourly.PrecipitationProbability[i]
		windspeed := weather.Hourly.Windspeed10m[i]
		fmt.Printf(
			"%s - %.1f°C %.0f, %.0f  %.1fkm/h \n",
			time,
			temperature,
			weathercode,
			precipitation_probability,
			windspeed,
		)
	}
}
