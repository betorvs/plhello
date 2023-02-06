package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetHealth(t *testing.T) {
	// Setup
	e := echo.New()
	t.Log("Given a echo instance")
	{
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/health")
		err := CheckHealth(c)
		if err != nil {
			t.Fatalf("\tCheckHealth function return an error %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("\tCheckHealth function should return a 200 status code but returns a %d", rec.Code)
		}
	}

}
