package openmeteo

import (
	"bytes"
	"errors"
	"io"
	"meteo/config"
	"meteo/internal/domain"
	"meteo/internal/services/openmeteo/mocks"
	"net/http"
	"reflect"
	"testing"
)

func Test_createURL(t *testing.T) {
	type args struct {
		lat float64
		lng float64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Null Island",
			args:    args{lat: 0.0, lng: 0.0},
			want:    "https://api.open-meteo.com/v1/forecast?latitude=0.000000&longitude=0.000000&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3&timeformat=unixtime",
			wantErr: false,
		},
		{
			name:    "Negative coordinates",
			args:    args{lat: -45.0, lng: -90.0},
			want:    "https://api.open-meteo.com/v1/forecast?latitude=-45.000000&longitude=-90.000000&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3&timeformat=unixtime",
			wantErr: false,
		},
		{
			name:    "Floating point precision",
			args:    args{lat: 37.7749, lng: -122.4194},
			want:    "https://api.open-meteo.com/v1/forecast?latitude=37.774900&longitude=-122.419400&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3&timeformat=unixtime",
			wantErr: false,
		},
		{
			name:    "Incorrect latitude",
			args:    args{lat: -95.0, lng: 0.0},
			wantErr: true,
		},
		{
			name:    "Incorrect longitude",
			args:    args{lat: 0.0, lng: -185},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createURL(tt.args.lat, tt.args.lng)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("createURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_openmeteo_fetchOpenmeteoData(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockResponse   string
		mockStatusCode int
		mockError      error
		wantErr        bool
		wantData       *domain.OpenmeteoWeatherData
	}{
		{
			name:           "successful API call",
			url:            "https://api.open-meteo.com/v1/forecast?mock=successful",
			mockResponse:   `{"latitude":52.52,"longitude":13.405,"hourly":{"time":[20060102150405],"temperature_2m":[10],"precipitation_probability":[20],"weathercode":[100],"windspeed_10m":[5]}}`,
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        false,
			wantData: &domain.OpenmeteoWeatherData{
				Latitude:  52.52,
				Longitude: 13.405,
				Hourly: domain.OpenmeteoHourlyData{
					Time:                     []int64{20060102150405},
					Temperature:              []float64{10},
					PrecipitationProbability: []float64{20},
					WeatherCode:              []int64{100},
					WindSpeed:                []float64{5},
				},
			},
		},
		{
			name:           "API returns non-200 status",
			url:            "https://api.open-meteo.com/v1/forecast?mock=error",
			mockResponse:   `Bad Request`,
			mockStatusCode: 400,
			mockError:      nil,
			wantErr:        true,
			wantData:       nil,
		},
		{
			name:           "HTTP client error",
			url:            "https://api.open-meteo.com/v1/forecast?mock=networkError",
			mockResponse:   "",
			mockStatusCode: 0,
			mockError:      errors.New("network error"),
			wantErr:        true,
			wantData:       nil,
		},
		{
			name:           "Error reading response body",
			url:            "https://api.open-meteo.com/v1/forecast?mock=readError",
			mockResponse:   "",
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        true,
			wantData:       nil,
		},
		{
			name:           "Error unmarshaling body",
			url:            "https://api.open-meteo.com/v1/forecast?mock=unmarshalError",
			mockResponse:   `{"broken json": {`,
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        true,
			wantData:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mocks.MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					r := io.NopCloser(bytes.NewReader([]byte(tt.mockResponse)))
					return &http.Response{
						StatusCode: tt.mockStatusCode,
						Body:       r,
					}, tt.mockError
				},
			}

			giver := openmeteo{client: client}
			data, err := giver.fetchOpenmeteoData(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("get() got = %+v, want %+v", data, tt.wantData)
			}
		})
	}
}

func Test_openmeteo_Get(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name         string
		args         args
		mockResponse string
		mockStatus   int
		want         *domain.WeatherData
		wantErr      bool
	}{
		{
			name: "successful API call",
			args: args{
				cfg: &config.Config{
					Latitude:  0.0,
					Longitude: 0.0,
				},
			},
			mockResponse: `{
				"latitude": 0.0,
				"longitude": 0.0,
				"hourly": {
					"time": [1609459200, 1609462800],
					"temperature_2m": [1.1, 2.2],
					"precipitation_probability": [0.0, 0.1],
					"weathercode": [0, 1],
					"windspeed_10m": [3.3, 4.4]
				}
			}`,
			mockStatus: http.StatusOK,
			want: &domain.WeatherData{
				Time:                     []int64{1609459200, 1609462800},
				Temperature:              []float64{1.1, 2.2},
				PrecipitationProbability: []float64{0.0, 0.1},
				WeatherState:             []string{"Clear sky", "Mainly clear"},
				WindSpeed:                []float64{3.3, 4.4},
			},
			wantErr: false,
		},
		{
			name: "API call returns error",
			args: args{
				cfg: &config.Config{
					Latitude:  0.0,
					Longitude: 0.0,
				},
			},
			mockResponse: `{"error": "invalid request"}`,
			mockStatus:   http.StatusBadRequest,
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Invalid latitude and longitude",
			args: args{
				cfg: &config.Config{
					Latitude:  100.0, // Invalid latitude
					Longitude: 200.0, // Invalid longitude
				},
			},
			mockResponse: "",
			mockStatus:   0,
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHttpClient := &mocks.MockClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.mockStatus,
						Body:       io.NopCloser(bytes.NewReader([]byte(tt.mockResponse))),
					}, nil
				},
			}

			om := NewOpenmeteo(mockHttpClient)

			got, err := om.Get(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("openmeteo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("openmeteo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
