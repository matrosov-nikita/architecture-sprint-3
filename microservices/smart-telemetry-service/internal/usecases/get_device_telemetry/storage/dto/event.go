package dto

import "time"

type StorageEvent struct {
	ID        int
	DeviceID  int
	EventType string
	Data      []byte
	OccuredOn time.Time
}
