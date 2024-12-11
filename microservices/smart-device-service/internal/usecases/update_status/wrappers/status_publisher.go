package wrappers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"smart-device-service/internal/usecases/dto"
)

const (
	statusTopic = "device_statuses"
)

type StatusChangedPublisher struct {
	kafkaWriter *kafka.Writer
}

func NewStatusChangedPublisher(brokerAddress string) *StatusChangedPublisher {
	return &StatusChangedPublisher{
		kafkaWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerAddress},
			Topic:   statusTopic,
		}),
	}
}

type PublishStatus struct {
	DeviceID int    `json:"device_id"`
	Name     string `json:"name"`
}

func (p *StatusChangedPublisher) PublishStatusChanged(ctx context.Context, deviceID int, status dto.DeviceStatus) error {
	statusJSON, _ := json.Marshal(PublishStatus{DeviceID: deviceID, Name: status.Status})
	if err := p.kafkaWriter.WriteMessages(ctx, []kafka.Message{
		{Value: statusJSON},
	}...); err != nil {
		return fmt.Errorf("write device status changed: %v", err)
	}

	return nil
}
