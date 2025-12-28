package message

import (
	"tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

// type SpreadsheetsAPI interface {
// 	AppendRow(ctx context.Context, sheetName string, row []string) error
// }

// type ReceiptsService interface {
// 	IssueReceipt(ctx context.Context, request entities.IssueReceiptPayload) error
// }

func NewWatermillRouter(
	// receiptsService ReceiptsService,
	// spreadsheetsAPI SpreadsheetsAPI,
	// rdb *redis.Client,
	eventProcessorConfig cqrs.EventProcessorConfig,
	eventHandler event.Handler,
	watermilLogger watermill.LoggerAdapter,
) *message.Router {
	router := message.NewDefaultRouter(watermilLogger)

	useMiddlewares(router, watermilLogger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"AppendToTracker",
			eventHandler.AppendToTracker,
		),
		cqrs.NewEventHandler(
			"TicketRefundToSheet",
			eventHandler.AppendCancelToTracker,
		),
		cqrs.NewEventHandler(
			"IssueReceipt",
			eventHandler.IssueReceipt,
		),
		cqrs.NewEventHandler(
			"StoreTicketToDatabase",
			eventHandler.SaveTicketToDatabase,
		),
		cqrs.NewEventHandler(
			"DeleteTicketFromDatabase",
			eventHandler.DeleteTicketFromDatabase,
		),
	)

	// handler := event.NewHandler(spreadsheetsAPI, receiptsService)

	// issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
	// 	Client:        rdb,
	// 	ConsumerGroup: "issue-receipt",
	// }, watermilLogger)
	// if err != nil {
	// 	panic(err)
	// }

	// appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
	// 	Client:        rdb,
	// 	ConsumerGroup: "append-to-tracker",
	// }, watermilLogger)
	// if err != nil {
	// 	panic(err)
	// }

	// bookingFailedSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
	// 	Client:        rdb,
	// 	ConsumerGroup: "failed-booking-subscriber",
	// }, watermilLogger)
	// if err != nil {
	// 	panic(err)
	// }

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

	// router.AddConsumerHandler("issue_receipt", "TicketBookingConfirmed", issueReceiptSub, func(msg *message.Message) error {

	// 	var payload entities.TicketBookingConfirmed

	// 	err := json.Unmarshal(msg.Payload, &payload)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return handler.IssueReceipt(msg.Context(), payload)
	// })

	// router.AddConsumerHandler("append_to_tracker", "TicketBookingConfirmed", appendToTrackerSub, func(msg *message.Message) error {

	// 	var payload entities.TicketBookingConfirmed

	// 	err := json.Unmarshal(msg.Payload, &payload)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return handler.AppendToTracker(msg.Context(), payload)
	// })

	// router.AddConsumerHandler("publish_failed_booking", "TicketBookingCanceled", bookingFailedSub, func(msg *message.Message) error {
	// 	var payload entities.TicketBookingCanceled

	// 	err := json.Unmarshal(msg.Payload, &payload)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	return handler.AppendCancelToTracker(msg.Context(), payload)
	// })

	return router
}
