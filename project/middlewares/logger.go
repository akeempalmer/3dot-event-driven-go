package middlewares

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
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
		logger := log.FromContext(msg.Context())
		logger.With("message_id", msg.UUID,
			"metadata", msg.Metadata,
			"handler", msg.Context())

		logger.With("message_id", log.CorrelationIDFromContext(msg.Context())).Info("Handling a message")

		return next(msg)
	}
}
