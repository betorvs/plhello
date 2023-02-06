package internal

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Answer defines a json format for a greeting response
type Answer struct {
	Say string `json:"say"`
}

// Greeting function resposible for answering what a customer whats to say
func Greeting(c echo.Context) error {
	hello := App.Customer.Hello()
	// Example of child span
	_, span := App.Tracer.Start(c.Request().Context(), "Greeting", trace.WithAttributes(attribute.String("say", hello)))
	defer span.End()
	res := Answer{
		Say: hello,
	}
	switch App.GreetingFormat {
	case "string":
		return c.String(http.StatusOK, hello)
	default:
		return c.JSON(http.StatusOK, res)
	}
}
