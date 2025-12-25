package main

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewEventBus(pub message.Publisher) (*cqrs.EventBus, error) {
	eventBus, err := cqrs.NewEventBus(pub, func(eventName string) string { return eventName }, cqrs.JSONMarshaler{})
	if err != nil {
		return nil, fmt.Errorf("could not create event bus: %w", err)
	}

	return eventBus, nil
}
