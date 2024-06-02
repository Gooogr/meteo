package domain

type MeteoblueWeatherData struct {
	MeteoblueMetadata MeteoblueMetadata
	MeteoblueData1h   MeteoblueData1h
}

type MeteoblueMetadata struct {
	Latitude  float64
	Longitude float64
}

type MeteoblueData1h struct {
	Time                     []int64
	Temperature              []float64
	PrecipitationProbability []float64
	WindSpeed                []float64
	Pictocode                []int64
}
