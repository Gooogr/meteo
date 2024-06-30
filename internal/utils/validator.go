package utils

import "fmt"

func ValidateLongitude(lng float64) error {
	if lng < -180 || lng > 180 {
		return fmt.Errorf("longitude must be between -180 and 180 degrees")
	}
	return nil
}

func ValidateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return fmt.Errorf("latitude must be between -90 and 90 degrees")
	}
	return nil
}
