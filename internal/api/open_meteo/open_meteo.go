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
	WeatherCode              string `json:"weathercode"`
	WindSpeed10m             string `json:"windspeed_10m"`
}

type HourlyData struct {
	Time                     []string  `json:"time"` // TODO: convert directly to time.Time
	Temperature2m            []float64 `json:"temperature_2m"`
	PrecipitationProbability []float64 `json:"precipitation_probability"`
	WeatherCode              []int     `json:"weathercode"`
	WindSpeed10m             []float64 `json:"windspeed_10m"`
}

func createURL(lat float64, lng float64) string {
	timezone := timezonemapper.LatLngToTimezoneString(lat, lng)
	url := fmt.Sprintf("%s?latitude=%f&longitude=%f&timezone=%s", apiURL, lat, lng, timezone)
	url = url + "&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3"
	return url
}

type httpClient interface {
	Get(string) (*http.Response, error)
}

type openMeteoGiver struct {
	client httpClient
}

func (giver openMeteoGiver) get(url string) ([]byte, error) {
	resp, err := giver.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetOpenMeteoData(lat float64, lng float64) (WeatherData, error) {
	var weather WeatherData

	client := &http.Client{}
	giver := openMeteoGiver{client: client}

	url := createURL(lat, lng)
	body, err := giver.get(url)
	if err != nil {
		return weather, err
	}
	err = json.Unmarshal(body, &weather)
	return weather, err
}
