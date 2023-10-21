package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/project-weave/weave-api/src/internal/validator"
)

type Event struct {
	ID                       uuid.UUID      `json:"event_id"`
	Name                     string         `json:"name"`
	StartDate                time.Time      `json:"start_date"`
	EndDate                  time.Time      `json:"end_date"`
	CreatedAt                time.Time      `json:"-"`
	CreatedBy                uuid.UUID      `json:"created_by"`
	IncludesTimeAvailability bool           `json:"boolean"`
	Availability             []Availability `json:"availability"`
}

type Availability struct {
	UserID        uuid.UUID `json:"user_id"`
	StartDateTime time.Time `json:"start_datetime"`
	EndDateTime   time.Time `json:"end_datetime"`
}

type EventService interface {
	AddEvent(event *Event) error
	GetEvent(eventID uuid.UUID) (*Event, error)
	GetUserEventAvailability(eventID uuid.UUID, userID uuid.UUID) ([]Availability, error)
	UpdateUserEventAvailability(eventID uuid.UUID, userID uuid.UUID, availabilities []Availability) error
}

func ValidateEvent(e *Event) map[string]string {
	eventValidator := validator.Validator{}
	currentDate := time.Now().Truncate(24 * time.Hour)
	return eventValidator.Validate(
		validator.Condition{
			Ok:        e.Name != "",
			FieldName: "name",
			Message:   "event name should be set",
		},
		validator.Condition{
			Ok:        e.EndDate.After(e.StartDate),
			FieldName: "end date",
			Message:   "end date should be after start date",
		},
		validator.Condition{
			Ok:        e.StartDate.Equal(currentDate) || e.StartDate.After(currentDate),
			FieldName: "start date",
			Message:   "start date must be equal to, or after today's date",
		},
	)
}
