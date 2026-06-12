package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/domain"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/queue"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/validator"
)

type EventHandler struct {
	publisher *queue.Publisher
}

func NewEventHandler(p *queue.Publisher) *EventHandler {
	return &EventHandler{publisher: p}
}

func (h *EventHandler) IngestEvent(c *gin.Context) {
	var event domain.Event
	// structural validation (required fields, types)
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid event payload",
			"details": err.Error(),
		})
		return
	}

	// semantic validation (severity values, etc.)
	if err := validator.ValidateEvent(event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// forward to queue.
	if err := h.publisher.Publish(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to queue event",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}
