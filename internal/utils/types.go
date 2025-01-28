package utils

import (
	"strings"
	"time"
)

type LocalDateTimeWithoutZone time.Time

func (ct *LocalDateTimeWithoutZone) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}
	*ct = LocalDateTimeWithoutZone(t)
	return nil
}

func (ct LocalDateTimeWithoutZone) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(ct).Format("2006-01-02T15:04:05") + `"`), nil
}

func (ct LocalDateTimeWithoutZone) String() string {
	return time.Time(ct).Format("2006-01-02T15:04:05")
}
