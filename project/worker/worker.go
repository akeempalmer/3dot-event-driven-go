package worker

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
)

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

type PubMessage struct {
	ID      string
	Payload []byte
}

type Message struct {
	Task     Task
	TicketID string
}

type Worker struct {
	queue chan Message

	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

func NewWorker(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) *Worker {
	return &Worker{
		queue: make(chan Message, 100),

		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}
}

func (w *Worker) Send(msgs ...Message) {
	for _, msg := range msgs {
		w.queue <- msg
	}
}

func (w *Worker) Run(ctx context.Context, receiptSub, sheetSub *redisstream.Subscriber) {

	go func() {
		messages, err := receiptSub.Subscribe(context.Background(), "issue-receipt")
		if err != nil {
			panic(err)
		}

		for msg := range messages {
			err := w.receiptsService.IssueReceipt(msg.Context(), string(msg.Payload))
			if err != nil {
				slog.With("error", err).Error("failed to issue the receipt")
				msg.Nack()
				continue
			}

			msg.Ack()
		}
	}()

	go func() {
		messages, err := sheetSub.Subscribe(context.Background(), "append-to-tracker")
		if err != nil {
			panic(err)
		}

		for msg := range messages {
			err := w.spreadsheetsAPI.AppendRow(msg.Context(), "tickets-to-print", []string{string(msg.Payload)})
			if err != nil {
				slog.With("error", err).Error("failed to append to tracker")
				msg.Nack()
				continue
			}

			msg.Ack()
		}
	}()

	// for msg := range w.queue {
	// 	switch msg.Task {
	// 	case TaskIssueReceipt:
	// 		err := w.receiptsService.IssueReceipt(ctx, msg.TicketID)
	// 		if err != nil {
	// 			slog.With("error", err).Error("failed to issue the receipt")
	// 			w.Send(msg)
	// 		}
	// 	case TaskAppendToTracker:
	// 		err := w.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{msg.TicketID})
	// 		if err != nil {
	// 			slog.With("error", err).Error("failed to append to tracker")
	// 			w.Send(msg)
	// 		}
	// 	}
	// }
}
