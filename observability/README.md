# LLM Chat Agent Observability Stack

Full observability setup for the LLM Chat Agent, providing metrics, tracing, and logging.

## Components

| Component | Purpose | Port |
|-----------|---------|------|
| OpenTelemetry Collector | Receives and routes telemetry | 4317 (gRPC), 4318 (HTTP) |
| Prometheus | Metrics storage | 9090 |
| Grafana | Dashboards and visualization | 3000 |
| Tempo | Distributed tracing | 3200 |
| Loki | Log aggregation | 3100 |
| Jaeger | Alternative tracing UI | 16686 |

## Quick Start

```bash
# Start the observability stack
cd observability
docker-compose up -d

# Access dashboards
open http://localhost:3000  # Grafana (admin/admin)
open http://localhost:9090  # Prometheus
open http://localhost:16686 # Jaeger
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    LLM Chat Agent                            │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              OpenTelemetry SDK                        │   │
│  │  Traces │ Metrics │ Logs                             │   │
│  └─────────────────────┬────────────────────────────────┘   │
└────────────────────────┼────────────────────────────────────┘
                         │ OTLP (gRPC/HTTP)
                         ▼
┌────────────────────────────────────────────────────────────┐
│              OpenTelemetry Collector                        │
├────────────────────────────────────────────────────────────┤
│  Receivers: OTLP                                           │
│  Processors: Batch, Memory Limiter                         │
│  Exporters: Prometheus, Tempo, Loki                        │
└─────────┬──────────────┬──────────────┬────────────────────┘
          │              │              │
          ▼              ▼              ▼
    ┌──────────┐  ┌──────────┐  ┌──────────┐
    │Prometheus│  │  Tempo   │  │   Loki   │
    │ Metrics  │  │ Traces   │  │   Logs   │
    └────┬─────┘  └────┬─────┘  └────┬─────┘
         │             │             │
         └─────────────┴─────────────┘
                       │
                       ▼
               ┌──────────────┐
               │   Grafana    │
               │  Dashboards  │
               └──────────────┘
```

## Available Dashboards

### LLM Overview Dashboard
- Request rate by provider
- P50/P95/P99 latency
- Token usage (input/output)
- Error rates
- Active requests

### Metrics Available

| Metric | Type | Description |
|--------|------|-------------|
| `llm_requests_total` | Counter | Total LLM requests |
| `llm_request_duration_seconds` | Histogram | Request latency |
| `llm_tokens_total` | Counter | Tokens processed |
| `llm_errors_total` | Counter | Errors by type |
| `llm_active_requests` | Gauge | Currently active requests |
| `http_requests_total` | Counter | HTTP requests |
| `http_request_duration_seconds` | Histogram | HTTP latency |

## Instrumentation

The application is instrumented using OpenTelemetry:

```go
import "github.com/lvonguyen/llm-chat-agent/internal/telemetry"

// Initialize telemetry
provider, err := telemetry.NewProvider(telemetry.Config{
    ServiceName:    "llm-chat-agent",
    ServiceVersion: "1.0.0",
    Environment:    "production",
    OTLPEndpoint:   "localhost:4317",
})

// Record LLM metrics
provider.RecordLLMRequest(ctx, telemetry.LLMRequestMetrics{
    Provider:     "anthropic",
    Model:        "claude-opus-4-5-20250514",
    InputTokens:  500,
    OutputTokens: 1200,
    Duration:     2 * time.Second,
    Success:      true,
})

// Create trace spans
ctx, span := provider.StartSpan(ctx, "process-chat-request")
defer span.End()
```

## Configuration

### Application Environment Variables

```bash
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=llm-chat-agent
OTEL_SERVICE_VERSION=1.0.0
```

### Grafana Access
- URL: http://localhost:3000
- Username: admin
- Password: admin

## Alerting (Future)

Prometheus alerting rules can be added for:
- High error rate (>5%)
- Latency SLO breach (P95 > 5s)
- Token usage anomalies
- Budget threshold alerts

