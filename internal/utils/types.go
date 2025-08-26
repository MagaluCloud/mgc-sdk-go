package utils

import (
	"strings"
	"time"
)

// LocalDateTimeWithoutZone represents a local date-time without timezone information.
//
// It is a thin wrapper around time.Time that customizes JSON marshaling and
// unmarshaling to use the layout defined by
// `LocalDateTimeWithoutZoneLayout`.
type LocalDateTimeWithoutZone time.Time

// LocalDateTimeWithoutZoneLayout is the default format used when marshaling
// and formatting `LocalDateTimeWithoutZone` values. It includes microseconds.
const LocalDateTimeWithoutZoneLayout = "2006-01-02T15:04:05.000000"

// UnmarshalJSON implements the json.Unmarshaler interface for
// LocalDateTimeWithoutZone.
//
// It attempts to parse the incoming JSON string using multiple layouts:
// - the microsecond layout defined by LocalDateTimeWithoutZoneLayout
// - a seconds-only layout without microseconds
// - time.RFC3339 as a fallback
//
// This allows it to accept timestamps both with and without microseconds.
func (ct *LocalDateTimeWithoutZone) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	layouts := []string{
		"2006-01-02T15:04:05.000000", // With microseconds
		"2006-01-02T15:04:05",        // Without microseconds
		time.RFC3339,
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, s)
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	*ct = LocalDateTimeWithoutZone(t)
	return nil
}

func (ct LocalDateTimeWithoutZone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ct).Format(LocalDateTimeWithoutZoneLayout) + `"`), nil
}

func (ct LocalDateTimeWithoutZone) String() string {
	return time.Time(ct).Format(LocalDateTimeWithoutZoneLayout)
}
