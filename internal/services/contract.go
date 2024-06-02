package services

import (
	"meteo/config"
	"meteo/internal/domain"
)

type Contract interface {
	Get(cfg *config.Config) (*domain.WeatherData, error)
}
