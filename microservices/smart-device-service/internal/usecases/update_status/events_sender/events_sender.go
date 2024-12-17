package events_sender

import (
	"context"
	"errors"
	"fmt"
	"time"

	"smart-device-service/internal/usecases/dto"
	"smart-device-service/internal/usecases/update_status/storage"
	storageDto "smart-device-service/internal/usecases/update_status/storage/dto"
)

type eventsStorage interface {
	GetNewEvent(ctx context.Context) (storageDto.Event, error)
	SetDone(ctx context.Context, eventID int64) error
}

type publisher interface {
	PublishStatusChanged(ctx context.Context, deviceID int, status dto.DeviceStatus) error
}

type Sender struct {
	storage   eventsStorage
	publisher publisher
}

func NewSender(storage eventsStorage, publisher publisher) *Sender {
	return &Sender{storage: storage, publisher: publisher}
}

func (s *Sender) StartProcessEvents(ctx context.Context, timeout time.Duration) {
	ticker := time.NewTicker(time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("stopping event processing")
				return
			case <-ticker.C:
			}

			event, err := s.storage.GetNewEvent(ctx)
			if err != nil {
				if !errors.Is(err, storage.ErrEventNotFound) {
					fmt.Printf("failed to get new event: %v", err)
				}
				continue
			}

			if err := s.sendMessage(ctx, event); err != nil {
				fmt.Printf("send message to queue: %v", err)
				continue
			}

			if err := s.storage.SetDone(ctx, event.ID); err != nil {
				fmt.Printf("failed to set event done: %v", err)
			}
		}
	}()
}

func (s *Sender) sendMessage(ctx context.Context, event storageDto.Event) error {
	deviceStatus := dto.DeviceStatus{
		Status: event.NewStatus,
	}
	if err := s.publisher.PublishStatusChanged(ctx, event.DeviceID, deviceStatus); err != nil {
		return fmt.Errorf("publish status update to queue: %w", err)
	}

	fmt.Printf("device new status: %v for device id: %d published\n", event.NewStatus, event.DeviceID)

	return nil
}
