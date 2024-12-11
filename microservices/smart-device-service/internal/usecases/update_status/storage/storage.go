package storage

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) UpdateDeviceStatus(deviceID int, newStatus string) error {
	query := `UPDATE devices SET status = $1 WHERE id = $2`

	_, err := s.db.Exec(query, newStatus, deviceID)
	if err != nil {
		return fmt.Errorf("failed to update device status in storage: %v", err)
	}

	return nil
}
