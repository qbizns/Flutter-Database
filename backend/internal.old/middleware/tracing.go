package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/your-org/pos-backend/internal/tracing"
	appctx "github.com/your-org/pos-backend/internal/pkg/context"
)

// TracingMiddleware provides distributed tracing for HTTP requests
type TracingMiddleware struct {
	tracer *tracing.Tracer
}

// NewTracingMiddleware creates a new tracing middleware
func NewTracingMiddleware(tracer *tracing.Tracer) *TracingMiddleware {
	return &TracingMiddleware{tracer: tracer}
}

// Handler returns middleware that creates spans for each request
func (tm *TracingMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming request
			ctx := propagation.TraceContext{}.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

			// Start span
			ctx, span := tm.tracer.StartSpan(ctx, r.Method+" "+r.URL.Path,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethod(r.Method),
					semconv.HTTPRoute(r.URL.Path),
					semconv.HTTPURL(r.URL.String()),
					semconv.HTTPScheme(r.URL.Scheme),
					semconv.HTTPTarget(r.URL.RequestURI()),
					semconv.NetHostName(r.Host),
					semconv.UserAgentOriginal(r.UserAgent()),
				),
			)
			defer span.End()

			// Add organization and user context if available
			if orgID, err := appctx.GetOrganizationID(ctx); err == nil {
				span.SetAttributes(tracing.AttrOrganizationID.String(orgID.String()))
			}
			if userID, err := appctx.GetUserID(ctx); err == nil {
				span.SetAttributes(tracing.AttrUserID.String(userID.String()))
			}

			// Wrap response writer to capture status code
			wrapped := &tracingResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Execute request
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Set response attributes
			span.SetAttributes(semconv.HTTPStatusCode(wrapped.statusCode))

			// Mark span as error if status >= 400
			if wrapped.statusCode >= 400 {
				span.SetStatus(codes.Error, http.StatusText(wrapped.statusCode))
			} else {
				span.SetStatus(codes.Ok, "")
			}
		})
	}
}

// tracingResponseWriter wraps http.ResponseWriter to capture status code
type tracingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (trw *tracingResponseWriter) WriteHeader(code int) {
	trw.statusCode = code
	trw.ResponseWriter.WriteHeader(code)
}

func (trw *tracingResponseWriter) Write(b []byte) (int, error) {
	return trw.ResponseWriter.Write(b)
}
