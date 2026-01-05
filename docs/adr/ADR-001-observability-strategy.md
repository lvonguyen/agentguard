# ADR-001: Observability Strategy for AI Agent Tracing

**Status:** Accepted  
**Date:** January 2026  
**Author:** Liem Vo-Nguyen  
**Reviewers:** [Security Architecture Team]

---

## Context

AgentGuard requires comprehensive observability for AI agent execution to enable:
1. Security signal detection (injection attempts, tool abuse, data exposure)
2. Audit trail generation for compliance evidence
3. Anomaly detection for behavioral drift
4. Performance monitoring and cost tracking

The market has multiple LLM observability vendors, but none specifically address security use cases.

## Decision Drivers

| Driver | Weight | Description |
|--------|--------|-------------|
| Security signal support | High | Must capture security-relevant events (injection scores, PII detection) |
| Self-hosting option | High | FedRAMP environments require data residency control |
| Open standards (OTEL) | Medium | Avoid vendor lock-in, enable ecosystem integration |
| Agent framework coverage | Medium | Support LangChain, CrewAI, AutoGen at minimum |
| Cost at scale | Medium | 100M+ spans/month at enterprise scale |
| Time-to-value | Low | Can invest in custom development if needed |

## Options Considered

### Option 1: LangSmith (LangChain native)

**Pros:**
- Best-in-class LangChain integration
- Excellent trace visualization UI
- Prompt playground for debugging
- Active development, strong community

**Cons:**
- LangChain ecosystem lock-in
- No self-hosted option (cloud-only)
- No security-specific signals
- No GRC integration
- Closed source

**Verdict:** ❌ Rejected - No self-hosting disqualifies for FedRAMP environments

### Option 2: Langfuse (Open Source)

**Pros:**
- Open source (MIT license)
- Self-hostable on any cloud
- Growing ecosystem, not LangChain-locked
- Cost tracking built-in
- Clean API for extension
- OTEL export support

**Cons:**
- Security features immature
- No compliance mapping
- Smaller team than LangSmith
- Less polished UI

**Verdict:** ✅ Selected as base layer - Self-hosting + extensibility critical

### Option 3: Arize Phoenix (ML Observability)

**Pros:**
- Strong ML observability heritage
- Embedding drift detection
- Open source
- Good visualization

**Cons:**
- LLM features still maturing
- Complex setup (requires Phoenix + online store)
- No agent-specific features
- Limited security focus

**Verdict:** ⚠️ Partial - Integrate embedding drift detection capability only

### Option 4: Helicone (Proxy-Based)

**Pros:**
- Simple proxy setup
- Good latency analytics
- Cost tracking
- Easy onboarding

**Cons:**
- Limited trace depth
- No agent-specific features
- Proxy architecture limits flexibility
- No self-hosting

**Verdict:** ❌ Rejected - Insufficient depth for security analysis

### Option 5: Build from Scratch

**Pros:**
- Full control over data model
- Security-first design
- No vendor dependencies
- Custom retention policies

**Cons:**
- Significant development investment (3-6 months)
- Duplicates commodity functionality
- Maintenance burden
- Delayed time-to-value

**Verdict:** ❌ Rejected - Not efficient, better to extend existing solution

## Decision

**Integrate Langfuse as base observability layer, extend with AgentGuard security enrichment.**

Architecture:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Agent Runtime (LangChain, etc.)              │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              AgentGuard SDK Middleware                   │   │
│  │                                                         │   │
│  │  1. Pre-invoke: Policy check, injection detection       │   │
│  │  2. Wrap callbacks: Trace collection                    │   │
│  │  3. Post-invoke: Security enrichment                    │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                              │
            ┌─────────────────┴─────────────────┐
            │                                   │
            ▼                                   ▼
┌───────────────────────┐          ┌───────────────────────┐
│      Langfuse         │          │   AgentGuard Backend  │
│   (Base Telemetry)    │          │  (Security Enrichment)│
│                       │          │                       │
│ • Trace storage       │   OTEL   │ • Security signals    │
│ • Cost tracking       │ ◄────────│ • Anomaly detection   │
│ • Basic visualization │          │ • Policy violations   │
│ • Prompt management   │          │ • Audit trail         │
└───────────────────────┘          └───────────────────────┘
                                              │
                                              ▼
                                   ┌───────────────────────┐
                                   │     ClickHouse        │
                                   │  (Security Analytics) │
                                   │                       │
                                   │ • Time-series queries │
                                   │ • Anomaly aggregation │
                                   │ • Long-term retention │
                                   └───────────────────────┘
```

**Integration Approach:**

1. **SDK Middleware** wraps agent framework callbacks
2. **Dual emission:** Traces sent to both Langfuse (base telemetry) and AgentGuard (security)
3. **Security enrichment** adds signals not available in Langfuse:
   - Injection attempt scores (from Lakera)
   - PII detection flags
   - Tool abuse indicators
   - Policy violation events
4. **OTEL export** from AgentGuard to Langfuse for unified visualization
5. **ClickHouse** for security-specific analytics and long-term retention

## Consequences

### Positive
- Self-hostable for FedRAMP compliance ✓
- Leverage Langfuse community for base telemetry
- Security-first enrichment layer we control
- Clean separation of concerns (base vs. security)
- OTEL standards enable future flexibility

### Negative
- Two systems to maintain (Langfuse + AgentGuard)
- Potential trace duplication costs
- Langfuse UI won't show security signals (need our own)
- Dependency on Langfuse roadmap for base features

### Risks
| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Langfuse project abandoned | Low | High | OTEL export enables migration |
| Langfuse schema changes break integration | Medium | Medium | Pin versions, integration tests |
| Security enrichment adds latency | Medium | Low | Async processing, sampling |

## Implementation Notes

- Start with LangChain SDK, add CrewAI in Phase 2
- Use Langfuse's Python SDK as foundation
- ClickHouse for security analytics (better than Postgres for time-series)
- OTEL Collector for trace routing flexibility

## References

- [Langfuse Documentation](https://langfuse.com/docs)
- [OpenTelemetry Semantic Conventions for LLM](https://opentelemetry.io/docs/specs/semconv/gen-ai/)
- [ClickHouse for Observability](https://clickhouse.com/docs/en/guides/developer/observability)
