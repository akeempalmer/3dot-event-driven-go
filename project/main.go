package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	rdsClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	err = rdsClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	ctx := context.Background()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := router.Run(ctx)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	logger := watermill.NewSlogLogger(nil)
	router := message.NewDefaultRouter(logger)

	receiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "receipt_service",
	}, logger)
	if err != nil {
		panic(err)
	}

	sheetSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "sheet_service",
	}, logger)
	if err != nil {
		panic(err)
	}

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddConsumerHandler("receipt_handler", "issue-receipt", receiptSub, func(msg *message.Message) error {
		return receiptsService.IssueReceipt(msg.Context(), string(msg.Payload))
	})

	router.AddConsumerHandler("sheet_handler", "append-to-tracker", sheetSub, func(msg *message.Message) error {
		err := spreadsheetsAPI.AppendRow(msg.Context(), "tickets-to-print", []string{string(msg.Payload)})
		if err != nil {
			return fmt.Errorf("failed to append to tracker: %w", err)
		}
		return err
	})

	err = service.New(
		spreadsheetsAPI,
		receiptsService,
		router,
		publisher,
	).Run(context.Background(), router)
	if err != nil {
		panic(err)
	}
}
