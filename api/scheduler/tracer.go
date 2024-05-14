package scheduler

import (
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("scheduler")
