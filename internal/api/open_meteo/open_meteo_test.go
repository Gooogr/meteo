package openmeteo

import (
	"bytes"
	"errors"
	"io"
	"meteo/internal/api/open_meteo/mocks"
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
		name string
		args args
		want string
	}{
		{
			name: "Null Island",
			args: args{lat: 0.0, lng: 0.0},
			want: "https://api.open-meteo.com/v1/forecast?latitude=0.000000&longitude=0.000000&timezone=Africa/Sao_Tome&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3",
		},
		{
			name: "Negative coordinates",
			args: args{lat: -45.0, lng: -90.0},
			want: "https://api.open-meteo.com/v1/forecast?latitude=-45.000000&longitude=-90.000000&timezone=America/Santiago&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3",
		},
		{
			name: "Floating point precision",
			args: args{lat: 37.7749, lng: -122.4194},
			want: "https://api.open-meteo.com/v1/forecast?latitude=37.774900&longitude=-122.419400&timezone=America/Los_Angeles&hourly=temperature_2m,precipitation_probability,weathercode,windspeed_10m&forecast_days=3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createURL(tt.args.lat, tt.args.lng); got != tt.want {
				t.Errorf("createURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeSlice_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    TimeSlice
		wantErr bool
	}{
		{
			name:    "Valid time format",
			args:    args{data: []byte(`["2023-01-01T12:00", "2023-01-01T13:00"]`)},
			want:    TimeSlice{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC), time.Date(2023, 1, 1, 13, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name:    "Invalid time format",
			args:    args{data: []byte(`["2023-01-01 12:00", "2023-01-01T13:00"]`)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty time slice",
			args:    args{data: []byte(`[]`)},
			want:    TimeSlice{},
			wantErr: false,
		},
		{
			name:    "Single valid time",
			args:    args{data: []byte(`["2023-01-01T12:00"]`)},
			want:    TimeSlice{time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			args:    args{data: []byte(`[2023-01-01T12:00]`)},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts TimeSlice
			err := ts.UnmarshalJSON(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimeSlice.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(ts, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", ts, tt.want)
			}
		})
	}
}

func Test_openMeteoGiver_get(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockResponse   string
		mockStatusCode int
		mockError      error
		wantErr        bool
	}{
		{
			name:           "successful API call",
			url:            "https://api.open-meteo.com/v1/forecast?mock=successful",
			mockResponse:   `{"success":true,"data":"valid weather data"}`,
			mockStatusCode: 200,
			mockError:      nil,
			wantErr:        false,
		},
		{
			name:           "API returns non-200 status",
			url:            "https://api.open-meteo.com/v1/forecast?mock=error",
			mockResponse:   `Bad Request`,
			mockStatusCode: 400,
			mockError:      nil,
			wantErr:        true,
		},
		{
			name:           "HTTP client error",
			url:            "https://api.open-meteo.com/v1/forecast?mock=networkError",
			mockResponse:   "",
			mockStatusCode: 0,
			mockError:      errors.New("network error"),
			wantErr:        true,
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

			giver := openMeteoGiver{client: client}
			responseBytes, err := giver.get(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("giver.get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				expectedResponse := tt.mockResponse
				if string(responseBytes) != expectedResponse {
					t.Errorf("giver.get() expected %s, got %s", expectedResponse, string(responseBytes))
				}
			}
		})
	}
}
