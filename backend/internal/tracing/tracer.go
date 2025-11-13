package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/your-org/pos-backend/internal/logging"
	"go.uber.org/zap"
)

// Tracer wraps OpenTelemetry tracer
type Tracer struct {
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
	logger   *logging.Logger
}

// Config holds tracing configuration
type Config struct {
	Enabled      bool
	ServiceName  string
	Environment  string
	JaegerURL    string
	SampleRate   float64
}

// New initializes a new tracer
func New(cfg Config, logger *logging.Logger) (*Tracer, error) {
	if !cfg.Enabled {
		logger.Info("distributed tracing disabled")
		return &Tracer{
			tracer: trace.NewNoopTracerProvider().Tracer("noop"),
			logger: logger,
		}, nil
	}

	// Create Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerURL)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment(cfg.Environment),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	logger.Info("distributed tracing initialized",
		zap.String("service", cfg.ServiceName),
		zap.String("jaeger_url", cfg.JaegerURL),
		zap.Float64("sample_rate", cfg.SampleRate),
	)

	return &Tracer{
		tracer:   tp.Tracer(cfg.ServiceName),
		provider: tp,
		logger:   logger,
	}, nil
}

// Shutdown shuts down the tracer and flushes all spans
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t.provider != nil {
		return t.provider.Shutdown(ctx)
	}
	return nil
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// AddEvent adds an event to the current span
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetAttributes sets attributes on the current span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// Common attribute keys
var (
	AttrOrganizationID = attribute.Key("organization.id")
	AttrUserID         = attribute.Key("user.id")
	AttrDocumentType   = attribute.Key("document.type")
	AttrDocumentID     = attribute.Key("document.id")
	AttrQueryType      = attribute.Key("db.query_type")
	AttrTableName      = attribute.Key("db.table")
	AttrCacheKey       = attribute.Key("cache.key")
	AttrCacheHit       = attribute.Key("cache.hit")
)
