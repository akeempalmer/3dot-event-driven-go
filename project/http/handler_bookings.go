package http

import (
	"fmt"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) CreateNewBooking(c echo.Context) error {
	var request entities.BookingRequest
	err := c.Bind(&request)
	if err != nil {
		fmt.Println("Error binding request:", err)
		return err
	}

	booking_id := uuid.New().String()

	err = h.bookingRepo.AddBooking(c.Request().Context(), entities.Booking{
		BookingID:       booking_id,
		ShowID:          request.ShowID,
		NumberOfTickets: request.NumberOfTickets,
		CustomerEmail:   request.CustomerEmail,
	})
	if err != nil {
		return err
	}

	return c.JSON(201, map[string]interface{}{"booking_id": booking_id})
}
