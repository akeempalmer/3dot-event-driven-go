package event

import (
	"context"
	"tickets/db/tickets"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
	apiClients      *clients.Clients
	ticketRepo      *tickets.TicketRepository
}

func NewHandler(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	// filesAPI FilesAPI,
	db *sqlx.DB,
	apiClients *clients.Clients,
) Handler {
	if spreadsheetsAPI == nil {
		panic("missing spreadsheetsAPI")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	// if filesAPI == nil {
	// 	panic("missing filesAPI")
	// }

	ticketRepo := tickets.NewTicketRepository(db)

	return Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
		ticketRepo:      ticketRepo,
		apiClients:      apiClients,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptPayload) error
}

type FilesAPI interface {
	UploadFile(ctx context.Context, fileID string, fileContent string) error
}
