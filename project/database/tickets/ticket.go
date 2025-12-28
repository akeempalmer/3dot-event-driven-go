package tickets

import (
	"context"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type TicketRepo interface {
	Save(ticket entities.TicketBookingConfirmed) error
}

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (tr *TicketRepository) Save(ctx context.Context, ticket *entities.TicketBookingConfirmed) error {
	query := `INSERT INTO tickets (ticket_id, price_amount, price_currency, customer_email) VALUES ($1, $2, $3, $4)`
	_, err := tr.db.Exec(query, ticket.TicketID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail)

	if err != nil {
		return err
	}

	return nil
}
