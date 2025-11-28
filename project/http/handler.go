package http

import "github.com/ThreeDotsLabs/watermill/message"

type Handler struct {
	// worker *worker.Worker
	publisher message.Publisher
}
