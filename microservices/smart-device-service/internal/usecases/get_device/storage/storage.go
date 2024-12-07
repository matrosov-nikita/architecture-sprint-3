package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"smart-device-service/internal/usecases/get_device/storage/dto"
)

var ErrDeviceNotFound = errors.New("device not found")

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetDevice(ctx context.Context, deviceID int) (dto.Device, error) {
	query := `
			SELECT id, user_id, name, serial_number, status, created_at
			FROM devices
			WHERE id = $1;`

	var device dto.Device
	err := s.db.QueryRowContext(ctx, query, deviceID).Scan(
		&device.ID,
		&device.UserID,
		&device.Name,
		&device.SerialNumber,
		&device.Status,
		&device.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.Device{}, ErrDeviceNotFound
		}
		return dto.Device{}, fmt.Errorf("query row failed: %w", err)
	}

	return device, nil
}
