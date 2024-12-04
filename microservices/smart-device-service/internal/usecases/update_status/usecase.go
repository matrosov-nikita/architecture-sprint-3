package update_status

import (
	"context"
	"fmt"

	"smart-device-service/internal/usecases/dto"
)

type publisher interface {
	PublishStatusChanged(ctx context.Context, deviceID string, status dto.DeviceStatus) error
}
type UpdateStatusUsecase struct {
	publisher publisher
}

func NewUpdateStatusUsecase(publisher publisher) *UpdateStatusUsecase {
	return &UpdateStatusUsecase{
		publisher: publisher,
	}
}

func (u *UpdateStatusUsecase) UpdateDeviceStatus(ctx context.Context, deviceID string, update dto.DeviceStatus) error {
	// 1. Обновляем девайс в БД

	// 2. Публикуем сообщение в очередь
	if err := u.publisher.PublishStatusChanged(ctx, deviceID, update); err != nil {
		return fmt.Errorf("publish status update to queue: %v", err)
	}

	return nil
}
