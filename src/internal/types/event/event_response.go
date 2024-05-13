package event

import "github.com/google/uuid"

type EventResponse struct {
	EventID       EventUUID       `json:"-" param:"event_id"`
	UserID        uuid.UUID       `json:"user_id"`
	Alias         string          `json:"alias" validate:"required"`
	Availabiliies []EventDateTime `json:"availabilities" validate:"required"`
}
