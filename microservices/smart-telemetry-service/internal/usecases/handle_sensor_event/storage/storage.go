package storage

import (
	"context"
	"database/sql"
	"fmt"

	"smart-telemetry-service/internal/usecases/handle_sensor_event/storage/dto"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveEvents(ctx context.Context, events []dto.StorageEvent) error {
	query := `
		INSERT INTO devices_events (device_id, event_type, data, occured_on)
		VALUES ($1, $2, $3, $4);
	`

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("tx begin: %v", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, e := range events {
		_, err := tx.ExecContext(ctx, query, e.DeviceID, e.EventType, e.Data, e.OccuredOn)
		if err != nil {
			return fmt.Errorf("insert event to table: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("tx commit: %v", err)
	}

	return nil
}
