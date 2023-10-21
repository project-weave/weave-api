package postgres

import (
	"github.com/google/uuid"
	"github.com/project-weave/weave-api/src/internal/domain"
)

type EventService struct {
	db *DB
}

var service domain.EventService
var _ = service.(*EventService)

func NewEventService(db *DB) *EventService {
	return &EventService{
		db: db,
	}
}

func (e *EventService) AddEvent(event *domain.Event) error {
	return nil
}

func (e *EventService) GetEvent(eventId uuid.UUID) (*domain.Event, error) {
	return nil, nil
}

func (e *EventService) GetUserEventAvailability(eventID uuid.UUID, userID uuid.UUID) ([]domain.Availability, error) {
	return nil, nil
}

func (e *EventService) UpdateUserEventAvailability(eventID uuid.UUID, userID uuid.UUID, availabilities []domain.Availability) error {
	return nil
}
