package subscribers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	usecaseDto "smart-device-service/internal/usecases/dto"
)

const (
	deviceCommandsTopic = "device_commands"
)

type commandHandler interface {
	SendCommand(deviceID int, command usecaseDto.DeviceCommand) error
}

type CommandSubscriber struct {
	kafkaReader    *kafka.Reader
	commandHandler commandHandler
}

func NewCommandSubscriber(kafkaBrokerAddress string, handler commandHandler) *CommandSubscriber {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaBrokerAddress},
		Topic:       deviceCommandsTopic,
		Partition:   0,
		MaxBytes:    10e6, // 10MB
		GroupID:     "my-group",
		StartOffset: kafka.LastOffset,
	})

	return &CommandSubscriber{kafkaReader: r, commandHandler: handler}
}

func (s *CommandSubscriber) Stop() error {
	return s.kafkaReader.Close()
}

func (s *CommandSubscriber) Run(ctx context.Context) error {
	for {
		m, err := s.kafkaReader.ReadMessage(ctx)
		if err != nil {
			break
		}

		var cmd usecaseDto.DeviceCommand
		if err := json.Unmarshal(m.Value, &cmd); err != nil {
			fmt.Printf("command subscriber decode msg value: %v\n", err)
			continue
		}

		if err := s.commandHandler.SendCommand(cmd.DeviceID, usecaseDto.DeviceCommand{
			Command: cmd.Command,
			UserID:  cmd.UserID,
		}); err != nil {
			return fmt.Errorf("send command: %v\n", err)
		}
	}

	return nil
}
