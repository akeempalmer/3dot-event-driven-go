package service

import (
	"context"
	"errors"
	stdHTTP "net/http"

	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
	"tickets/worker"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(
	spreadsheetsAPI worker.SpreadsheetsAPI,
	receiptsService worker.ReceiptsService,
	router *message.Router,
	publisher *redisstream.Publisher,
) Service {

	worker := worker.NewWorker(spreadsheetsAPI, receiptsService)

	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	go worker.Run(context.Background(), router, nil, nil)
	return Service{
		echoRouter: echoRouter,
	}
}

func (s Service) Run(ctx context.Context, router *message.Router) error {

	go func() {
		router.Run(context.Background())
	}()

	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
