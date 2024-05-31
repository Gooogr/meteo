package internal

import "meteo/internal/domain"

type Contract interface {
	Get(lat float64, lng float64) (*domain.WeatherData, error)
}
