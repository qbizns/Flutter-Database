package tracing

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Tracer provides distributed tracing capabilities
// This is a stub implementation - can be enhanced with Jaeger/Zipkin
type Tracer struct {
	tracer trace.Tracer
}

// New creates a new tracer instance
func New() *Tracer {
	return &Tracer{
		tracer: trace.NewNoopTracerProvider().Tracer("noop"),
	}
}

// Start creates a new span
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// StartSpan creates a new span (alias for Start for backward compatibility)
func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(ctx, name, opts...)
}

// GetTracer returns the underlying tracer
func (t *Tracer) GetTracer() trace.Tracer {
	return t.tracer
}

// Attribute helpers
var (
	AttrOrganizationID = attribute.Key("organization.id")
	AttrUserID         = attribute.Key("user.id")
)
