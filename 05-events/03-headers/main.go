package main

import (
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type MessageHeader struct {
	ID         string `json:"id"`
	EventName  string `json:"event_name"`
	OccurredAt string `json:"occurred_at"`
}

func NewMessageHeader(eventName string) MessageHeader {
	return MessageHeader{
		ID:         uuid.NewString(),
		EventName:  eventName,
		OccurredAt: time.Now().Format(time.RFC3339),
	}
}

type ProductOutOfStock struct {
	ProductID     string        `json:"product_id"`
	MessageHeader MessageHeader `json:"header"`
}

type ProductBackInStock struct {
	ProductID     string        `json:"product_id"`
	MessageHeader MessageHeader `json:"header"`
	Quantity      int           `json:"quantity"`
}

type Publisher struct {
	pub message.Publisher
}

func NewPublisher(pub message.Publisher) Publisher {
	return Publisher{
		pub: pub,
	}
}

func (p Publisher) PublishProductOutOfStock(productID string) error {
	event := ProductOutOfStock{
		ProductID:     productID,
		MessageHeader: NewMessageHeader("ProductOutOfStock"),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}

func (p Publisher) PublishProductBackInStock(productID string, quantity int) error {
	event := ProductBackInStock{
		ProductID:     productID,
		MessageHeader: NewMessageHeader("ProductBackInStock"),
		Quantity:      quantity,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	return p.pub.Publish("product-updates", msg)
}
