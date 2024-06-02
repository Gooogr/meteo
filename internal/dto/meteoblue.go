package dto

type MeteoblueWeatherData struct {
	MeteoblueMetadata MeteoblueMetadata `json:"metadata"`
	MeteoblueData1h   MeteoblueData1h   `json:"data_1h"`
}

type MeteoblueMetadata struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MeteoblueData1h struct {
	Time                     []int64   `json:"time"`
	Temperature              []float64 `json:"temperature"`
	PrecipitationProbability []float64 `json:"precipitation_probability"`
	WindSpeed                []float64 `json:"windspeed"`
	Pictocode                []int64   `json:"pictocode"`
}
