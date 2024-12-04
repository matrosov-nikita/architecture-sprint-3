package update_status

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"smart-device-service/internal/usecases/dto"
)

type statusUpdater interface {
	UpdateDeviceStatus(ctx context.Context, deviceID string, update dto.DeviceStatus) error
}

type Handler struct {
	statusUpdater statusUpdater
}

func NewHandler(statusUpdater statusUpdater) *Handler {
	return &Handler{statusUpdater: statusUpdater}
}

func (h *Handler) Handle(c *gin.Context) {
	deviceID := c.Param("deviceId")

	var statusUpdate dto.DeviceStatus
	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	if statusUpdate.Status != "on" && statusUpdate.Status != "off" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Обновляем статус устройства
	if err := h.statusUpdater.UpdateDeviceStatus(c.Request.Context(), deviceID, statusUpdate); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "status update internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}
