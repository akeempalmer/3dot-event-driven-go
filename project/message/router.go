package message

import (
	"context"
	"encoding/json"
	"tickets/entities"
	"tickets/message/event"
	"tickets/middlewares"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/redis/go-redis/v9"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptPayload) error
}

func NewWatermillRouter(
	receiptsService ReceiptsService,
	spreadsheetsAPI SpreadsheetsAPI,
	rdb *redis.Client,
	watermilLogger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(watermilLogger)

	handler := event.NewHandler(spreadsheetsAPI, receiptsService)

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, watermilLogger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, watermilLogger)
	if err != nil {
		panic(err)
	}

	bookingFailedSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "failed-booking-subscriber",
	}, watermilLogger)
	if err != nil {
		panic(err)
	}
	router.AddMiddleware(middlewares.LogHeader{}.CorrelationMiddleware)

	router.AddMiddleware(middlewares.LogHeader{
		MessageID: watermill.NewUUID(),
		Payload:   "Handling a message",
	}.LoggerMiddleware)

	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermilLogger,
	}.Middleware)

	// router.AddConsumerHandler(
	// 	"issue_receipt",
	// 	"issue-receipt",
	// 	issueReceiptSub,
	// 	func(msg *message.Message) error {
	// 		var ticketReceiptEvent entities.IssueReceiptPayload

	// 		err := json.Unmarshal(msg.Payload, &ticketReceiptEvent)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to unmarshal ticket event")
	// 		}

	// 		err = receiptsService.IssueReceipt(msg.Context(), ticketReceiptEvent)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to issue receipt: %w", err)
	// 		}
	// 		return nil
	// 	},
	// )

	// router.AddConsumerHandler(
	// 	"append_to_tracker",
	// 	"append-to-tracker",
	// 	appendToTrackerSub,
	// 	func(msg *message.Message) error {
	// 		var ticketEvent entities.AppendToTrackerPayload

	// 		err := json.Unmarshal(msg.Payload, &ticketEvent)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to unmarashal ticket event")
	// 		}

	// 		err = spreadsheetsAPI.AppendRow(msg.Context(), "tickets-to-print", []string{ticketEvent.TicketID, ticketEvent.CustomerEmail, ticketEvent.Price.Amount, ticketEvent.Price.Currency})
	// 		if err != nil {
	// 			return fmt.Errorf("failed to append to tracker: %w", err)
	// 		}
	// 		return nil
	// 	},
	// )

	router.AddConsumerHandler("issue_receipt", "TicketBookingConfirmed", issueReceiptSub, func(msg *message.Message) error {

		var payload entities.TicketBookingConfirmed

		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		return handler.IssueReceipt(msg.Context(), payload)
	})

	router.AddConsumerHandler("append_to_tracker", "TicketBookingConfirmed", appendToTrackerSub, func(msg *message.Message) error {

		var payload entities.TicketBookingConfirmed

		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		return handler.AppendToTracker(msg.Context(), payload)
	})

	router.AddConsumerHandler("publish_failed_booking", "TicketBookingCanceled", bookingFailedSub, func(msg *message.Message) error {
		var payload entities.TicketBookingCanceled

		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		return handler.AppendCancelToTracker(msg.Context(), payload)
	})

	return router
}
