package db_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"tickets/db/tickets"
	"tickets/entities"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB
var getDBOnce sync.Once

func getDB() *sqlx.DB {
	getDBOnce.Do(func() {
		var err error
		db, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}

		query := `CREATE TABLE IF NOT EXISTS tickets (
		ticket_id UUID PRIMARY KEY,
		price_amount NUMERIC(10, 2) NOT NULL,
		price_currency CHAR(3) NOT NULL,
		customer_email VARCHAR(255) NOT NULL
		)`

		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
	})
	return db
}

func TestTicketRepository_Add_Idempotency(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	db := getDB()

	ticketRepo := tickets.NewTicketRepository(db)

	ticket := &entities.TicketBookingConfirmed{
		TicketID:      uuid.NewString(),
		Price:         entities.Money{Amount: "100.00", Currency: "USD"},
		CustomerEmail: "amco@example.com",
	}

	require.EventuallyWithT(t,
		func(collect *assert.CollectT) {
			err := ticketRepo.Save(ctx, ticket)
			if !assert.NoError(collect, err) {
				return
			}
			returnedTickets, err := ticketRepo.FindAll(ctx)
			if !assert.NoError(collect, err) {
				return
			}
			assert.GreaterOrEqual(collect, len(returnedTickets), 1)
		},
		time.Second*10,
		time.Millisecond*50,
	)

	err := ticketRepo.Save(ctx, ticket)

	if err != nil {
		t.Fatal(err)
	}
	assert.NoError(t, err)

	err = ticketRepo.Save(ctx, ticket)

	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)

}

// POSTGRES_URL=postgres://user:password@localhost:5432/db?sslmode=disable go test ./tests/db/ -v
