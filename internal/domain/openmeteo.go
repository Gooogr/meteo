package domain

type OpenmeteoWeatherData struct {
	Latitude  float64
	Longitude float64
	Hourly    OpenmeteoHourlyData
}

type OpenmeteoHourlyData struct {
	Time                     []int64
	Temperature              []float64
	PrecipitationProbability []float64
	WeatherCode              []int64
	WindSpeed                []float64
}
