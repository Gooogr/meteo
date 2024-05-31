package domain

import "time"

type OpenmeteoWeatherData struct {
	Latitude  float64
	Longitude float64
	Hourly    HourlyData
}

type TimeSlice []time.Time

type HourlyData struct {
	Time                     TimeSlice
	Temperature2m            []float64
	PrecipitationProbability []float64
	WeatherCode              []int
	WindSpeed10m             []float64
}
