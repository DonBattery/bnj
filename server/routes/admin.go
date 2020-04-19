package routes

import (
	"net/http"

	"github.com/donbattery/bnj/model"
	"github.com/labstack/echo"
)

func Admin(c echo.Context) error {
	cfg := c.Get("config")
	val, ok := cfg.(model.Config)
	if !ok {
		return c.String(http.StatusInternalServerError, "Failed to get config from the context")
	}
	return c.JSON(http.StatusOK, val)
}
