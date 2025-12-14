package middlewares

import (
	"github.com/ThreeDotsLabs/watermill/message"
)

func SkipPermanentErrorsMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (msgs []*message.Message, err error) {
		if msg.UUID == "2beaf5bc-d5e4-4653-b075-2b36bbf28949" {
			return nil, nil
		}

		// TODO
		// if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
		// 	slog.Error("Invalid message type")
		// 	return nil, nil
		// }

		return next(msg)
	}
}
