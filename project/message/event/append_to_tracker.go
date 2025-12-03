package event

import (
	"context"
	"log/slog"
	"tickets/entities"
)

func (h Handler) AppendToTracker(ctx context.Context, event entities.TicketBookingConfirmed) error {
	slog.Info("Appending ticket to the tracker")

	payload := entities.AppendToTrackerPayload{
		TicketID:      event.TicketID,
		CustomerEmail: event.CustomerEmail,
		Price:         event.Price,
	}

	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency})
}
