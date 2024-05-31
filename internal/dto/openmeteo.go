package dto

import (
	"encoding/json"
	"time"
)

type WeatherData struct {
	Latitude  float64    `json:"latitude"`
	Longitude float64    `json:"longitude"`
	Hourly    HourlyData `json:"hourly"`
}

type TimeSlice []time.Time

type HourlyData struct {
	Time                     TimeSlice `json:"time"`
	Temperature2m            []float64 `json:"temperature_2m"`
	PrecipitationProbability []float64 `json:"precipitation_probability"`
	WeatherCode              []int     `json:"weathercode"` // Специфичен для этого API
	WindSpeed10m             []float64 `json:"windspeed_10m"`
}

func (ts *TimeSlice) UnmarshalJSON(data []byte) error {
	var timeStrings []string
	if err := json.Unmarshal(data, &timeStrings); err != nil {
		return err
	}
	*ts = make([]time.Time, len(timeStrings))
	for idx, t := range timeStrings {
		parsedTime, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			return err
		}
		(*ts)[idx] = parsedTime
	}
	return nil
}
