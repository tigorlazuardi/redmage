package reddit

import (
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("reddit")
