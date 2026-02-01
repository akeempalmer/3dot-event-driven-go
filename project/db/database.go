package db

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

	query = `CREATE TABLE IF NOT EXISTS shows (
		show_id UUID PRIMARY KEY,
		dead_nation_id UUID NOT NULL,
		number_of_tickets INT NOT NULL,
		start_time TIMESTAMP NOT NULL,
		title VARCHAR(255) NOT NULL,
		venue VARCHAR(255) NOT NULL,

		UNIQUE (dead_nation_id)
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	query = `CREATE TABLE IF NOT EXISTS bookings (
    booking_id UUID PRIMARY KEY,
    show_id UUID NOT NULL,
    number_of_tickets INT NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    FOREIGN KEY (show_id) REFERENCES shows(show_id)
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return db
}
