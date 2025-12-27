package http

import (
	"encoding/json"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type TicketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
}

type ticketsConfirmationRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	// for _, ticket := range request.Tickets {
	// 	h.worker.Send(worker.Message{Task: worker.TaskIssueReceipt, TicketID: ticket})
	// 	h.worker.Send(worker.Message{Task: worker.TaskAppendToTracker, TicketID: ticket})
	// }

	for _, ticket := range request.Tickets {

		if ticket.Price.Currency == "" {
			ticket.Price.Currency = "USD"
		}

		switch ticket.Status {
		case "confirmed":
			payload := entities.TicketBookingConfirmed{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
				Status:        ticket.Status,
			}

			payloadEvent, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), []byte(payloadEvent))
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			msg.Metadata.Set("type", "TicketBookingConfirmed")
			// err = h.publisher.Publish("TicketBookingConfirmed", msg)

			log.CorrelationIDFromContext(c.Request().Context())
			err = h.eventBus.Publish(c.Request().Context(), payload)
			if err != nil {
				continue
			}

		case "canceled":
			payload := entities.TicketBookingCanceled{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payloadEvent, err := json.Marshal(payload)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), []byte(payloadEvent))
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			msg.Metadata.Set("type", "TicketBookingCanceled")
			// err = h.publisher.Publish("TicketBookingCanceled", msg)

			log.CorrelationIDFromContext(c.Request().Context())
			err = h.eventBus.Publish(c.Request().Context(), payload)
			if err != nil {
				continue
			}
		default:
			continue
		}

		// ticketReceipt := entities.IssueReceiptPayload{
		// 	TicketID: ticket.TicketID,
		// 	Price:    ticket.Price,
		// }

		// ticketReceiptPayload, err := json.Marshal(ticketReceipt)
		// if err != nil {
		// 	return err
		// }

		// msg := message.NewMessage(watermill.NewUUID(), []byte(ticketReceiptPayload))

		// err = h.publisher.Publish("issue-receipt", msg)
		// if err != nil {
		// 	return err
		// }

		// ticketEvent := entities.AppendToTrackerPayload{
		// 	TicketID:      ticket.TicketID,
		// 	CustomerEmail: ticket.CustomerEmail,
		// 	Price:         ticket.Price,
		// }

		// ticketPayload, err := json.Marshal(ticketEvent)
		// if err != nil {
		// 	continue
		// }

		// msg = message.NewMessage(watermill.NewUUID(), []byte(ticketPayload))
		// err = h.publisher.Publish("append-to-tracker", msg)
		// if err != nil {
		// 	return err
		// }
	}

	return c.NoContent(http.StatusOK)
}
