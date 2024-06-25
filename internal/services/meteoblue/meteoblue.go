package meteoblue

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"meteo/config"
	"meteo/internal/domain"
	"meteo/internal/dto"
	"meteo/internal/services"
)

var meteobluePictocodes = map[int64]string{
	1:  "Clear, cloudless sky",
	2:  "Clear, few cirrus",
	3:  "Clear with cirrus",
	4:  "Clear with few low clouds",
	5:  "Clear with few low clouds and few cirrus",
	6:  "Clear with few low clouds and cirrus",
	7:  "Partly cloudy",
	8:  "Partly cloudy and few cirrus",
	9:  "Partly cloudy and cirrus",
	10: "Mixed with some thunderstorm clouds possible",
	11: "Mixed with few cirrus with some thunderstorm clouds possible",
	12: "Mixed with cirrus and some thunderstorm clouds possible",
	13: "Clear but hazy",
	14: "Clear but hazy with few cirrus",
	15: "Clear but hazy with cirrus",
	16: "Fog/low stratus clouds",
	17: "Fog/low stratus clouds with few cirrus",
	18: "Fog/low stratus clouds with cirrus",
	19: "Mostly cloudy",
	20: "Mostly cloudy and few cirrus",
	21: "Mostly cloudy and cirrus",
	22: "Overcast",
	23: "Overcast with rain",
	24: "Overcast with snow",
	25: "Overcast with heavy rain",
	26: "Overcast with heavy snow",
	27: "Rain, thunderstorms likely",
	28: "Light rain, thunderstorms likely",
	29: "Storm with heavy snow",
	30: "Heavy rain, thunderstorms likely",
	31: "Mixed with showers",
	32: "Mixed with snow showers",
	33: "Overcast with light rain",
	34: "Overcast with light snow",
	35: "Overcast with mixture of snow and rain",
}

type meteoblue struct {
	client httpClient
}

func NewMeteoblue(client httpClient) services.Contract {
	return &meteoblue{
		client: client,
	}
}

func (mb *meteoblue) fetchMeteoblueData(url string) (*domain.MeteoblueWeatherData, error) {
	weatherDto := dto.MeteoblueWeatherData{}

	// Get data from meteoblue.
	resp, err := mb.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &weatherDto)
	if err != nil {
		return nil, err
	}

	// Convert dto to domain.
	data := &domain.MeteoblueWeatherData{
		MeteoblueMetadata: domain.MeteoblueMetadata{
			Latitude:  weatherDto.MeteoblueMetadata.Latitude,
			Longitude: weatherDto.MeteoblueMetadata.Longitude,
		},
		MeteoblueData1h: domain.MeteoblueData1h{
			Time:                     weatherDto.MeteoblueData1h.Time,
			Temperature:              weatherDto.MeteoblueData1h.Temperature,
			PrecipitationProbability: weatherDto.MeteoblueData1h.PrecipitationProbability,
			Pictocode:                weatherDto.MeteoblueData1h.Pictocode,
			WindSpeed:                weatherDto.MeteoblueData1h.WindSpeed,
		},
	}

	return data, nil
}

func (mb *meteoblue) Get(cfg *config.Config) (*domain.WeatherData, error) {
	url, err := createURL(cfg.Latitude, cfg.Longitude, cfg.MeteoblueAPIKey, cfg.MeteoblueAPISharedSecret)
	if err != nil {
		return nil, err
	}

	data, err := mb.fetchMeteoblueData(url)
	if err != nil {
		return nil, err
	}

	weatherState := make([]string, len(data.MeteoblueData1h.Pictocode))
	for i, code := range data.MeteoblueData1h.Pictocode {
		state := ""

		wc, ok := meteobluePictocodes[code]
		if ok {
			state = wc
		}

		weatherState[i] = state
	}
	return &domain.WeatherData{
		Time:                     data.MeteoblueData1h.Time,
		Temperature:              data.MeteoblueData1h.Temperature,
		PrecipitationProbability: data.MeteoblueData1h.PrecipitationProbability,
		WeatherState:             weatherState,
		WindSpeed:                data.MeteoblueData1h.WindSpeed,
	}, nil
}

func createURL(lat, lng float64, apiKey, sharedSecret string) (string, error) {
	if err := services.ValidateLongitude(lng); err != nil {
		return "", err
	}
	if err := services.ValidateLatitude(lat); err != nil {
		return "", err
	}

	query := fmt.Sprintf(
		"/packages/basic-1h?lat=%.6f&lon=%.6f&apikey=%s&expire=1924948800&forecast_days=3&temperature=C&timeformat=timestamp_utc",
		lat, lng, apiKey,
	)

	sig := generateSignature(query, sharedSecret)
	url := fmt.Sprintf("https://my.meteoblue.com%s&sig=%s", query, sig)

	return url, nil
}

func generateSignature(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
