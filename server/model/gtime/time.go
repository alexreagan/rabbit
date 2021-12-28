package gtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const timeFormatLocal = "2006-01-02T15:04:05+08:00"
const zTimeFormat = "2006-01-02T15:04:05Z"

type GTime struct {
	time.Time
}

func NewGTime(v time.Time) GTime {
	return GTime{Time: v}
}

func ZeroTime() GTime {
	return GTime{Time: time.Time{}}
}

func Now() GTime {
	return GTime{Time: time.Now()}
}

func (t GTime) MarshalJSON() ([]byte, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return []byte("\"\""), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(timeFormat))
	return []byte(formatted), nil
}

func (t *GTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "" {
		*t = ZeroTime()
		return nil
	}

	tm, err := time.ParseInLocation(timeFormat, string(data[1:len(data)-1]), time.Local)
	if err != nil {
		tm, err = time.ParseInLocation(timeFormatLocal, string(data[1:len(data)-1]), time.Local)
		if err != nil {
			tm, err = time.ParseInLocation(zTimeFormat, string(data[1:len(data)-1]), time.Local)
			*t = NewGTime(tm)
			return err
		} else {
			*t = NewGTime(tm)
			return err
		}
	} else {
		*t = NewGTime(tm)
		return err
	}
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
