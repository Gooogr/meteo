package openmeteo

import (
	"encoding/json"
	"fmt"
	"io"

	"meteo/config"
	"meteo/internal/domain"
	"meteo/internal/dto"
	"meteo/internal/services"
)

const baseURL = "https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f"

var openmeteoWeatherCodes = map[int64]string{
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

type openmeteo struct {
	client httpClient
}

func NewOpenmeteo(client httpClient) services.Contract {
	return &openmeteo{
		client: client,
	}
}

func (om *openmeteo) fetchOpenmeteoData(url string) (*domain.OpenmeteoWeatherData, error) {
	weatherDto := dto.OpenmeteoWeatherData{}

	// Get data from openmeteo.
	resp, err := om.client.Get(url)
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

	err = json.Unmarshal(body, &weatherDto)
	if err != nil {
		return nil, err
	}

	// Convert dto to domain.
	data := &domain.OpenmeteoWeatherData{
		Latitude:  weatherDto.Latitude,
		Longitude: weatherDto.Longitude,
		Hourly: domain.OpenmeteoHourlyData{
			Time:                     weatherDto.Hourly.Time,
			Temperature:              weatherDto.Hourly.Temperature,
			PrecipitationProbability: weatherDto.Hourly.PrecipitationProbability,
			WeatherCode:              weatherDto.Hourly.WeatherCode,
			WindSpeed:                weatherDto.Hourly.WindSpeed,
		},
	}

	return data, nil
}

func (om *openmeteo) Get(cfg *config.Config) (*domain.WeatherData, error) {
	url, err := createURL(cfg.Latitude, cfg.Longitude)
	if err != nil {
		return nil, err
	}

	data, err := om.fetchOpenmeteoData(url)
	if err != nil {
		return nil, err
	}

	weatherState := make([]string, len(data.Hourly.WeatherCode))
	for i, code := range data.Hourly.WeatherCode {
		state := ""

		wc, ok := openmeteoWeatherCodes[code]
		if ok {
			state = wc
		}

		weatherState[i] = state
	}

	return &domain.WeatherData{
		Time:                     data.Hourly.Time,
		Temperature:              data.Hourly.Temperature,
		PrecipitationProbability: data.Hourly.PrecipitationProbability,
		WeatherState:             weatherState,
		WindSpeed:                data.Hourly.WindSpeed,
	}, nil
}

func createURL(lat float64, lng float64) (string, error) {
	if lat < -90 || lat > 90 {
		return "", fmt.Errorf("latitude must be between -90 and 90 degrees")
	}
	if lng < -180 || lng > 180 {
		return "", fmt.Errorf("longitude must be between -180 and 180 degrees")
	}

	url := fmt.Sprintf(baseURL, lat, lng)
	url = url + "&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3&timeformat=unixtime"

	return url, nil
}
