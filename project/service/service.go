package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	watermillLog "github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
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
	db *sqlx.DB,
	apiClients *clients.Clients,
) Service {
	watermillLogger := watermill.NewSlogLogger(slog.Default())

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)

	// Applying the correlation ID decorator
	redisPublisher = watermillLog.CorrelationPublisherDecorator{
		Publisher: redisPublisher,
	}

	eventBus := event.NewBus(redisPublisher)

	eventsHandler := event.NewHandler(
		spreadsheetsAPI,
		receiptsService,
		db,
		apiClients,
	)

	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := message.NewWatermillRouter(
		eventProcessorConfig,
		eventsHandler,
		watermillLogger,
	)

	echoRouter := ticketsHttp.NewHttpRouter(eventBus, db)

	return Service{
		echoRouter,
		watermillRouter,
	}
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
