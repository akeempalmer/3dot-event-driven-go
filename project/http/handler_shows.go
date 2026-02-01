package http

import (
	"fmt"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) CreateNewShow(c echo.Context) error {
	var request entities.ShowRequest
	err := c.Bind(&request)
	if err != nil {
		fmt.Println("Error binding request:", err)
		return err
	}

	show_id := uuid.New().String()

	err = h.showRepo.AddShow(c.Request().Context(), entities.Show{
		ShowID:          show_id,
		DeadNationID:    request.DeadNationID,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       request.StartTime,
		Title:           request.Title,
		Venue:           request.Venue,
	})
	if err != nil {
		return err
	}

	return c.JSON(201, map[string]interface{}{"show_id": show_id})
}
