package event

import (
	"context"
	"log/slog"
	"tickets/entities"
)

func (h Handler) AppendToTracker(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	slog.Info("Appending ticket to the tracker")

	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})
}

func (h Handler) AppendCancelToTracker(ctx context.Context, event *entities.TicketBookingCanceled) error {
	slog.Info("Appending ticket to the canceled tracker")

	return h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-refund", []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})
}

func (h Handler) SaveTicketToDatabase(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	slog.Info("Saving ticket to the database")

	err := h.ticketRepo.Save(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (h Handler) DeleteTicketFromDatabase(ctx context.Context, event *entities.TicketBookingCanceled) error {
	slog.Info("Deleting ticket from the database")

	err := h.ticketRepo.Delete(ctx, event.TicketID)
	if err != nil {
		return err
	}

	return nil
}
