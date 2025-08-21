package utils

import (
	"strings"
	"time"
)

type LocalDateTimeWithoutZone time.Time

const LocalDateTimeWithoutZoneLayout = "2006-01-02T15:04:05.000000"

func (ct *LocalDateTimeWithoutZone) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)

	// Try multiple timestamp formats
	layouts := []string{
		"2006-01-02T15:04:05.000000", // With microseconds
		"2006-01-02T15:04:05",        // Without microseconds
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
