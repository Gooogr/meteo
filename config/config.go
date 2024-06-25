package config

import (
	"log"
	"os"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Latitude                 float64 `mapstructure:"latitude" validate:"required"`
	Longitude                float64 `mapstructure:"longitude" validate:"required"`
	DefaultAPI               string  `mapstructure:"default-api" validate:"required"`
	MeteoblueAPIKey          string  `mapstructure:"meteoblue-api-key"`
	MeteoblueAPISharedSecret string  `mapstructure:"meteoblue-shared-secret"`
}

func ReadConfigFile() *Config {
	// Read default yaml config file
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	configDir := os.Getenv("METEO_CONFIG_PATH")
	if configDir == "" {
		log.Fatalf("METEO_CONFIG_PATH environment variable not set")
	}

	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read error %v", err)
	}

	// Validate config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to unmarshall the config %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		log.Fatalf("missing required attributes %v\n", err)
	}
	return &cfg
}

func UpdateConfigFile(key string, value interface{}) {
	// Read default yaml config file
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	configDir := os.Getenv("METEO_CONFIG_PATH")
	if configDir == "" {
		log.Fatalf("METEO_CONFIG_PATH environment variable not set")
	}

	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read error %v", err)
	}

	// Convert config to map
	configMap := viper.AllSettings()

	// Update value in map by key
	if _, ok := configMap[key]; !ok {
		log.Fatalf("can't find key in existed config file: %v", key)
	}
	if reflect.TypeOf(configMap[key]) != reflect.TypeOf(value) {
		log.Fatalf("inconsistent value type for key: %v", key)
	}

	configMap[key] = value

	// Re-write yaml file
	updatedConfigMap, err := yaml.Marshal(&configMap)
	if err != nil {
		log.Fatalf("failed to marshal updated YAML data: %v", err)
	}

	err = os.WriteFile(viper.ConfigFileUsed(), updatedConfigMap, 0644)
	if err != nil {
		log.Fatalf("failed to write updated YAML file: %v", err)
	}
}
