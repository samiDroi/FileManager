package utils

import (
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + ct.Time.Format("2006-01-02 15:04:05") + `"`), nil
}

func ToCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}
