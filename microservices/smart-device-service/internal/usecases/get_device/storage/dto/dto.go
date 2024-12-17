package dto

import "time"

type Device struct {
	ID           int
	Name         string
	SerialNumber string
	UserID       string
	Status       string
	CreatedAt    time.Time
}
