package domain

type WeatherData struct {
	Time                     []int64
	Temperature              []float64
	PrecipitationProbability []float64
	WeatherState             []string
	WindSpeed                []float64
}
