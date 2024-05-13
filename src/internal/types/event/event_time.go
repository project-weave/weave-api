package event

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type EventTime struct {
	time.Time
}

func (et EventTime) String() string {
	return et.Format(time.TimeOnly)
}

func (et *EventTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.TimeOnly, s)
	if err != nil {
		return err
	}
	et.Time = t
	return nil
}

func (et EventTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", et.String())), nil
}

func (d EventTime) Value() (driver.Value, error) {
	return d.Time, nil
}

func (et *EventTime) Scan(value interface{}) error {
	timeStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot convert %v to string", value)
	}
	t, err := time.Parse(time.TimeOnly, timeStr)
	if err != nil {
		return err
	}
	et.Time = t
	return nil
}
