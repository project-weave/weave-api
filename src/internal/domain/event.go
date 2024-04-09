package domain

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type EventService interface {
	AddEvent(ctx context.Context, event *Event) (uuid.UUID, error)
	GetEvent(ctx context.Context, eventID uuid.UUID) (*Event, []EventResponse, error)
	UpsertUserEventAvailability(ctx context.Context, eventResponse EventResponse) error
}

type Event struct {
	ID              uuid.UUID   `json:"id"`
	Name            string      `json:"name" validate:"required,min=1"`
	StartTime       EventTime   `json:"start_time" validate:"required"`
	EndTime         EventTime   `json:"end_time" validate:"required"`
	IsSpecificDates bool        `json:"is_specific_dates" validate:"required"`
	EventDates      []EventDate `json:"event_dates" validate:"required,min=1"`
	// TimeZone        string    `json:"time_zone" validate:"required"`
}

type EventResponse struct {
	EventID       uuid.UUID       `json:"-" param:"event_id"`
	UserID        uuid.UUID       `json:"user_id"`
	Alias         string          `json:"alias" validate:"required"`
	Availabiliies []EventDateTime `json:"availabilities" validate:"required"`
}

func RegisterEventServiceValidators(validate *validator.Validate) {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		event := sl.Current().Interface().(Event)
		ValidateStartAndEndTime(event.StartTime, event.EndTime, sl)
	}, Event{})
}

func ValidateStartAndEndTime(startTime EventTime, endTime EventTime, sl validator.StructLevel) {
	if startTime.After(endTime.Time) {
		sl.ReportError(endTime, "end_time", "EndTime", "end_time_after_start_time", "")
	}
}

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
