package api

import (
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("api")
