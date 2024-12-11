package get_device

import (
	"context"
	"errors"

	"smart-device-service/internal/usecases/dto"
	storageErrors "smart-device-service/internal/usecases/get_device/storage"
	storageDto "smart-device-service/internal/usecases/get_device/storage/dto"
)

var ErrDeviceNotFound = errors.New("device not found")

type storage interface {
	GetDevice(ctx context.Context, deviceID int) (storageDto.Device, error)
}

type GetDeviceUsecase struct {
	storage storage
}

func NewGetDeviceUsecase(storage storage) *GetDeviceUsecase {
	return &GetDeviceUsecase{storage: storage}
}

func (u *GetDeviceUsecase) GetDevice(ctx context.Context, deviceID int) (dto.Device, error) {
	storageDevice, err := u.storage.GetDevice(ctx, deviceID)
	if err != nil {
		if errors.Is(err, storageErrors.ErrDeviceNotFound) {
			return dto.Device{}, ErrDeviceNotFound
		}
		return dto.Device{}, err
	}

	d := dto.Device{
		ID:           storageDevice.ID,
		SerialNumber: storageDevice.SerialNumber,
		UserID:       storageDevice.UserID,
		Name:         storageDevice.Name,
		CreatedAt:    storageDevice.CreatedAt,
		Status:       storageDevice.Status,
	}

	return d, nil
}
