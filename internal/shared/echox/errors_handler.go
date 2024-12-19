package echox

import (
	"errors"
	"net/http"
	"story-pulse/contracts"
	"story-pulse/internal/shared/log"

	"github.com/labstack/echo/v4"
	apperrors "story-pulse/internal/shared/error"
)

type HTTPError struct {
	Message    string `json:"message"`
	IncidentID string `json:"incident_id,omitempty"`
}

func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appError *apperrors.Error

	if !errors.As(err, &appError) {
		appError = apperrors.InternalWithoutStackTrace(err)
	}

	httpError := contracts.HTTPError{
		Message:    appError.SafeError(),
		IncidentID: appError.IncidentID,
	}

	logger := log.FromContext(c.Request().Context())

	if appError.Code == apperrors.InternalCode {
		logger.Error("server error",
			"message", err.Error(),
			"incident_id", appError.IncidentID,
			"method", c.Request().Method,
			"url", c.Request().RequestURI,
			"stack_trace", appError.StackTrace,
		)
	} else {
		logger.Warn("client error", "message", err.Error())
	}

	if err := c.JSON(toHTTPStatus(appError.Code), httpError); err != nil {
		c.Logger().Error(err)
	}
}

func toHTTPStatus(code apperrors.Code) int {
	switch code {
	case apperrors.InternalCode:
		return http.StatusInternalServerError
	case apperrors.BadRequestCode:
		return http.StatusBadRequest
	case apperrors.NotFoundCode:
		return http.StatusNotFound
	case apperrors.UnauthorizedCode:
		return http.StatusUnauthorized
	case apperrors.ForbiddenCode:
		return http.StatusForbidden
	case apperrors.AlreadyExistsCode, apperrors.VersionMismatchCode:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
