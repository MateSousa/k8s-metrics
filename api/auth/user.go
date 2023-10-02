package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateUser
// @Summary Create user
func CreateUser(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// GetUser
// @Summary Get user
func GetUser(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
