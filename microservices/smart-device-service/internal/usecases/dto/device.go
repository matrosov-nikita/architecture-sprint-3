package dto

type DeviceStatus struct {
	Status string `json:"status"`
}

type DeviceCommand struct {
	Command string `json:"command"`
	UserID  int    `json:"user_id"`
}
