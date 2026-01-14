// Package telemetry provides HTTP middleware for observability
package telemetry

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMetrics holds HTTP-specific metrics
type HTTPMetrics struct {
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	requestSize     metric.Int64Histogram
	responseSize    metric.Int64Histogram
}

// NewHTTPMetrics creates HTTP metrics
func NewHTTPMetrics(meter metric.Meter) (*HTTPMetrics, error) {
	m := &HTTPMetrics{}
	var err error

	m.requestCounter, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	m.requestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration"),
	)
	if err != nil {
		return nil, err
	}

	m.requestSize, err = meter.Int64Histogram(
		"http_request_size_bytes",
		metric.WithDescription("HTTP request size"),
	)
	if err != nil {
		return nil, err
	}

	m.responseSize, err = meter.Int64Histogram(
		"http_response_size_bytes",
		metric.WithDescription("HTTP response size"),
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// responseWriter wraps http.ResponseWriter to capture status and size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += int64(n)
	return n, err
}

// Middleware returns HTTP middleware for metrics and tracing
func (m *HTTPMetrics) Middleware(tracer trace.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Start span
			ctx, span := tracer.Start(r.Context(), r.URL.Path,
				trace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.user_agent", r.UserAgent()),
				),
			)
			defer span.End()

			// Wrap response writer
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			// Process request
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Record metrics
			duration := time.Since(start)
			attrs := []attribute.KeyValue{
				attribute.String("method", r.Method),
				attribute.String("path", r.URL.Path),
				attribute.Int("status", rw.status),
			}

			m.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
			m.requestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
			m.responseSize.Record(ctx, rw.size, metric.WithAttributes(attrs...))

			if r.ContentLength > 0 {
				m.requestSize.Record(ctx, r.ContentLength, metric.WithAttributes(attrs...))
			}

			// Add response attributes to span
			span.SetAttributes(
				attribute.Int("http.status_code", rw.status),
				attribute.Int64("http.response_size", rw.size),
			)
		})
	}
}


