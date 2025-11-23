package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {

	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		logger.Error("failed to create subscriber", err, nil)
	}

	messages, err := subscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		logger.Error("failed to subscribe to topic", err, nil)
	}

	for msg := range messages {
		fmt.Printf("Message ID: %s - %s", string(msg.UUID), string(msg.Payload))
		msg.Ack()
	}

}
