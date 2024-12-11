package send_command

import (
	"fmt"

	"smart-device-service/internal/usecases/dto"
)

type SendCommandUsecase struct{}

func NewSendCommandUsecase() *SendCommandUsecase {
	return &SendCommandUsecase{}
}

// SendCommand Заглушка отправления команды на устройство.
func (s *SendCommandUsecase) SendCommand(deviceID int, command dto.DeviceCommand) error {
	fmt.Printf("sending command: %v by user %d to device id: %d\n", command.Command, command.UserID, deviceID)
	return nil
}
