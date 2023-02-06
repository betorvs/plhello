package internal

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// health represents the HealthCheck response
type health struct {
	Status string `json:"status"`
}

// CheckHealth handles the application Health Check
func CheckHealth(c echo.Context) error {
	health := health{}
	health.Status = "UP"
	return c.JSON(http.StatusOK, health)
}
