package dto

import "time"

type StorageEvent struct {
	DeviceID  int
	EventType string
	Data      []byte
	OccuredOn time.Time
}
