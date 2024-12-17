package storage

import (
	"context"
	"database/sql"
	"fmt"

	storageDto "smart-telemetry-service/internal/usecases/get_device_telemetry/storage/dto"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetEventsByDeviceID(ctx context.Context, deviceID int64) ([]storageDto.StorageEvent, error) {
	query := `
		SELECT id, device_id, event_type, data, occured_on
		FROM devices_events
		WHERE device_id = $1;
	`

	rows, err := s.db.QueryContext(ctx, query, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var events []storageDto.StorageEvent
	for rows.Next() {
		var event storageDto.StorageEvent

		err := rows.Scan(&event.ID, &event.DeviceID, &event.EventType, &event.Data, &event.OccuredOn)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return events, nil
}
