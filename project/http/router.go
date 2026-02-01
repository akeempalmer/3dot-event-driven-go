package http

import (
	shows "tickets/db/show"
	"tickets/db/tickets"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	// w *worker.Worker,
	// publisher message.Publisher,
	eventBus *cqrs.EventBus,
	db *sqlx.DB,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		eventBus:   eventBus,
		ticketRepo: tickets.NewTicketRepository(db),
		showRepo:   shows.NewShowRepository(db),
	}

	// e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)
	e.POST("/tickets-status", handler.PostTicketsConfirmation)

	e.GET("/tickets", handler.GetAllTickets)

	e.GET("/health", handler.GetHealthHandler)

	e.POST("/shows", handler.CreateNewShow)

	return e
}

func NewWatermillRouter() *message.Router {
	logger := watermill.NewSlogLogger(nil)
	router := message.NewDefaultRouter(logger)
	return router
}
