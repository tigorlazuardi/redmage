package handlers

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("handlers")
