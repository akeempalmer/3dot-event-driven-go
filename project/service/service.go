package service

import (
	"context"
	"errors"
	"log"
	"log/slog"
	stdHTTP "net/http"

	watermillLog "github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/event"
)

type Service struct {
	echoRouter      *echo.Echo
	watermillRouter *watermillMessage.Router
}

func New(
	spreadsheetsAPI event.SpreadsheetsAPI,
	receiptsService event.ReceiptsService,
	redisClient *redis.Client,
) Service {
	watermillLogger := watermill.NewSlogLogger(slog.Default())

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)

	// Applying the correlation ID decorator
	redisPublisher = watermillLog.CorrelationPublisherDecorator{
		Publisher: redisPublisher,
	}

	watermillRouter := message.NewWatermillRouter(
		receiptsService,
		spreadsheetsAPI,
		redisClient,
		watermillLogger,
	)

	// Create the Event Bus
	eventBus, err := NewEventBus(redisPublisher)
	if err != nil {
		log.Fatalf("could not create event bus: %v", err)
	}

	echoRouter := ticketsHttp.NewHttpRouter(eventBus)

	return Service{
		echoRouter,
		watermillRouter,
	}
}

func NewEventBus(pub watermillMessage.Publisher) (*cqrs.EventBus, error) {
	return cqrs.NewEventBusWithConfig(
		pub,
		cqrs.EventBusConfig{
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				return params.EventName, nil
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
		},
	)
}

func (s Service) Run(ctx context.Context) error {
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")
		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}
		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
