package config

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type Config struct {
	Latitude  float64 `yaml:"Latitude" validate:"required"`
	Longitude float64 `yaml:"Longitude" validate:"required"`
}

const (
	ConfigFolder = "config"
)

func ReadConfig() *Config {
	vp := viper.New()
	vp.AddConfigPath(ConfigFolder)

	var cfg Config

	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Read error %v", err)
	}
	if err := vp.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to unmarshall the config %v", err)
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		log.Fatalf("Missing required attributes %v\n", err)
	}

	return &cfg
}
