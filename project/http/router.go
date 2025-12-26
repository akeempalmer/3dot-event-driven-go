package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	// w *worker.Worker,
	// publisher message.Publisher,
	eventBus *cqrs.EventBus,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		publisher: eventBus,
	}

	// e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)
	e.POST("/tickets-status", handler.PostTicketsConfirmation)

	e.GET("/health", handler.GetHealthHandler)

	return e
}

func NewWatermillRouter() *message.Router {
	logger := watermill.NewSlogLogger(nil)
	router := message.NewDefaultRouter(logger)
	return router
}
