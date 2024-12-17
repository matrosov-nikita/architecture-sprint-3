package get_device

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-device-service/internal/usecases/dto"
	"smart-device-service/internal/usecases/get_device"
)

type deviceGetter interface {
	GetDevice(ctx context.Context, deviceID int) (dto.Device, error)
}

type Handler struct {
	deviceGetter deviceGetter
}

func NewHandler(deviceGetter deviceGetter) *Handler {
	return &Handler{deviceGetter: deviceGetter}
}

func (h *Handler) Handle(c *gin.Context) {
	paramDeviceID := c.Param("deviceId")

	deviceID, err := strconv.Atoi(paramDeviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.deviceGetter.GetDevice(c.Request.Context(), deviceID)
	if err != nil {
		if errors.Is(err, get_device.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
			return
		}
		fmt.Printf("get device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Status update internal error"})
		return
	}

	c.JSON(http.StatusOK, device)
}
