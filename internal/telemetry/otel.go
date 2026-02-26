// Package telemetry provides OpenTelemetry instrumentation
package telemetry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

// Config holds telemetry configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string
	MetricsPort    int
}

// Provider manages OpenTelemetry providers
type Provider struct {
	config         Config
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	tracer         trace.Tracer
	meter          metric.Meter

	// LLM-specific metrics
	requestCounter  metric.Int64Counter
	requestDuration metric.Float64Histogram
	tokenCounter    metric.Int64Counter
	errorCounter    metric.Int64Counter
	activeRequests  metric.Int64UpDownCounter
}

// NewProvider creates a new telemetry provider
func NewProvider(cfg Config) (*Provider, error) {
	ctx := context.Background()

	// Create resource with service info
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup trace exporter — use TLS by default, plaintext only when OTEL_INSECURE=true
	exporterOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
	}
	if strings.EqualFold(os.Getenv("OTEL_INSECURE"), "true") {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithInsecure())
	} else {
		exporterOpts = append(exporterOpts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	traceExporter, err := otlptracegrpc.New(ctx, exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Setup tracer provider
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Setup Prometheus exporter for metrics
	promExporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(promExporter),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	p := &Provider{
		config:         cfg,
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		tracer:         tracerProvider.Tracer(cfg.ServiceName),
		meter:          meterProvider.Meter(cfg.ServiceName),
	}

	// Initialize metrics
	if err := p.initMetrics(); err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	return p, nil
}

func (p *Provider) initMetrics() error {
	var err error

	p.requestCounter, err = p.meter.Int64Counter(
		"llm_requests_total",
		metric.WithDescription("Total number of LLM requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	p.requestDuration, err = p.meter.Float64Histogram(
		"llm_request_duration_seconds",
		metric.WithDescription("LLM request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	p.tokenCounter, err = p.meter.Int64Counter(
		"llm_tokens_total",
		metric.WithDescription("Total tokens processed"),
		metric.WithUnit("{token}"),
	)
	if err != nil {
		return err
	}

	p.errorCounter, err = p.meter.Int64Counter(
		"llm_errors_total",
		metric.WithDescription("Total LLM errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return err
	}

	p.activeRequests, err = p.meter.Int64UpDownCounter(
		"llm_active_requests",
		metric.WithDescription("Currently active LLM requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return err
	}

	return nil
}

// Tracer returns the tracer instance
func (p *Provider) Tracer() trace.Tracer {
	return p.tracer
}

// Meter returns the meter instance
func (p *Provider) Meter() metric.Meter {
	return p.meter
}

// Shutdown gracefully shuts down telemetry providers.
// Both tracer and meter are shut down regardless of individual failures.
func (p *Provider) Shutdown(ctx context.Context) error {
	var errs []error
	if err := p.tracerProvider.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("tracer provider shutdown: %w", err))
	}
	if err := p.meterProvider.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("meter provider shutdown: %w", err))
	}
	return errors.Join(errs...)
}

// LLMRequestMetrics records metrics for an LLM request
type LLMRequestMetrics struct {
	Provider     string
	Model        string
	InputTokens  int64
	OutputTokens int64
	Duration     time.Duration
	Success      bool
	ErrorType    string
}

// RecordLLMRequest records metrics for an LLM request
func (p *Provider) RecordLLMRequest(ctx context.Context, m LLMRequestMetrics) {
	attrs := []attribute.KeyValue{
		attribute.String("provider", m.Provider),
		attribute.String("model", m.Model),
		attribute.Bool("success", m.Success),
	}

	p.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	p.requestDuration.Record(ctx, m.Duration.Seconds(), metric.WithAttributes(attrs...))

	inputAttrs := make([]attribute.KeyValue, len(attrs), len(attrs)+1)
	copy(inputAttrs, attrs)
	inputAttrs = append(inputAttrs, attribute.String("type", "input"))
	p.tokenCounter.Add(ctx, m.InputTokens, metric.WithAttributes(inputAttrs...))

	outputAttrs := make([]attribute.KeyValue, len(attrs), len(attrs)+1)
	copy(outputAttrs, attrs)
	outputAttrs = append(outputAttrs, attribute.String("type", "output"))
	p.tokenCounter.Add(ctx, m.OutputTokens, metric.WithAttributes(outputAttrs...))

	if !m.Success {
		errAttrs := make([]attribute.KeyValue, len(attrs), len(attrs)+1)
		copy(errAttrs, attrs)
		errAttrs = append(errAttrs, attribute.String("error_type", m.ErrorType))
		p.errorCounter.Add(ctx, 1, metric.WithAttributes(errAttrs...))
	}
}

// StartRequest marks the start of an LLM request
func (p *Provider) StartRequest(ctx context.Context, provider, model string) {
	attrs := []attribute.KeyValue{
		attribute.String("provider", provider),
		attribute.String("model", model),
	}
	p.activeRequests.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// EndRequest marks the end of an LLM request
func (p *Provider) EndRequest(ctx context.Context, provider, model string) {
	attrs := []attribute.KeyValue{
		attribute.String("provider", provider),
		attribute.String("model", model),
	}
	p.activeRequests.Add(ctx, -1, metric.WithAttributes(attrs...))
}

// StartSpan starts a new span
func (p *Provider) StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return p.tracer.Start(ctx, name, opts...)
}
