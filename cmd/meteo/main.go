package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Gooogr/meteo/internal/api/display"
	meteo "github.com/Gooogr/meteo/internal/api/open_meteo"
	"github.com/zsefvlol/timezonemapper"
)

// TODO:
// Ask user to type coordinates or city at first run
// Create json config file based on this input
// Take this coords from config next times. Also user can pass lat and long or city directly

func main() {
	// Get coordinates
	latFlag := flag.Float64("lat", 55.7522, "Latitude coordinate")
	lngFlag := flag.Float64("lng", 37.6156, "Longitude coordinate")
	flag.Parse()

	lat := *latFlag
	lng := *lngFlag

	// Check for valid latitude and longitude
	if lat < -90 || lat > 90 {
		fmt.Println("Invalid latitude. Latitude must be between -90 and 90 degrees.")
		return
	}
	if lng < -180 || lng > 180 {
		fmt.Println("Invalid longitude. Longitude must be between -180 and 180 degrees.")
		return
	}

	// Request forecast
	weatherData, err := meteo.GetWeatherData(meteo.Location{Lat: lat, Lng: lng})
	if err != nil {
		fmt.Printf("Error fetching weather data: %v\n", err)
		os.Exit(1)
	}

	timezone := timezonemapper.LatLngToTimezoneString(lat, lng)
	display.PrintWeatherForecast(weatherData, timezone)

}
