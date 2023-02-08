package internal

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Answer defines a json format for a greeting response
type Answer struct {
	Say string `json:"say"`
}

// Greeting function resposible for answering what a customer whats to say
func (a Application) Greeting(c echo.Context) error {
	greetingFormat := "json"
	query := c.QueryParam("format")
	if query != "" {
		a.Logger.Debug("msg", "format query parameter received")
		greetingFormat = strings.ToLower(query)
	}
	hello := a.Customer.Hello()
	// Example of child span
	_, span := a.Tracer.Start(c.Request().Context(), "Greeting", trace.WithAttributes(attribute.String("say", hello), attribute.String("format", greetingFormat)))
	defer span.End()
	res := Answer{
		Say: hello,
	}

	switch greetingFormat {
	case "string":
		return c.String(http.StatusOK, hello)
	default:
		return c.JSON(http.StatusOK, res)
	}
}
