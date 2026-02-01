package shows

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type ShowRepository struct {
	db *sqlx.DB
}

func NewShowRepository(db *sqlx.DB) *ShowRepository {
	if db == nil {
		panic("db is nil")
	}

	return &ShowRepository{db: db}
}

func (s ShowRepository) AddShow(ctx context.Context, show entities.Show) error {
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO 
			shows (show_id, dead_nation_id, number_of_tickets, start_time, title, venue)
		VALUES (:show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue)
		`, show)

	if err != nil {
		return fmt.Errorf("could not add show: %w", err)
	}

	return nil

}
