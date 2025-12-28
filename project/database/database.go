package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitializeSchema() *sqlx.DB {
	query := `CREATE TABLE IF NOT EXISTS tickets (
		ticket_id UUID PRIMARY KEY,
		price_amount NUMERIC(10, 2) NOT NULL,
		price_currency CHAR(3) NOT NULL,
		customer_email VARCHAR(255) NOT NULL
	)`
	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	return db
}
