package http

import (
	"net/http"
	"tickets/entities"

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
		msg := message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))

		err = h.publisher.Publish("issue-receipt", msg)
		if err != nil {
			return err
		}

		msg = message.NewMessage(watermill.NewUUID(), []byte(ticket.TicketID))
		err = h.publisher.Publish("append-to-tracker", msg)
		if err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
