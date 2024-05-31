package domain

type WeatherData struct {
	Time                     TimeSlice
	Temperature2m            []float64
	PrecipitationProbability []float64
	WeatherState             []string
	WindSpeed10m             []float64
}
