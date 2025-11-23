package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "smoke_sensor")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		orderID := string(msg.UUID)
		fmt.Println("Received message", orderID, string(msg.Payload))

		if string(msg.Payload) == "1" {
			err := alarmClient.StartAlarm()
			if err != nil {
				fmt.Println("Failed to start alarm:", err)
				msg.Nack()
				continue
			}
		} else if string(msg.Payload) == "0" {
			err := alarmClient.StopAlarm()
			if err != nil {
				fmt.Println("Failed to stop alarm:", err)
				msg.Nack()
				continue
			}
		}

		msg.Ack()
	}
}
