package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Test(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello World")
}
