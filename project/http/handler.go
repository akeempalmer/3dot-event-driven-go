package http

import (
	"tickets/db/bookings"
	shows "tickets/db/show"
	"tickets/db/tickets"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	// worker *worker.Worker
	// publisher message.Publisher
	eventBus    *cqrs.EventBus
	ticketRepo  *tickets.TicketRepository
	showRepo    *shows.ShowRepository
	bookingRepo *bookings.BookingRepository
}
