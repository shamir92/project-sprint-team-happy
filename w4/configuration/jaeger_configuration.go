package configuration

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type jaegerConfiguration struct {
	tracer trace.Tracer
}

func NewJaegerConfiguration() *jaegerConfiguration {

	// Initialize JaegerDSN
	var jaegerDsn string = os.Getenv("JAEGER_DSN")
	if jaegerDsn == "" {
		jaegerDsn = "http://localhost:14268/api/traces"
	}
	log.Println(jaegerDsn)

	// Initialize Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerDsn)))
	if err != nil {
		log.Fatalf("failed to initialize Jaeger exporter: %v", err)
	}
	log.Println("shamir")
	log.Println(exporter)

	// Initialize trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("project-sprint-w4"),
		)),
	)
	// defer func() { _ = tp.Shutdown(context.Background()) }()

	otel.SetTracerProvider(tp)
	return &jaegerConfiguration{
		tracer: otel.Tracer("project-sprint-w4"),
	}
}

type IJaegerConfiguration interface {
	GetTracer() trace.Tracer
}

func (ac *jaegerConfiguration) GetTracer() trace.Tracer {
	return ac.tracer
}
