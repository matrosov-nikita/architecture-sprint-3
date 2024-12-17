package get_device_telemetry

import (
	"context"
	"encoding/json"
	"fmt"

	"smart-telemetry-service/internal/usecases/dto"
	storageDto "smart-telemetry-service/internal/usecases/get_device_telemetry/storage/dto"
)

type sensorType string

const (
	sensorTemperatureType sensorType = "temperature"
)

type eventsStorage interface {
	GetEventsByDeviceID(ctx context.Context, deviceID int64) ([]storageDto.StorageEvent, error)
}

type GetEventsUsecase struct {
	eventsStorage eventsStorage
}

func NewGetEventsUsecase(eventsStorage eventsStorage) *GetEventsUsecase {
	return &GetEventsUsecase{eventsStorage: eventsStorage}
}

type temperatureData struct {
	Temperature float64 `json:"temperature"`
}

func (s *GetEventsUsecase) GetEventsByDeviceID(ctx context.Context, deviceID int64) ([]dto.SensorTemperatureEvent, error) {
	resultEvents := make([]dto.SensorTemperatureEvent, 0)
	storageEvents, err := s.eventsStorage.GetEventsByDeviceID(ctx, deviceID)

	if err != nil {
		return nil, fmt.Errorf("get events by device id: %w", err)
	}

	for _, event := range storageEvents {
		sensorEvent := dto.SensorTemperatureEvent{
			DeviceId:  event.DeviceID,
			Type:      event.EventType,
			OccuredOn: event.OccuredOn,
		}

		if event.EventType == string(sensorTemperatureType) {
			var t temperatureData
			if err := json.Unmarshal(event.Data, &t); err == nil {
				sensorEvent.Temperature = t.Temperature
			}
		}

		resultEvents = append(resultEvents, sensorEvent)
	}

	return resultEvents, nil
}
