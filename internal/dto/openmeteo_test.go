package dto

import (
	"reflect"
	"testing"
	"time"
)

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
