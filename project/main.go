package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/redis/go-redis/v9"

	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	// Create a new context, and pass it to signal.NotifyContext. The incoming interrupt signal will cancel the context.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apiClients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	err = redisClient.Ping(context.Background()).Err()

	if err != nil {
		panic(err)
	}

	defer redisClient.Close()

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	err = service.New(
		spreadsheetsAPI,
		receiptsService,
		redisClient,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
