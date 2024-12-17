package subscribers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"smart-telemetry-service/internal/usecases/dto"
)

const (
	sensorDataTopic = "sensor_data"
	groupID         = "my-group-2"
	readerMaxBytes  = 10 * 1024 * 1024
)

type dataHandler interface {
	HandleSensorEvents(ctx context.Context, events []dto.SensorTemperatureEvent) error
}

type SensorDataSubscriber struct {
	kafkaReader *kafka.Reader
	dataHandler dataHandler
}

func NewSensorDataSubscriber(kafkaBrokerAddress string, dataHandler dataHandler) *SensorDataSubscriber {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaBrokerAddress},
		Topic:       sensorDataTopic,
		Partition:   0,
		MaxBytes:    readerMaxBytes, // 10MB
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})

	return &SensorDataSubscriber{kafkaReader: r, dataHandler: dataHandler}
}

func (s *SensorDataSubscriber) Stop() error {
	return fmt.Errorf("stop sensor data subscriber: %w", s.kafkaReader.Close())
}

func (s *SensorDataSubscriber) Run(ctx context.Context) error {
	for {
		m, err := s.kafkaReader.ReadMessage(ctx)
		if err != nil {
			break
		}

		var e dto.SensorTemperatureEvent
		if err := json.Unmarshal(m.Value, &e); err != nil {
			fmt.Printf("sensor data subscriber decode msg value: %v\n", err)
			continue
		}

		e.OccuredOn = m.Time

		fmt.Printf("got telemetry event: %v\n", e)

		if err := s.dataHandler.HandleSensorEvents(context.Background(), []dto.SensorTemperatureEvent{e}); err != nil {
			return fmt.Errorf("handle sensor events: %w", err)
		}
	}

	return nil
}
