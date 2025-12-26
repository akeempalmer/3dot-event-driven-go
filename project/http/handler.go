package http

import "github.com/ThreeDotsLabs/watermill/components/cqrs"

type Handler struct {
	// worker *worker.Worker
	// publisher message.Publisher
	publisher *cqrs.EventBus
}
