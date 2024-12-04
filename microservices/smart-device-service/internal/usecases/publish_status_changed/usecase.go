package publish_status_changed

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

func (p *StatusChangedPublisher) PublishStatusChanged(ctx context.Context, deviceID string, status dto.DeviceStatus) error {
	fmt.Println("HERE")
	statusJSON, _ := json.Marshal(status)
	if err := p.kafkaWriter.WriteMessages(ctx, []kafka.Message{
		{
			Key:   []byte(deviceID),
			Value: statusJSON,
		},
	}...); err != nil {
		return fmt.Errorf("write device status changed: %v", err)
	}

	return nil
}
