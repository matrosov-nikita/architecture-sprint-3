package dto

import "time"

type DeviceStatus struct {
	Status string `json:"status"`
}

type DeviceCommand struct {
	DeviceID int    `json:"device_id"`
	Command  string `json:"command"`
	UserID   int    `json:"user_id"`
}

type Device struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	SerialNumber string    `json:"serial_number"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
}
