package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetHealthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
