package get_device_telemetry

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"smart-telemetry-service/internal/usecases/dto"
)

type eventsGetter interface {
	GetEventsByDeviceID(ctx context.Context, deviceID int64) ([]dto.SensorTemperatureEvent, error)
}

type Handler struct {
	eventsGetter eventsGetter
}

func NewHandler(eventsGetter eventsGetter) *Handler {
	return &Handler{eventsGetter: eventsGetter}
}

func (h *Handler) Handle(c *gin.Context) {
	paramDeviceID := c.Param("deviceId")

	deviceID, err := strconv.Atoi(paramDeviceID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid device ID"})
		return
	}

	events, err := h.eventsGetter.GetEventsByDeviceID(c.Request.Context(), int64(deviceID))
	if err != nil {
		fmt.Printf("get telemetry events by device id: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Get events internal error"})
		return
	}

	c.JSON(http.StatusOK, events)
}
