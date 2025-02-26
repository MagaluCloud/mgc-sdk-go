package utils

import (
	"strings"
	"time"
)

type LocalDateTimeWithoutZone time.Time

const LocalDateTimeWithoutZoneLayout = "2006-01-02T15:04:05.000000"

func (ct *LocalDateTimeWithoutZone) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse(LocalDateTimeWithoutZoneLayout, s)
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
