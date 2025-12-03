package event

import "tickets/adapters"

type Handler struct {
	spreadsheetsAPI *adapters.SpreadsheetsAPIClient
	receiptsService *adapters.ReceiptsServiceClient
}
