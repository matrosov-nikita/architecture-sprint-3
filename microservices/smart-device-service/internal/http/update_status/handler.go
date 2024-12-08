package update_status

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-device-service/internal/usecases/dto"
)

type statusUpdater interface {
	UpdateDeviceStatus(ctx context.Context, deviceID int, update dto.DeviceStatus) error
}

type Handler struct {
	statusUpdater statusUpdater
}

func NewHandler(statusUpdater statusUpdater) *Handler {
	return &Handler{statusUpdater: statusUpdater}
}

func (h *Handler) Handle(c *gin.Context) {
	paramDeviceID := c.Param("deviceId")

	deviceID, err := strconv.Atoi(paramDeviceID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid device ID"})
		return
	}

	var statusUpdate dto.DeviceStatus
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Обновляем статус устройства
	if err := h.statusUpdater.UpdateDeviceStatus(c.Request.Context(), deviceID, statusUpdate); err != nil {
		// Для упрощения логируем прямо в консоль.
		fmt.Printf("update device status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Status update internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}
