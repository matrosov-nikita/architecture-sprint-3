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
		MaxBytes:    10e6, // 10MB
		GroupID:     "my-group-2",
		StartOffset: kafka.LastOffset,
	})

	return &SensorDataSubscriber{kafkaReader: r, dataHandler: dataHandler}
}

func (s *SensorDataSubscriber) Stop() error {
	return s.kafkaReader.Close()
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

		if err := s.dataHandler.HandleSensorEvents(context.Background(), []dto.SensorTemperatureEvent{e}); err != nil {
			return fmt.Errorf("handle sensor event: %v\n", err)
		}
	}

	return nil
}
