package middlewares

import (
	"fmt"
	"log/slog"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v3"
)

func (lh LogHeader) CorrelationMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")

		if correlationID == "" {
			correlationID = fmt.Sprintf("gen_%s", shortuuid.New())
		}

		ctx := log.ToContext(msg.Context(), slog.With("correlation_id", correlationID))
		msg.SetContext(ctx)
		return next(msg)
	}
}
