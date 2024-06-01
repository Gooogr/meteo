package main

import (
	"fmt"
	"net/http"
	"os"

	"meteo/config"
	"meteo/internal/display"
	"meteo/internal/services/openmeteo"

	"github.com/zsefvlol/timezonemapper"
)

// TODO:
// Ask user to type coordinates or city at first run
// Create json config file based on this input
// Take this coords from config next times. Also user can pass lat and long or city directly

func main() {
	cfg := config.ReadConfig()
	lat := cfg.Latitude
	lng := cfg.Longitude

	// Init http client.
	httpClient := &http.Client{}

	// Init services.
	weatherService := openmeteo.NewOpenmeteo(httpClient)

	// Get weather data
	weatherData, err := weatherService.Get(lat, lng)
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		os.Exit(1)
	}

	// Render table
	timezone := timezonemapper.LatLngToTimezoneString(lat, lng)
	display.DisplayTable(weatherData, timezone)
}
