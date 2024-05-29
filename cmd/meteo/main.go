package main

import (
	"fmt"
	openmeteo "meteo/internal/api/open_meteo"
	"meteo/internal/config"
	"meteo/internal/display"
	"os"

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

	if lat < -90 || lat > 90 {
		fmt.Println("Invalid latitude. Latitude must be between -90 and 90 degrees.")
		return
	}
	if lng < -180 || lng > 180 {
		fmt.Println("Invalid longitude. Longitude must be between -180 and 180 degrees.")
		return
	}

	weatherData, err := openmeteo.GetOpenMeteoData(lat, lng)
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		os.Exit(1)
	}

	timezone := timezonemapper.LatLngToTimezoneString(lat, lng)
	display.DisplayTable(weatherData, timezone)

}
