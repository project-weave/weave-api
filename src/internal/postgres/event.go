package postgres

import (
	"context"
	"log"
	"sort"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/project-weave/weave-api/src/internal/types/event"
)

type EventService struct {
	db *DB
}

func NewEventService(db *DB) *EventService {
	return &EventService{
		db: db,
	}
}

func (es *EventService) AddEvent(ctx context.Context, e *event.Event) (event.EventUUID, error) {
	sql, args, err := sq.Insert("events").
		Columns("name", "is_specific_dates", "start_time", "end_time", "dates").
		Values(e.Name, e.IsSpecificDates, e.StartTime, e.EndTime.Time, e.Dates).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return event.EventUUID{}, err
	}

	row := es.db.pool.QueryRow(ctx, sql, args...)
	err = row.Scan(&e.ID)
	if err != nil {
		return event.EventUUID{}, err
	}

	return e.ID, err
}

func (es *EventService) GetEvent(ctx context.Context, eID event.EventUUID) (*event.Event, []event.EventResponse, error) {
	tx, err := es.db.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Fatalf("tx.Rollback failed: %v", err)
		}
	}()

	sql, args, err := sq.Select("id", "name", "is_specific_dates", "start_time", "end_time", "dates").
		From("events").
		Where(sq.Eq{"id": eID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, err
	}

	row := tx.QueryRow(ctx, sql, args...)
	e := event.Event{}
	err = row.Scan(&e.ID, &e.Name, &e.IsSpecificDates, &e.StartTime, &e.EndTime, &e.Dates)
	if err != nil {
		return nil, nil, err
	}

	sql, args, err = sq.Select("user_id", "event_id", "alias", "availabilities").
		From("event_responses").
		Where(sq.Eq{"event_id": eID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, nil, err
	}
	responses := []event.EventResponse{}
	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, nil, err
	}
	for rows.Next() {
		response := event.EventResponse{}
		timeRanges := [][]event.EventDateTime{}

		err = rows.Scan(&response.UserID, &response.EventID, &response.Alias, &timeRanges)
		if err != nil {
			return nil, nil, err
		}

		timeSlots, err := convertTimeRangesToTimeSlots(timeRanges)
		if err != nil {
			return nil, nil, err
		}
		response.Availabiliies = timeSlots
		responses = append(responses, response)
	}

	return &e, responses, nil
}

var timeSlotGapMins = 30

func convertTimeRangesToTimeSlots(ranges [][]event.EventDateTime) ([]event.EventDateTime, error) {
	var timeSlots []event.EventDateTime
	for _, timeRange := range ranges {
		startTime := timeRange[0]
		endTime := timeRange[1]
		curr := startTime.Time

		for !curr.After(endTime.Time) {
			timeSlots = append(timeSlots, event.EventDateTime{Time: curr})
			curr = curr.Add(time.Duration(timeSlotGapMins) * time.Minute)
		}
	}
	return timeSlots, nil
}

func convertTimeSlotsToRanges(timeSlots []event.EventDateTime) ([][]event.EventDateTime, error) {
	if len(timeSlots) == 0 {
		return [][]event.EventDateTime{}, nil
	}

	sort.Slice(timeSlots, func(i, j int) bool {
		return timeSlots[i].Before(timeSlots[j].Time)
	})

	// TODO: consider potential bug when earliest time is 12:45pm and timeGap changes from 15 to 30
	startOfRange := timeSlots[0]
	prev := startOfRange
	timeRanges := [][]event.EventDateTime{}
	for i := 1; i < len(timeSlots); i++ {
		curr := timeSlots[i]
		if curr == prev || curr.IsZero() {
			continue
		}

		if curr.Sub(prev.Time) != time.Duration(timeSlotGapMins)*time.Minute {
			timeRanges = append(timeRanges, []event.EventDateTime{startOfRange, prev})
			startOfRange = curr
		}
		prev = curr
	}
	if len(timeRanges) == 0 || !timeRanges[len(timeRanges)-1][1].Equal(prev.Time) {
		timeRanges = append(timeRanges, []event.EventDateTime{startOfRange, prev})
	}

	return timeRanges, nil
}

func (e *EventService) UpsertUserEventAvailability(ctx context.Context, response event.EventResponse) error {
	timeRanges, err := convertTimeSlotsToRanges(response.Availabiliies)
	if err != nil {
		return err
	}

	sql, args, err := sq.Insert("event_responses").
		Columns("user_id", "event_id", "alias", "availabilities").
		Values(response.UserID, response.EventID, response.Alias, timeRanges).
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
