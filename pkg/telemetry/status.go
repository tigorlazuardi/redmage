package telemetry

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// EndWithStatus ends the span with the status of the error if not nil
// otherwise it will set the status to OK.
//
// This function should be used for ending spans, not for starting spans
// or spans that will have children, to avoid duplicate error recordings.
//
// Do not defer this function directly since err might be nil at the
// start of defer call. Instead it should be wrapped in a function to
// capture the error correctly.
//
// Example:
//
//	var err error
//	ctx, span := tracer.Start(ctx, "my-operation")
//	defer func() { telemetry.EndWithStatus(span, err) }()
func EndWithStatus(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	span.End()
}
