package dto

type OpenmeteoWeatherData struct {
	Latitude  float64             `json:"latitude"`
	Longitude float64             `json:"longitude"`
	Hourly    OpenmeteoHourlyData `json:"hourly"`
}

type OpenmeteoHourlyData struct {
	Time                     []int64   `json:"time"`
	Temperature              []float64 `json:"temperature_2m"`
	PrecipitationProbability []float64 `json:"precipitation_probability"`
	WeatherCode              []int64   `json:"weathercode"`
	WindSpeed                []float64 `json:"windspeed_10m"`
}
