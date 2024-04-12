package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Telemetry struct {
	tracer *trace.TracerProvider
}

func New() Telemetry {
	otel.SetTextMapPropagator(createPropagator())
	provider := createProvider()
	otel.SetTracerProvider(provider)
	return Telemetry{
		tracer: provider,
	}
}

func createPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}

func createProvider() *trace.TracerProvider {
	return trace.NewTracerProvider()
}

func (te *Telemetry) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return te.tracer.Shutdown(ctx)
}
