package helper

import "go.opentelemetry.io/otel"

var TracerOtel = otel.Tracer("fiber-server")
