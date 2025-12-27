package message

import (
	"tickets/middlewares"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func useMiddlewares(router *message.Router, watermillLogger watermill.LoggerAdapter) {
	router.AddMiddleware(middleware.Recoverer)

	router.AddMiddleware(middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermillLogger,
	}.Middleware)

	router.AddMiddleware(middlewares.LogHeader{}.CorrelationMiddleware)

	router.AddMiddleware(middlewares.LogHeader{
		MessageID: watermill.NewUUID(),
		Payload:   "Handling a message",
	}.LoggerMiddleware)

	router.AddMiddleware(middlewares.SkipPermanentErrorsMiddleware)
}
