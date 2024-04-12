package telemetry

import (
	"context"
	"time"

	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Telemetry struct {
	tracer *sdktrace.TracerProvider
}

func New(ctx context.Context, cfg *config.Config) (tele Telemetry, err error) {
	otel.SetTextMapPropagator(createPropagator())
	provider, err := createProvider(ctx, cfg)
	if err != nil {
		return tele, err
	}
	tele.tracer = provider
	otel.SetTracerProvider(provider)
	return tele, err
}

func createPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}

func createProvider(ctx context.Context, cfg *config.Config) (*sdktrace.TracerProvider, error) {
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.Float64("telemetry.trace.ratio"))),
	}

	if cfg.Bool("telemetry.openobserve.trace.enable") {
		url := cfg.String("telemetry.openobserve.trace.url")
		o2exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpointURL(url),
			otlptracehttp.WithHeaders(map[string]string{
				"Authorization": cfg.String("telemetry.openobserve.trace.auth"),
			}),
		)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		opts = append(opts, sdktrace.WithBatcher(o2exporter))
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("redmage"),
		semconv.ServiceVersionKey.String(cfg.String("runtime.version")),
		attribute.String("environment", cfg.String("runtime.environment")),
	)

	opts = append(opts, sdktrace.WithResource(res))

	return sdktrace.NewTracerProvider(opts...), nil
}

func (te *Telemetry) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return te.tracer.Shutdown(ctx)
}
