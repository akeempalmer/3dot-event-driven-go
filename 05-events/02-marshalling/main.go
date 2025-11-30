package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type PaymentCompleted struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	CompletedAt string `json:"completed_at"`
}

func main() {
	logger := watermill.NewSlogLogger(nil)

	router := message.NewDefaultRouter(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddHandler("payment_completed_handler", "payment-completed", sub, "order-confirmed", pub, func(msg *message.Message) ([]*message.Message, error) {
		var paymentPayload PaymentCompleted

		err := json.Unmarshal(msg.Payload, &paymentPayload)
		if err != nil {
			return nil, err
		}

		orderEvent := struct {
			OrderID     string `json:"order_id"`
			ConfirmedAt string `json:"confirmed_at"`
		}{
			OrderID:     paymentPayload.OrderID,
			ConfirmedAt: paymentPayload.CompletedAt,
		}

		orderMsgPayload, err := json.Marshal(orderEvent)
		if err != nil {
			return nil, err
		}

		return []*message.Message{message.NewMessage(watermill.NewUUID(), []byte(orderMsgPayload))}, nil
	})

	err = router.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
