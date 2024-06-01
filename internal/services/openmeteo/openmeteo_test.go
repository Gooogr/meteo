package openmeteo

import (
	"bytes"
	"errors"
	"io"
	"meteo/internal/domain"
	"meteo/internal/services/openmeteo/mocks"
	"net/http"
	"reflect"
	"testing"
	"time"
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
			want:    "https://api.open-meteo.com/v1/forecast?latitude=0.000000&longitude=0.000000&timezone=Africa/Sao_Tome&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3",
			wantErr: false,
		},
		{
			name:    "Negative coordinates",
			args:    args{lat: -45.0, lng: -90.0},
			want:    "https://api.open-meteo.com/v1/forecast?latitude=-45.000000&longitude=-90.000000&timezone=America/Santiago&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3",
			wantErr: false,
		},
		{
			name:    "Floating point precision",
			args:    args{lat: 37.7749, lng: -122.4194},
			want:    "https://api.open-meteo.com/v1/forecast?latitude=37.774900&longitude=-122.419400&timezone=America/Los_Angeles&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3",
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

func Test_openmeteo_get(t *testing.T) {
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
			mockResponse:   `{"latitude":52.52,"longitude":13.405,"hourly":{"time":["2023-01-01T00:00"],"temperature_2m":[10],"precipitation_probability":[20],"weathercode":[100],"windspeed_10m":[5]}}`,
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        false,
			wantData: &domain.OpenmeteoWeatherData{
				Latitude:  52.52,
				Longitude: 13.405,
				Hourly: domain.HourlyData{
					Time:                     domain.TimeSlice{time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)},
					Temperature2m:            []float64{10},
					PrecipitationProbability: []float64{20},
					WeatherCode:              []int{100},
					WindSpeed10m:             []float64{5},
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
			data, err := giver.get(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("get() got = %+v, want %+v", data, tt.wantData)
			}
		})
	}
}

//func Test_Example(t *testing.T) {
//	type args struct {
//		lat float64
//		lng float64
//	}
//
//	// Init vars.
//	arguments := args{lat: 0.0, lng: 0.0}
//
//	// Mocks.
//	mockHttpClient := mocks.MockClient{}
//
//	// Init service.
//	svc := NewOpenmeteo(mockHttpClient)
//
//	// Test.
//	res, err := svc.Get(arguments.lat, arguments.lng)
//	// asserts
//}
