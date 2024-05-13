package event

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type EventDateTime struct {
	time.Time
}

func (edt EventDateTime) String() string {
	return edt.Format(time.DateTime)
}

func (edt *EventDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	edt.Time = t
	return nil
}

func (edt EventDateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", edt.String())), nil
}

func (edt EventDateTime) Value() (driver.Value, error) {
	return edt.Time, nil
}

func (edt *EventDateTime) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot convert %v to time.Time", value)
	}
	edt.Time = t
	return nil
}
