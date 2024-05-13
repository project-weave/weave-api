package event

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type EventDate struct {
	time.Time
}

func (ed EventDate) String() string {
	return ed.Format(time.DateOnly)
}

func (ed *EventDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	ed.Time = t
	return nil
}

func (ed EventDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ed.String())), nil
}

func (ed EventDate) Value() (driver.Value, error) {
	return ed.Time, nil
}

func (ed *EventDate) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot convert %v to time.Time", value)
	}
	ed.Time = t
	return nil
}
