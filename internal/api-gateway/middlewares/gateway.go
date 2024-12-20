package middlewares

import (
	"github.com/labstack/echo/v4"
)

func NewGatewayMiddleware(gatewayHandler echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return gatewayHandler(c)
		}
	}
}
