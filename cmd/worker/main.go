package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/consumer"
	"github.com/lhilove/apigateway-event-ingestion-pipeline/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dbpool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer dbpool.Close()

	repo := repository.NewEventRepository(dbpool)

	c, err := consumer.NewConsumer(os.Getenv("NATS_URL"), repo)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	fmt.Println("ingestion worker started")

	if err := c.Start(ctx); err != nil {
		log.Fatalf("consumer error: %v", err)
	}

	fmt.Println("ingestion worker shutting down")
}
