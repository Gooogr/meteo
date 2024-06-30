package services

import (
	"meteo/internal/domain"
)

type Contract interface {
	Get(latitude, longitude float64) (*domain.WeatherData, error)
}
