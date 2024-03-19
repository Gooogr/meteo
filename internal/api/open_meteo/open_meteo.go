package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
	Weathercode              []int     `json:"weathercode"`
	Windspeed10m             []float64 `json:"windspeed_10m"`
}

type Location struct {
	Lat, Lng float64
}

func createURL(loc Location) string {
	// Add coordinates
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f", apiURL, loc.Lat, loc.Lng)
	// Add timezone
	timezone := timezonemapper.LatLngToTimezoneString(loc.Lat, loc.Lng)
	url = fmt.Sprintf("%s&timezone=%s", url, timezone)
	// Add other constant filters
	url = url + "&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3"
	return url
}

func GetWeatherData(loc Location) (WeatherData, error) {
	var weather WeatherData
	res, err := http.Get(createURL(loc))
	if err != nil {
		return weather, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return weather, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return weather, err
	}

	err = json.Unmarshal(body, &weather)
	return weather, err
}
