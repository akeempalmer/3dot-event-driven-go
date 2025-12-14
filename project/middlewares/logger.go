package middlewares

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v3"
)

type LogHeader struct {
	MessageID string
	Payload   string
	Metadata  message.Metadata
	Handler   context.Context
}

func (lh LogHeader) LoggerMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")

		if correlationID == "" {
			correlationID = fmt.Sprintf("gen_%s", shortuuid.New())
		}

		ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
		msg.SetContext(ctx)

		logger := slog.With(
			"message_id", msg.UUID,
			"metadata", msg.Metadata,
			"handler", msg.Context(),
		)

		logger.Info("Handling a message")
		return next(msg)
	}
}
