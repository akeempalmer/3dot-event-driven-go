package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
)

func (h Handler) AppendToTracker(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	slog.Info("Appending ticket to the tracker")

	err := h.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency})
	if err != nil {
		return err
	}

	err = h.StoreTicketContent(ctx, event)

	return err
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

func (h Handler) StoreTicketContent(ctx context.Context, event *entities.TicketBookingConfirmed) error {

	fileID := fmt.Sprintf("%s-ticket.html", event.TicketID)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal the payload")
	}

	fileContent := string(payload)

	resp, err := h.apiClients.Files.PutFilesFileIdContentWithTextBodyWithResponse(ctx, fileID, fileContent)
	if err != nil {
		return fmt.Errorf("failed to put file content")
	}

	if resp.StatusCode() == http.StatusConflict {
		log.FromContext(ctx).With("file", fileID).Info("file already exists")
		return nil
	}

	return nil
}
