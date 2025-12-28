package http

import (
	"tickets/database/tickets"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	// worker *worker.Worker
	// publisher message.Publisher
	eventBus   *cqrs.EventBus
	ticketRepo *tickets.TicketRepository
}
