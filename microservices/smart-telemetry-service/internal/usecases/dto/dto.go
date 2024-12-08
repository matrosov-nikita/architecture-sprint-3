package dto

import "time"

type SensorTemperatureEvent struct {
	DeviceId    int       `json:"device_id"`
	Temperature float64   `json:"temperature"`
	Type        string    `json:"type"`
	OccuredOn   time.Time `json:"occured_on"`
}