package internal

import (
	"os"

	"github.com/betorvs/plhello/internal/customer"
	"github.com/betorvs/plhello/internal/platform/logger"
	"go.opentelemetry.io/otel/trace"
)

// App stores the current configuration
var App Application

// Application contains the application instantiate structs and configurations
type Application struct {
	// GreetingFormat string
	GreetingFormat string
	// Logger logger.Logger
	Logger logger.Logger
	// Customer customer.Greeting
	Customer customer.Greeting
	// Tracer trace.Tracer
	Tracer trace.Tracer
}

// GetEnv gets an environment variable content or a default value
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func InitApplication(appCustomer customer.Greeting, logger logger.Logger, tracer trace.Tracer) {
	App.GreetingFormat = GetEnv("GREETING_FORMAT", "json")
	App.Logger = logger
	App.Customer = appCustomer
	App.Tracer = tracer
}
