package middlewares

import (
	"context"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
)

type LogHeader struct {
	MessageID string
	Payload   string
	Metadata  message.Metadata
	Handler   context.Context
}

func (lh LogHeader) LoggerMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		logger := slog.With(
			"message_id", msg.UUID,
		)

		logger.Info("Handling a message")
		return next(msg)
	}
}
