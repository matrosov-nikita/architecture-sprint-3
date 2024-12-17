package update_status

import (
	"context"
	"fmt"

	"smart-device-service/internal/usecases/dto"
)

type storage interface {
	UpdateDeviceStatus(ctx context.Context, deviceID int, newStatus string) error
}

type UpdateStatusUsecase struct {
	storage storage
}

func NewUpdateStatusUsecase(storage storage) *UpdateStatusUsecase {
	return &UpdateStatusUsecase{storage: storage}
}

func (u *UpdateStatusUsecase) UpdateDeviceStatus(ctx context.Context, deviceID int, update dto.DeviceStatus) error {
	if err := u.storage.UpdateDeviceStatus(ctx, deviceID, update.Status); err != nil {
		return fmt.Errorf("update device status: %w", err)
	}

	return nil
}
