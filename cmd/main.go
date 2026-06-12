package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/handler"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/middleware"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/queue"
)

func main() {
	publisher, err := queue.NewPublisher("nats://localhost:4222")
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}

	router := gin.Default()
	eventHandler := handler.NewEventHandler(publisher)

	// Group routes that require authentication
	api := router.Group("/api/v1")
	api.Use(middleware.RateLimit(5, 10)) // 5 req/sec, burst of 10
	api.Use(middleware.APIKeyAuth())
	{
		api.POST("/events", eventHandler.IngestEvent)
	}

	router.Run(":8080")
}
