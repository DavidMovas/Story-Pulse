package middlewares

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewLoggerMiddleware(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Infow("Request", "METHOD", c.Request().Method, "URL", c.Request().URL.String(), "QUERY", c.Request().URL.RawQuery)
			return next(c)
		}
	}
}
