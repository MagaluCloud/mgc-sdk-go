package utils

import (
	"testing"
	"time"
)

func TestLocalDateTimeWithoutZone_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid time",
			data:    []byte(`"2023-01-02T12:34:56.000000"`),
			want:    time.Date(2023, time.January, 2, 12, 34, 56, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format with space",
			data:    []byte(`"2023-01-02 12:34:56.000000"`),
			wantErr: true,
		},
		{
			name:    "invalid month",
			data:    []byte(`"2023-13-02T12:34:56.000000"`),
			wantErr: true,
		},
		{
			name:    "empty string",
			data:    []byte(`""`),
			wantErr: true,
		},
		{
			name:    "null input",
			data:    []byte(`null`),
			wantErr: true,
		},
		{
			name:    "malformed time",
			data:    []byte(`"2023-01-02T12:34"`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ct LocalDateTimeWithoutZone
			err := ct.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !time.Time(ct).Equal(tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", time.Time(ct), tt.want)
			}
		})
	}
}

func TestLocalDateTimeWithoutZone_MarshalJSON(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	tests := []struct {
		name    string
		ct      LocalDateTimeWithoutZone
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid UTC time",
			ct:      LocalDateTimeWithoutZone(time.Date(2023, time.January, 2, 12, 34, 56, 0, time.UTC)),
			want:    []byte(`"2023-01-02T12:34:56.000000"`),
			wantErr: false,
		},
		{
			name:    "valid non-UTC time",
			ct:      LocalDateTimeWithoutZone(time.Date(2023, time.January, 2, 12, 34, 56, 0, loc)),
			want:    []byte(`"2023-01-02T12:34:56.000000"`),
			wantErr: false,
		},
		{
			name:    "zero time",
			ct:      LocalDateTimeWithoutZone(time.Time{}),
			want:    []byte(`"0001-01-01T00:00:00.000000"`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ct.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("MarshalJSON() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLocalDateTimeWithoutZone_String(t *testing.T) {
	tests := []struct {
		name string
		ct   LocalDateTimeWithoutZone
		want string
	}{
		{
			name: "valid time",
			ct:   LocalDateTimeWithoutZone(time.Date(2023, time.January, 2, 12, 34, 56, 0, time.UTC)),
			want: "2023-01-02T12:34:56.000000",
		},
		{
			name: "zero time",
			ct:   LocalDateTimeWithoutZone(time.Time{}),
			want: "0001-01-01T00:00:00.000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
