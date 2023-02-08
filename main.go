/*
PLHELLO service requires CUSTOMER_NAME and CUSTOMER_GREETING environment variables to start
and it will expose two endpoints:
  - /v1/health that returns 200 ok
  - /v1/greeting that return a greeting from a customer

Others variables:
  - PORT: default 9090
  - APP_NAME: default plhello, use only "_" because it will be used by prometheus and so on.
  - LOG_FORMAT: default logfmt
  - LOG_LEVEL: default INFO
  - TRACE_ENDPOINT: default disabled. Options: stdout or GRPC endpoint

## Directory and File structure

Main purpose of these directories is to provide a package oriented organization. Which means, main contains low-level implementations and it can import from any place. As we move to internal, it cannot import from main, and the same applies for customer and platform packages.

  - main.go: Only file outside internal, it should provide support for all internal functionalities. It should start and shutdown, as an instantiate web server, configurations and so on.
  - internal: Package that provides all functionalities for plhello service.
  - internal/customer: Package for customer and customer behaviours.
  - internal/platform: Directory of packages with external purposes, like logging, databases and so on.
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/betorvs/plhello/internal"
	"github.com/betorvs/plhello/internal/customer"
	"github.com/betorvs/plhello/internal/platform/logger"
	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var build = "develop"

// disabled string constant is used to disable otel tracing provider
const disabled = "disabled"

func main() {
	// Generating application default values
	port := getEnv("PORT", "9090")
	appName := getEnv("APP_NAME", "plhello")
	logFormat := getEnv("LOG_FORMAT", "logfmt")
	logLevel := getEnv("LOG_LEVEL", "INFO")
	log := logger.InitLogger(appName, logFormat, logLevel)
	// Options: stdout for stdout output
	// grpc endpoint
	traceEndpoint := getEnv("TRACE_ENDPOINT", disabled)

	log.Info("msg", "loading service build info "+build)
	// Checking variables for customer
	name := os.Getenv("CUSTOMER_NAME")
	greeting := os.Getenv("CUSTOMER_GREETING")
	appCustomer, err := customer.NewCustomer(name, greeting)
	if err != nil {
		log.Error("msg", "customer with empty values", "err", err.Error())
		os.Exit(1)
	}
	log.Info("msg", "loading customer "+name)

	// starting echo instance
	e := echo.New()
	e.HideBanner = true

	// Init Values.Application with instantiated resources
	// enabling opentelemetry tracing
	tp, err := initTracer(appName, traceEndpoint, log)
	if err != nil {
		log.Error("err", err)
	}
	defer func() {
		if traceEndpoint == disabled {
			return
		}
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Info("msg", "Error shutting down tracer provider", "err", err)
		}
	}()
	tracer := otel.Tracer(appName)
	e.Use(otelecho.Middleware(appName))

	// enabling metrics
	p := prometheus.NewPrometheus(appName, nil)
	p.Use(e)

	// enabling pprof
	pprof.Register(e)

	// exporting customer, log and tracer
	// internal.InitApplication(appCustomer, log, tracer)
	app := &internal.Application{
		Customer: appCustomer,
		Logger:   log,
		Tracer:   tracer,
	}

	g := e.Group("/v1")
	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderContentType},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
	}))

	g.GET("/health", app.CheckHealth)
	g.GET("/greeting", app.Greeting)

	// Start server
	// From echo docs: https://echo.labstack.com/cookbook/graceful-shutdown/
	go func() {
		log.Info("msg", "starting web server on port "+port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			fmt.Println("shutting down the server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.SIGINT (Ctrl+C) and syscall.SIGTERM (K8S shutdown syscall)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Info("msg", "stopping web server")

}

func initTracer(appName, endpoint string, log logger.Logger) (*sdktrace.TracerProvider, error) {
	var tp *sdktrace.TracerProvider
	switch endpoint {
	case disabled:
		// tp = trace.NewNoopTracerProvider()
		// tp2 := sdktrace.NewTracerProvider()
		return tp, nil
	case "stdout":
		log.Info("msg", "starting tracer provider with stdout")
		exporter, err := stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			log.Error("error", err)
			os.Exit(1)
		}

		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(appName))),
		)
	default:
		log.Info("msg", "starting tracer provider with "+endpoint)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
		defer cancel()
		conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Error("error", err)
			os.Exit(1)
		}
		exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			log.Error("error", err)
			os.Exit(1)
		}
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(appName))),
		)
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

// getEnv gets an environment variable content or a default value
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
