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
	"meteo/internal/utils"
)

type meteoblue struct {
	client httpClient
	cfg    config.Config
}

func NewMeteoblue(
	client httpClient,
	config config.Config,
) services.Contract {
	return &meteoblue{
		client: client,
		cfg:    config,
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

func (mb *meteoblue) Get(latitude, longitude float64) (*domain.WeatherData, error) {
	cfgMeteoblue := mb.cfg.MeteoblueConfig()

	url, err := createURL(
		latitude,
		longitude,
		cfgMeteoblue.APIKey,
		cfgMeteoblue.APISharedSecret,
	)
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
	if err := utils.ValidateLongitude(lng); err != nil {
		return "", err
	}
	if err := utils.ValidateLatitude(lat); err != nil {
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
