package router

import (
	"net/http"

	"github.com/labstack/echo"
)

func GetUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}
