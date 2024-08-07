package event

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Event struct {
	ID              EventUUID   `json:"id"`
	Name            string      `json:"name" validate:"required,min=1"`
	StartTime       EventTime   `json:"start_time"`
	EndTime         EventTime   `json:"end_time"`
	IsSpecificDates bool        `json:"is_specific_dates"`
	Dates           []EventDate `json:"dates" validate:"required,min=1"`
	// TimeZone        string    `json:"time_zone" validate:"required"`
}

func RegisterEventServiceValidators(validate *validator.Validate) {
	validate.RegisterStructValidation(func(sl validator.StructLevel) {
		event := sl.Current().Interface().(Event)
		ValidateStartAndEndTime(event.StartTime, event.EndTime, sl)
	}, Event{})
}

func isEndTimeMidnight(time time.Time) bool {
	return time.Hour() == 0 && time.Minute() == 0 && time.Second() == 0
}

func ValidateStartAndEndTime(startTime EventTime, endTime EventTime, sl validator.StructLevel) {
	if startTime.IsZero() {
		sl.ReportError(startTime, "start_time", "StartTime", "required", "")
	}
	if endTime.IsZero() {
		sl.ReportError(endTime, "end_time", "EndTime", "required", "")
	}

	// if end time is midnight then assume it is the next day
	if !isEndTimeMidnight(endTime.Time) && !startTime.Before(endTime.Time) {
		sl.ReportError(endTime, "end_time", "EndTime", "end_time_after_start_time", "")
	}
}
