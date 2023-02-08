package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/betorvs/plhello/internal/customer"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TestGetGreeting(t *testing.T) {
	// Setup
	// "mocking" tracing
	_ = trace.NewNoopTracerProvider()
	tracer := otel.Tracer("test")
	// creating echo instance
	e := echo.New()
	// creating a customer
	customer, _ := customer.NewCustomer("A", "Hello")
	app := Application{
		Customer: customer,
		Tracer:   tracer,
	}
	t.Log("Given a echo instance and a valid customer")
	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/greeting")
		err := app.Greeting(c)
		if err != nil {
			t.Fatalf("\tGreeting function return an error %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("\tGreeting function should return a 200 status code but returns a %d", rec.Code)
		}
	}
	t.Log("Given a echo instance and a valid customer and a query parameter")
	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/greeting?format=string")
		c.QueryParams().Add("format", "string")
		err := app.Greeting(c)
		if err != nil {
			t.Fatalf("\tGreeting function return an error %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("\tGreeting function should return a 200 status code but returns a %d", rec.Code)
		}
	}

}
