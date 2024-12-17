package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"smart-device-service/internal/usecases/update_status/storage/dto"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) UpdateDeviceStatus(ctx context.Context, deviceID int, newStatus string) error {
	updateDeviceStatusQuery := `UPDATE devices SET status = $1 WHERE id = $2`

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	defer func() { _ = tx.Rollback() }()

	if _, err = tx.ExecContext(ctx, updateDeviceStatusQuery, newStatus, deviceID); err != nil {
		return fmt.Errorf("update device status: exec context: %w", err)
	}
	if err := s.saveEvent(ctx, tx, deviceID, newStatus); err != nil {
		return fmt.Errorf("save outbox event: exec context: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to update device status in storage: tx commit: %w", err)
	}

	return nil
}

type eventPayload struct {
	DeviceID  int    `json:"device_id"`
	NewStatus string `json:"new_status"`
}

func (s *Storage) GetNewEvent(ctx context.Context) (dto.Event, error) {
	selectNewEventQuery := `
		SELECT id, payload
		FROM outbox_events
		WHERE status = 'new'
		ORDER BY created_at ASC
		LIMIT 1;
	`

	var (
		id      int64
		payload string
	)

	err := s.db.QueryRowContext(ctx, selectNewEventQuery).Scan(&id, &payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Event{}, ErrEventNotFound
		}
		return dto.Event{}, fmt.Errorf("get new event: query row context: %w", err)
	}

	var p eventPayload
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		return dto.Event{}, fmt.Errorf("get new event: unmarshal payload: %w", err)
	}

	return dto.Event{
		ID:        id,
		DeviceID:  p.DeviceID,
		NewStatus: p.NewStatus,
	}, nil
}

func (s *Storage) SetDone(ctx context.Context, eventID int64) error {
	updateStatusQuery := `
		UPDATE outbox_events
		SET status = 'done'
		WHERE id = $1;
	`

	_, err := s.db.ExecContext(ctx, updateStatusQuery, eventID)
	if err != nil {
		return fmt.Errorf("set event as done: exec context: %w", err)
	}

	return nil
}

func (s *Storage) saveEvent(ctx context.Context, tx *sql.Tx, deviceID int, newStatus string) error {
	insertOutboxEventQuery := `
		INSERT INTO outbox_events (payload) VALUES ($1);
	`

	payload, err := json.Marshal(eventPayload{DeviceID: deviceID, NewStatus: newStatus})
	if err != nil {
		return fmt.Errorf("marshal event payload: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertOutboxEventQuery, payload); err != nil {
		return fmt.Errorf("save event: exec context: %w", err)
	}

	return nil
}
