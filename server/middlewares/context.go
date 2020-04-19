package middlewares

import (
	"context"

	"github.com/labstack/echo"

	"github.com/donbattery/bnj/model"
	"github.com/donbattery/bnj/utils"
)

func NewInjectorMiddleware(ctx context.Context, core model.Core) func(echo.HandlerFunc) echo.HandlerFunc {
	cfg := utils.Conf(ctx)
	db := utils.DB(ctx)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", cfg)
			c.Set("database", db)
			c.Set("core", core)
			return next(c)
		}
	}
}
