package http

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
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
		h.pub.Publish("issue-receipt", message.NewMessage(watermill.NewUUID(), []byte(ticket)))
		h.pub.Publish("append-to-tracker", message.NewMessage(watermill.NewUUID(), []byte(ticket)))
	}

	return c.NoContent(http.StatusOK)
}
