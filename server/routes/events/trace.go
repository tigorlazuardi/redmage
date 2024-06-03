package events

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("server/routes/events")
