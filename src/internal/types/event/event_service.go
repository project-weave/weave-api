package event

import "context"

type EventService interface {
	AddEvent(ctx context.Context, event *Event) (EventUUID, error)
	GetEvent(ctx context.Context, eventID EventUUID) (*Event, []EventResponse, error)
	UpsertUserEventAvailability(ctx context.Context, eventResponse EventResponse) error
}
