package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetHealth(t *testing.T) {
	// Setup
	// "mocking" tracing
	// _ = trace.NewNoopTracerProvider()
	// tracer := otel.Tracer("test")
	// creating echo instance
	e := echo.New()
	// creating a customer
	// customer, _ := customer.NewCustomer("A", "Hello")
	// app := Application{
	// 	Customer: customer,
	// 	Tracer:   tracer,
	// }
	app := Application{}
	t.Log("Given a echo instance")
	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/health")
		err := app.CheckHealth(c)
		if err != nil {
			t.Fatalf("\tCheckHealth function return an error %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("\tCheckHealth function should return a 200 status code but returns a %d", rec.Code)
		}
	}

}
