package http

import (
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
)

type Handler struct {
	// worker *worker.Worker
	pub redisstream.Publisher
}
