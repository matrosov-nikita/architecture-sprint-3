package update_status

import (
	"context"
	"fmt"

	"smart-device-service/internal/usecases/dto"
)

type publisher interface {
	PublishStatusChanged(ctx context.Context, deviceID int, status dto.DeviceStatus) error
}

type storage interface {
	UpdateDeviceStatus(deviceID int, newStatus string) error
}

type UpdateStatusUsecase struct {
	publisher publisher
	storage   storage
}

func NewUpdateStatusUsecase(publisher publisher, storage storage) *UpdateStatusUsecase {
	return &UpdateStatusUsecase{publisher: publisher, storage: storage}
}

func (u *UpdateStatusUsecase) UpdateDeviceStatus(ctx context.Context, deviceID int, update dto.DeviceStatus) error {
	// 1. Обновляем девайс в БД
	if err := u.storage.UpdateDeviceStatus(deviceID, update.Status); err != nil {
		return err
	}

	// 2. Публикуем сообщение в очередь об изменении статуса
	if err := u.publisher.PublishStatusChanged(ctx, deviceID, update); err != nil {
		return fmt.Errorf("publish status update to queue: %v", err)
	}

	return nil
}
