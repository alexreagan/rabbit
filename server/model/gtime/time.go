package gtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type GTime struct {
	time.Time
}

func NewGTime(v time.Time) GTime {
	return GTime{Time: v}
}

func (t GTime) MarshalJSON() ([]byte, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return []byte("\"\""), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(timeFormat))
	return []byte(formatted), nil
}

func (t GTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *GTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = GTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
