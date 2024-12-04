package dto

type DeviceStatus struct {
	Status string `json:"status"` // "on" or "off"
}

type DeviceCommand struct {
	DeviceID int    `json:"deviceId"`
	Command  string `json:"command"` // "turn_on" or "turn_off"
}
