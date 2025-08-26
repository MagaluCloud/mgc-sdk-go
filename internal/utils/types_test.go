package utils

import (
	"testing"
	"time"
)

func TestLocalDateTimeWithoutZone_UnmarshalJSON_NoMicroseconds(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`"2023-01-02T12:34:56"`)
	err := ct.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
	}
	want := time.Date(2023, time.January, 2, 12, 34, 56, 0, time.UTC)
	if !time.Time(ct).Equal(want) {
		t.Errorf("UnmarshalJSON() got = %v, want %v", time.Time(ct), want)
	}
}

func TestLocalDateTimeWithoutZone_UnmarshalJSON_RFC3339Z(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`"2023-01-02T12:34:56Z"`)
	err := ct.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
	}
	want := time.Date(2023, time.January, 2, 12, 34, 56, 0, time.UTC)
	if !time.Time(ct).Equal(want) {
		t.Errorf("UnmarshalJSON() got = %v, want %v", time.Time(ct), want)
	}
}

func TestLocalDateTimeWithoutZone_UnmarshalJSON_RFC3339WithOffset(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`"2023-01-02T12:34:56-05:00"`)
	err := ct.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
	}
	want := time.Date(2023, time.January, 2, 17, 34, 56, 0, time.UTC)
	if !time.Time(ct).Equal(want) {
		t.Errorf("UnmarshalJSON() got = %v, want %v", time.Time(ct), want)
	}
}

func TestLocalDateTimeWithoutZone_UnmarshalJSON_RFC3339WithFractionalSeconds_Supported(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`"2023-01-02T12:34:56.789Z"`)
	err := ct.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
	}
	want := time.Date(2023, time.January, 2, 12, 34, 56, 789000000, time.UTC)
	if !time.Time(ct).Equal(want) {
		t.Errorf("UnmarshalJSON() got = %v, want %v", time.Time(ct), want)
	}
}

func TestLocalDateTimeWithoutZone_UnmarshalJSON_WhitespaceInsideQuotes_ShouldError(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`" 2023-01-02T12:34:56.000000 "`)
	err := ct.UnmarshalJSON(data)
	if err == nil {
		t.Fatalf("UnmarshalJSON() expected error with leading/trailing whitespace inside quotes, got none")
	}
}

func TestLocalDateTimeWithoutZone_UnmarshalJSON_DateOnly_ShouldError(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`"2023-01-02"`)
	err := ct.UnmarshalJSON(data)
	if err == nil {
		t.Fatalf("UnmarshalJSON() expected error for date-only string, got none")
	}
}

func TestLocalDateTimeWithoutZone_MarshalJSON_TruncatesToMicroseconds(t *testing.T) {
	ct := LocalDateTimeWithoutZone(time.Date(2023, time.January, 2, 12, 34, 56, 123456789, time.UTC))
	got, err := ct.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() unexpected error: %v", err)
	}
	want := []byte(`"2023-01-02T12:34:56.123456"`)
	if string(got) != string(want) {
		t.Errorf("MarshalJSON() got = %s, want %s", got, want)
	}
}

func TestLocalDateTimeWithoutZone_String_WithNonUTCLocation(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	ct := LocalDateTimeWithoutZone(time.Date(2023, time.January, 2, 12, 34, 56, 0, loc))
	got := ct.String()
	want := "2023-01-02T12:34:56.000000"
	if got != want {
		t.Errorf("String() = %v, want %v", got, want)
	}
}

func TestLocalDateTimeWithoutZone_RoundTrip_FromRFC3339WithOffset(t *testing.T) {
	var ct LocalDateTimeWithoutZone
	data := []byte(`"2023-01-02T12:34:56-05:00"`)
	if err := ct.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON() unexpected error: %v", err)
	}
	got, err := ct.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() unexpected error: %v", err)
	}
	want := []byte(`"2023-01-02T12:34:56.000000"`)
	if string(got) != string(want) {
		t.Errorf("RoundTrip MarshalJSON() got = %s, want %s", got, want)
	}
}
