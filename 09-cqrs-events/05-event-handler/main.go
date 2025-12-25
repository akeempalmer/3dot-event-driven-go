package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type FollowRequestSent struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type EventsCounter interface {
	CountEvent() error
}

type FollowRequestSentHandler struct {
	EventsCounter
}

func (h FollowRequestSentHandler) CountEvent(ctx context.Context, event *FollowRequestSent) error {
	err := h.EventsCounter.CountEvent()

	return err
}

func NewFollowRequestSentHandler(counter EventsCounter) cqrs.EventHandler {
	h := FollowRequestSentHandler{
		counter,
	}

	return cqrs.NewEventHandler(
		"HandleFollowRequestSent",
		h.CountEvent,
	)
}
