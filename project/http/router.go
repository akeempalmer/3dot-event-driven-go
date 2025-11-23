package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	// w *worker.Worker,
	pub *redisstream.Publisher,
) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		pub: *pub,
	}

	e.POST("/tickets-confirmation", handler.PostTicketsConfirmation)

	return e
}
