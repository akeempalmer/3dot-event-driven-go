package event

import (
	"context"
	"tickets/database/tickets"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	ticketRepo      *tickets.TicketRepository
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	db *sqlx.DB,
) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	ticketRepo := tickets.NewTicketRepository(db)

	return Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		ticketRepo:      ticketRepo,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptPayload) error
}
