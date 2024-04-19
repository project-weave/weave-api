package postgres

import (
	"context"
	"sort"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/project-weave/weave-api/src/internal/domain"
)

type EventService struct {
	db *DB
}

func NewEventService(db *DB) *EventService {
	return &EventService{
		db: db,
	}
}

func (e *EventService) AddEvent(ctx context.Context, event *domain.Event) (uuid.UUID, error) {
	sql, args, err := sq.Insert("events").
		Columns("name", "is_specific_dates", "start_time", "end_time", "dates").
		Values(event.Name, event.IsSpecificDates, event.StartTime, event.EndTime.Time, event.Dates).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return uuid.UUID{}, err
	}

	row := e.db.pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&event.ID)
	if err != nil {
		return uuid.UUID{}, err
	}

	return event.ID, err
}

func (e *EventService) GetEvent(ctx context.Context, eventID uuid.UUID) (*domain.Event, []domain.EventResponse, error) {
	tx, err := e.db.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	sql, args, err := sq.Select("id", "name", "is_specific_dates", "start_time", "end_time", "dates").
		From("events").
		Where(sq.Eq{"id": eventID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, err
	}

	row := tx.QueryRow(ctx, sql, args...)
	event := domain.Event{}
	err = row.Scan(&event.ID, &event.Name, &event.IsSpecificDates, &event.StartTime, &event.EndTime, &event.Dates)
	if err != nil {
		return nil, nil, err
	}

	sql, args, err = sq.Select("user_id", "event_id", "alias", "availabilities").
		From("event_responses").
		Where(sq.Eq{"event_id": eventID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, err
	}
	eventResponses := []domain.EventResponse{}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		eventResponse := domain.EventResponse{}
		timeRanges := [][]domain.EventDateTime{}

		err = rows.Scan(&eventResponse.UserID, &eventResponse.EventID, &eventResponse.Alias, &timeRanges)
		if err != nil {
			return nil, nil, err
		}

		timeSlots, err := convertTimeRangesToTimeSlots(timeRanges)
		if err != nil {
			return nil, nil, err
		}
		eventResponse.Availabiliies = timeSlots
		eventResponses = append(eventResponses, eventResponse)
	}

	return &event, eventResponses, nil
}

var timeSlotGapMins = 30

func convertTimeRangesToTimeSlots(ranges [][]domain.EventDateTime) ([]domain.EventDateTime, error) {
	var timeSlots []domain.EventDateTime
	for _, timeRange := range ranges {
		startTime := timeRange[0]
		endTime := timeRange[1]
		curr := startTime.Time

		for !curr.After(endTime.Time) {
			timeSlots = append(timeSlots, domain.EventDateTime{Time: curr})
			curr = curr.Add(time.Duration(timeSlotGapMins) * time.Minute)
		}
	}
	return timeSlots, nil
}

func convertTimeSlotsToRanges(timeSlots []domain.EventDateTime) ([][]domain.EventDateTime, error) {
	if len(timeSlots) == 0 {
		return [][]domain.EventDateTime{}, nil
	}

	sort.Slice(timeSlots, func(i, j int) bool {
		return timeSlots[i].Before(timeSlots[j].Time)
	})

	// TODO: consider potential bug when earliest time is 12:45pm and timeGap changes from 15 to 30
	startOfRange := timeSlots[0]
	prev := startOfRange
	timeRanges := [][]domain.EventDateTime{}
	for i := 1; i < len(timeSlots); i++ {
		curr := timeSlots[i]
		if curr == prev || curr.IsZero() {
			continue
		}

		if curr.Sub(prev.Time) != time.Duration(timeSlotGapMins)*time.Minute {
			timeRanges = append(timeRanges, []domain.EventDateTime{startOfRange, prev})
			startOfRange = curr
		}
		prev = curr
	}
	if len(timeRanges) == 0 || !timeRanges[len(timeRanges)-1][1].Equal(prev.Time) {
		timeRanges = append(timeRanges, []domain.EventDateTime{startOfRange, prev})
	}

	return timeRanges, nil
}

func (e *EventService) UpsertUserEventAvailability(ctx context.Context, eventResponse domain.EventResponse) error {
	timeRanges, err := convertTimeSlotsToRanges(eventResponse.Availabiliies)
	if err != nil {
		return err
	}

	sql, args, err := sq.Insert("event_responses").
		Columns("user_id", "event_id", "alias", "availabilities").
		Values(eventResponse.UserID, eventResponse.EventID, eventResponse.Alias, timeRanges).
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT (user_id, event_id, alias) DO UPDATE SET availabilities = ?", timeRanges).
		ToSql()

	if err != nil {
		return err
	}

	_, err = e.db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
