package internal

import (
	"github.com/betorvs/plhello/internal/customer"
	"github.com/betorvs/plhello/internal/platform/logger"
	"go.opentelemetry.io/otel/trace"
)

// App stores the current configuration
// var App Application

// Application contains the application instantiate structs and configurations
type Application struct {
	// Logger logger.Logger
	Logger logger.Logger
	// Customer customer.Greeting
	Customer customer.Greeting
	// Tracer trace.Tracer
	Tracer trace.Tracer
}

// func InitApplication(appCustomer customer.Greeting, logger logger.Logger, tracer trace.Tracer) {
// 	App.Logger = logger
// 	App.Customer = appCustomer
// 	App.Tracer = tracer
// }
