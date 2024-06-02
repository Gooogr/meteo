package meteoblue

import (
	"bytes"
	"errors"
	"io"
	"meteo/config"
	"meteo/internal/domain"
	"meteo/internal/services/meteoblue/mocks"
	"net/http"
	"reflect"
	"testing"
)

func Test_meteoblue_fetchMeteoblueData(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockResponse   string
		mockStatusCode int
		mockError      error
		wantErr        bool
		wantData       *domain.MeteoblueWeatherData
	}{
		{
			name:           "successful API call",
			url:            "https://my.meteoblue.com/packages/basic-1h?mock=successful",
			mockResponse:   `{"metadata":{"latitude":52.52,"longitude":13.405},"data_1h":{"time":[20060102150405],"temperature":[10],"precipitation_probability":[20],"pictocode":[100],"windspeed":[5]}}`,
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        false,
			wantData: &domain.MeteoblueWeatherData{
				MeteoblueMetadata: domain.MeteoblueMetadata{
					Latitude:  52.52,
					Longitude: 13.405,
				},

				MeteoblueData1h: domain.MeteoblueData1h{
					Time:                     []int64{20060102150405},
					Temperature:              []float64{10},
					PrecipitationProbability: []float64{20},
					Pictocode:                []int64{100},
					WindSpeed:                []float64{5},
				},
			},
		},
		{
			name:           "API returns non-200 status",
			url:            "https://my.meteoblue.com/packages/basic-1h?mock=error",
			mockResponse:   `Bad Request`,
			mockStatusCode: 400,
			mockError:      nil,
			wantErr:        true,
			wantData:       nil,
		},
		{
			name:           "HTTP client error",
			url:            "https://my.meteoblue.com/packages/basic-1h?mock=networkError",
			mockResponse:   "",
			mockStatusCode: 0,
			mockError:      errors.New("network error"),
			wantErr:        true,
			wantData:       nil,
		},
		{
			name:           "Error reading response body",
			url:            "https://my.meteoblue.com/packages/basic-1h?mock=readError",
			mockResponse:   "",
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        true,
			wantData:       nil,
		},
		{
			name:           "Error unmarshaling body",
			url:            "https://my.meteoblue.com/packages/basic-1h?mock=unmarshalError",
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

			mb := meteoblue{client: client}
			data, err := mb.fetchMeteoblueData(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("get() got = %+v, want %+v", data, tt.wantData)
			}
		})
	}
}

func Test_meteoblue_Get(t *testing.T) {
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
					Latitude:                 0.0,
					Longitude:                0.0,
					MeteoblueAPIKey:          "meteoblue-api-key",
					MeteoblueAPISharedSecret: "meteoblue-shared-secret",
				},
			},
			mockResponse: `{
				"metadata": {
					"latitude": 0.0,
					"longitude": 0.0
				},
				"data_1h": {
					"time": [1609459200, 1609462800],
					"temperature": [1.1, 2.2],
					"precipitation_probability": [0.0, 0.1],
					"pictocode": [1, 7],
					"windspeed": [3.3, 4.4]
				}
			}`,
			mockStatus: http.StatusOK,
			want: &domain.WeatherData{
				Time:                     []int64{1609459200, 1609462800},
				Temperature:              []float64{1.1, 2.2},
				PrecipitationProbability: []float64{0.0, 0.1},
				WeatherState:             []string{"Clear, cloudless sky", "Partly cloudy"},
				WindSpeed:                []float64{3.3, 4.4},
			},
			wantErr: false,
		},
		{
			name: "API call returns error",
			args: args{
				cfg: &config.Config{
					Latitude:                 0.0,
					Longitude:                0.0,
					MeteoblueAPIKey:          "meteoblue-api-key",
					MeteoblueAPISharedSecret: "meteoblue-shared-secret",
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
					Latitude:                 100.0, // Invalid latitude
					Longitude:                200.0, // Invalid longitude
					MeteoblueAPIKey:          "meteoblue-api-key",
					MeteoblueAPISharedSecret: "meteoblue-shared-secret",
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

			mb := NewMeteoblue(mockHttpClient)

			got, err := mb.Get(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("meteoblue.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("meteoblue.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createURL(t *testing.T) {
	type args struct {
		lat          float64
		lng          float64
		apiKey       string
		sharedSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid coordinates",
			args: args{
				lat:          37.7749,
				lng:          -122.4194,
				apiKey:       "testApiKey",
				sharedSecret: "testSecret",
			},
			want:    "https://my.meteoblue.com/packages/basic-1h?lat=37.774900&lon=-122.419400&apikey=testApiKey&expire=1924948800&forecast_days=3&temperature=C&timeformat=timestamp_utc&sig=ce1763b9edd8fc1e68ec7af70b81bf9ec8b1679b795eb189d88ec270ed22716a",
			wantErr: false,
		},
		{
			name: "Invalid coordinates",
			args: args{
				lat:          100.0, // Invalid latitude
				lng:          200.0, // Invalid longitude
				apiKey:       "testApiKey",
				sharedSecret: "testSecret",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Empty API key and secret",
			args: args{
				lat:          37.7749,
				lng:          -122.4194,
				apiKey:       "",
				sharedSecret: "",
			},
			want:    "https://my.meteoblue.com/packages/basic-1h?lat=37.774900&lon=-122.419400&apikey=&expire=1924948800&forecast_days=3&temperature=C&timeformat=timestamp_utc&sig=dcf9f0a021f16c291f89bb0f3c9e4e561900e967b3d30fb159179f603b806eb8",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createURL(tt.args.lat, tt.args.lng, tt.args.apiKey, tt.args.sharedSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("createURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateSignature(t *testing.T) {
	type args struct {
		data   string
		secret string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test case 1",
			args: args{
				data:   "hello",
				secret: "world",
			},
			want: "3cfa76ef14937c1c0ea519f8fc057a80fcd04a7420f8e8bcd0a7567c272e007b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateSignature(tt.args.data, tt.args.secret); got != tt.want {
				t.Errorf("generateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
