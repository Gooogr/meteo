package cmd

import (
	"fmt"
	"meteo/config"
	"meteo/internal/display"
	"meteo/internal/services"
	"meteo/internal/services/meteoblue"
	"meteo/internal/services/openmeteo"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/zsefvlol/timezonemapper"
)

const (
	MeteoblueName = "meteoblue"
	OpenmeteoName = "openmeteo"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "meteo",
	Short: "CLI app for weather prediction",
	Long:  `CLI app for weather prediction`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse yaml config to struct
		cfg := config.ReadConfigFile()
		httpClient := &http.Client{}

		// Switch between weather API by flag
		apiFlagValue, _ := cmd.Flags().GetString("api")
		if apiFlagValue == "" {
			apiFlagValue = cfg.DefaultAPI
		}

		var weatherService services.Contract
		if apiFlagValue == MeteoblueName {
			weatherService = meteoblue.NewMeteoblue(httpClient)
		} else if apiFlagValue == OpenmeteoName {
			weatherService = openmeteo.NewOpenmeteo(httpClient)
		} else {
			fmt.Printf("Can't recognize weather API name %s", apiFlagValue)
			os.Exit(1)
		}

		// Update longitude and latitude from flags
		latFlagValue, err := cmd.Flags().GetFloat64("lat")
		if err == nil && latFlagValue != -10_000 {
			cfg.Latitude = latFlagValue
		}
		lonFlagValue, err := cmd.Flags().GetFloat64("lon")
		if err == nil && lonFlagValue != -10_000 {
			cfg.Longitude = lonFlagValue
		}

		// TODO: pass coordinates into display part
		fmt.Printf("Latitude: %f \n", cfg.Latitude)
		fmt.Printf("Longitude: %f \n", cfg.Longitude)

		// Get weather data
		weatherData, err := weatherService.Get(cfg)
		if err != nil {
			fmt.Printf("Error fetching weather data: %v\n", err)
			os.Exit(1)
		}

		// Render table
		timezone := timezonemapper.LatLngToTimezoneString(cfg.Latitude, cfg.Longitude)
		display.DisplayTable(weatherData, timezone)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	_ = config.ReadConfigFile() // pass config deeper
	rootCmd.PersistentFlags().Float64P("lat", "", -10_000, "Forecasting latitude")
	rootCmd.PersistentFlags().Float64P("lon", "", -10_000, "Forecasting longitude")
	rootCmd.PersistentFlags().StringP("api", "w", "", "Weather API client")
}
