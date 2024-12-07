package handle_sensor_event

import (
	"context"
	"encoding/json"
	"fmt"

	"smart-telemetry-service/internal/usecases/dto"
	storageDto "smart-telemetry-service/internal/usecases/handle_sensor_event/storage/dto"
)

type sensorType string

const (
	sensorTemperatureType sensorType = "temperature"
)

type eventsStorage interface {
	SaveEvents(ctx context.Context, events []storageDto.StorageEvent) error
}

type HandleSensorEventUsecase struct {
	eventsStorage eventsStorage
}

func NewHandleSensorEventUsecase(eventsStorage eventsStorage) *HandleSensorEventUsecase {
	return &HandleSensorEventUsecase{eventsStorage: eventsStorage}
}

type temperatureData struct {
	Temperature float64 `json:"temperature"`
}

func (s *HandleSensorEventUsecase) HandleSensorEvents(ctx context.Context, events []dto.SensorTemperatureEvent) error {
	storageEvents := make([]storageDto.StorageEvent, 0, len(events))
	for _, event := range events {
		stEvent := storageDto.StorageEvent{
			DeviceID:  event.DeviceId,
			EventType: event.Type,
			Data:      nil,
			OccuredOn: event.OccuredOn,
		}
		if event.Type == string(sensorTemperatureType) {
			data := temperatureData{
				Temperature: event.Temperature,
			}
			storageSensorData, err := json.Marshal(data)
			if err != nil {
				return fmt.Errorf("encode sensor temperature data: %v", err)
			}
			stEvent.Data = storageSensorData
		}

		storageEvents = append(storageEvents, stEvent)
	}
	if err := s.eventsStorage.SaveEvents(ctx, storageEvents); err != nil {
		return fmt.Errorf("save events: %v", err)
	}

	return nil
}
