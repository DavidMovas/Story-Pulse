package echox

import (
	"brain-wave/internal/shared/log"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestGroup := slog.Group("request",
			slog.String("method", c.Request().Method),
			slog.String("url", c.Request().URL.String()),
		)
		attrs := []any{requestGroup}

		logger := slog.Default().With(attrs...)
		ctx := log.WithLogger(c.Request().Context(), logger)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
