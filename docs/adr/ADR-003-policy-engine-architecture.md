# ADR-003: Policy Engine Architecture

**Status:** Accepted  
**Date:** 2026-01-03  
**Deciders:** Security Architecture Team  
**Technical Story:** Design policy-as-code engine for agent security guardrails

---

## Context

Agentic AI systems present unique security challenges that traditional content safety guardrails don't address:

1. **Tool access control**: Which tools can an agent invoke, with what parameters, under what conditions?
2. **Data flow governance**: What data can flow between agent components, external systems, and outputs?
3. **Capability escalation**: How do we prevent agents from acquiring capabilities beyond their intended scope?
4. **Human-in-the-loop**: When must an agent pause for human approval?
5. **Rate limiting**: How do we prevent runaway agent loops or resource exhaustion?

### Current Landscape Gaps

| Vendor/Tool | Content Safety | Tool Control | Data Flow | Capability Limits | HITL |
|-------------|----------------|--------------|-----------|-------------------|------|
| Lakera Guard | ✅ Excellent | ❌ None | ❌ None | ❌ None | ❌ None |
| AWS Guardrails | ✅ Good | ❌ None | ❌ None | ❌ None | ❌ None |
| NeMo Guardrails | ✅ Dialogue rails | ⚠️ Limited | ❌ None | ❌ None | ⚠️ Basic |
| Azure Content Safety | ✅ Good | ❌ None | ❌ None | ❌ None | ❌ None |
| Native frameworks | ❌ None | ⚠️ Basic | ❌ None | ❌ None | ⚠️ Manual |

**No vendor provides comprehensive agent-specific policy enforcement.**

---

## Decision Drivers

| Driver | Weight | Description |
|--------|--------|-------------|
| **Agent-specific policies** | Critical | Must address tool access, data flow, capability escalation |
| **Policy-as-code** | High | Version controlled, auditable, testable policies |
| **Low latency** | High | Policy evaluation in agent hot path (<10ms) |
| **Extensibility** | High | Custom policy types as needs evolve |
| **Ecosystem maturity** | Medium | Leverage existing tooling where possible |
| **Operational simplicity** | Medium | Minimize new skills required for policy authors |

---

## Considered Options

### Option 1: Open Policy Agent (OPA) with Rego

**Use OPA as policy engine with custom Rego policies for agent security.**

**Pros:**
- Industry standard for policy-as-code
- Mature ecosystem (5+ years in production at scale)
- Excellent performance (<1ms typical decisions)
- Built-in policy testing framework
- Native JSON/YAML policy data support
- Kubernetes-native integration (if needed)
- Large community and documentation

**Cons:**
- Rego learning curve for non-developers
- Not designed specifically for AI/agent use cases
- Requires custom policy libraries for agent concepts

### Option 2: Cedar (AWS-backed)

**Use Cedar policy language with custom agent extensions.**

**Pros:**
- Modern design, readable syntax
- Strong typing prevents policy errors
- AWS backing provides long-term support
- Designed for authorization scenarios

**Cons:**
- Newer, smaller ecosystem (2023 release)
- Limited tooling compared to OPA
- AWS-centric design may limit portability
- Fewer enterprise deployments to learn from

### Option 3: Custom DSL

**Build AgentGuard-specific policy language tailored to agent concepts.**

**Pros:**
- Perfect fit for agent security domain
- Intuitive syntax for security teams
- No unnecessary complexity from general-purpose engines

**Cons:**
- Significant engineering investment
- No existing tooling (IDE support, testing, etc.)
- Maintenance burden for language evolution
- No community support or external validation

### Option 4: LangChain/Framework-Native Guards

**Use each framework's native callback/guard mechanisms.**

**Pros:**
- Minimal integration overhead
- Framework-specific optimizations
- No additional runtime dependencies

**Cons:**
- Fragmented implementation across frameworks
- No unified policy language
- Limited capabilities (varies by framework)
- No policy versioning or audit trail
- Can't enforce policies at SDK level

### Option 5: Hybrid: OPA Core + AgentGuard Policy Libraries

**Use OPA as engine with pre-built Rego libraries for agent-specific patterns.**

**Pros:**
- Best of OPA ecosystem maturity
- AgentGuard libraries hide Rego complexity
- Policy authors use higher-level abstractions
- Full OPA power available when needed
- Testable policies with existing tooling

**Cons:**
- Two layers of abstraction to maintain
- Some Rego knowledge still needed for advanced cases
- Slight overhead from abstraction layer

---

## Decision

**Selected: Option 5 — Hybrid OPA Core + AgentGuard Policy Libraries**

### Rationale

1. **Ecosystem leverage**: OPA has proven scale (Netflix, Goldman Sachs, Atlassian) and won't be abandoned
2. **Performance**: OPA's compiled policies evaluate in microseconds—essential for agent hot path
3. **Abstraction opportunity**: We build the agent-specific policy patterns, users consume high-level APIs
4. **Testing infrastructure**: OPA's policy testing framework reduces our engineering burden
5. **Extensibility**: Advanced users can drop to Rego for custom needs

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        AGENTGUARD POLICY ENGINE                                 │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                 │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │                   POLICY AUTHORING LAYER                                 │  │
│  │                                                                          │  │
│  │  ┌────────────────────────────────────────────────────────────────────┐ │  │
│  │  │                 AgentGuard Policy YAML                             │ │  │
│  │  │                                                                    │ │  │
│  │  │  # High-level, agent-specific policy definitions                   │ │  │
│  │  │  policies:                                                         │ │  │
│  │  │    - name: restrict-database-tools                                 │ │  │
│  │  │      type: tool_access                                             │ │  │
│  │  │      rules:                                                        │ │  │
│  │  │        - tool: sql_query                                           │ │  │
│  │  │          allow_read: true                                          │ │  │
│  │  │          allow_write: false                                        │ │  │
│  │  │          max_rows: 1000                                            │ │  │
│  │  │          forbidden_tables: [users, credentials]                    │ │  │
│  │  │                                                                    │ │  │
│  │  │    - name: require-approval-for-external                           │ │  │
│  │  │      type: human_in_loop                                           │ │  │
│  │  │      triggers:                                                     │ │  │
│  │  │        - tool_category: external_api                               │ │  │
│  │  │        - data_classification: confidential                         │ │  │
│  │  │      approval:                                                     │ │  │
│  │  │        timeout: 5m                                                 │ │  │
│  │  │        fallback: deny                                              │ │  │
│  │  └────────────────────────────────────────────────────────────────────┘ │  │
│  │                                    │                                     │  │
│  │                                    ▼                                     │  │
│  │  ┌────────────────────────────────────────────────────────────────────┐ │  │
│  │  │               AgentGuard Policy Compiler                           │ │  │
│  │  │                                                                    │ │  │
│  │  │  • Validates policy YAML syntax and semantics                      │ │  │
│  │  │  • Generates Rego code using policy templates                      │ │  │
│  │  │  • Bundles policies with data files                                │ │  │
│  │  │  • Runs OPA policy tests                                           │ │  │
│  │  └────────────────────────────────────────────────────────────────────┘ │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                       │                                        │
│                                       ▼                                        │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │                      REGO POLICY LIBRARIES                               │  │
│  │                                                                          │  │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐ ┌────────────┐            │  │
│  │  │   Tool     │ │  Data Flow │ │  Human-in  │ │   Rate     │            │  │
│  │  │   Access   │ │  Control   │ │   -Loop    │ │  Limiting  │            │  │
│  │  │            │ │            │ │            │ │            │            │  │
│  │  │ • allow/   │ │ • PII      │ │ • Approval │ │ • Per-tool │            │  │
│  │  │   deny     │ │   routing  │ │   gates    │ │   limits   │            │  │
│  │  │ • param    │ │ • Classif  │ │ • Timeout  │ │ • Token    │            │  │
│  │  │   checks   │ │   checks   │ │   handling │ │   budgets  │            │  │
│  │  │ • context  │ │ • Redact   │ │ • Fallback │ │ • Circuit  │            │  │
│  │  │   aware    │ │   rules    │ │   policies │ │   breaker  │            │  │
│  │  └────────────┘ └────────────┘ └────────────┘ └────────────┘            │  │
│  │                                                                          │  │
│  │  ┌────────────┐ ┌────────────┐ ┌────────────┐                           │  │
│  │  │ Capability │ │  Anomaly   │ │   Audit    │                           │  │
│  │  │   Bounds   │ │  Response  │ │  Logging   │                           │  │
│  │  │            │ │            │ │            │                           │  │
│  │  │ • Max tool │ │ • Pattern  │ │ • Decision │                           │  │
│  │  │   calls    │ │   triggers │ │   records  │                           │  │
│  │  │ • Scope    │ │ • Auto     │ │ • Evidence │                           │  │
│  │  │   limits   │ │   response │ │   capture  │                           │  │
│  │  └────────────┘ └────────────┘ └────────────┘                           │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                       │                                        │
│                                       ▼                                        │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │                         OPA RUNTIME                                      │  │
│  │                                                                          │  │
│  │  ┌────────────────────┐  ┌────────────────────┐  ┌──────────────────┐   │  │
│  │  │   Policy Bundle    │  │    Decision API    │  │   Policy Data    │   │  │
│  │  │                    │  │                    │  │                  │   │  │
│  │  │ • Compiled Rego    │  │ • REST endpoint    │  │ • Tool catalog   │   │  │
│  │  │ • Optimized for    │  │ • gRPC endpoint    │  │ • Data classify  │   │  │
│  │  │   <1ms eval        │  │ • SDK embedded     │  │ • User contexts  │   │  │
│  │  │ • Hot reload       │  │                    │  │ • Dynamic rules  │   │  │
│  │  └────────────────────┘  └────────────────────┘  └──────────────────┘   │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Policy Types

### 1. Tool Access Policies

Control which tools agents can invoke and with what parameters.

```yaml
# Example: Database Tool Restrictions
policies:
  - name: database-access-controls
    type: tool_access
    scope:
      agents: ["data-analyst-agent", "report-generator"]
    rules:
      - tool: sql_query
        conditions:
          - operation: SELECT         # Only allow read
          - max_rows: 1000            # Limit result size
          - allowed_schemas: [public, analytics]
          - forbidden_tables: [users, credentials, audit_log]
          - require_where_clause: true  # No full table scans
        
      - tool: file_write
        deny: true                    # Never allow file writes
        
      - tool: http_request
        conditions:
          - allowed_domains: [api.internal.com]
          - methods: [GET]            # Read-only
          - require_auth_header: true
```

### 2. Data Flow Policies

Control how data moves through agent execution chains.

```yaml
# Example: PII Protection
policies:
  - name: pii-protection
    type: data_flow
    rules:
      - classification: PII
        allowed_destinations:
          - internal_database
          - audit_log
        denied_destinations:
          - external_api
          - llm_context        # Don't send PII back to LLM
          - agent_output
        on_violation: redact   # Options: deny, redact, alert
        
      - classification: CONFIDENTIAL
        require_encryption: true
        audit_all_access: true
```

### 3. Human-in-the-Loop Policies

Define when agents must pause for human approval.

```yaml
# Example: High-Risk Action Approval
policies:
  - name: high-risk-approval
    type: human_in_loop
    triggers:
      - tool_category: financial_transaction
      - data_classification: [CONFIDENTIAL, RESTRICTED]
      - cost_estimate: ">$1000"
      - irreversible_action: true
    approval:
      method: slack_dm          # Options: slack, email, portal
      approvers: ["security-team", "manager"]
      timeout: 10m
      require_reason: true
      fallback: deny
      escalation:
        after: 5m
        to: ["director"]
```

### 4. Rate Limiting Policies

Prevent runaway agent execution.

```yaml
# Example: Resource Controls
policies:
  - name: agent-resource-limits
    type: rate_limit
    scope:
      agents: ["*"]            # All agents
    limits:
      - resource: tool_invocations
        max: 100
        window: 1m
        action: throttle
        
      - resource: llm_tokens
        max: 50000
        window: 1h
        action: deny
        alert: true
        
      - resource: execution_time
        max: 5m
        action: terminate
        
      - resource: loop_detection
        max_iterations: 10
        same_tool_sequence: 3
        action: terminate
        alert: true
```

### 5. Capability Bound Policies

Prevent capability escalation beyond intended scope.

```yaml
# Example: Agent Scope Limits
policies:
  - name: agent-capability-bounds
    type: capability
    agents:
      - name: customer-support-agent
        allowed_capabilities:
          - read_customer_data
          - create_ticket
          - send_email
        denied_capabilities:
          - modify_billing
          - access_admin_panel
          - execute_code
        max_chain_depth: 5       # Limit nested agent calls
        allowed_sub_agents: []   # Cannot spawn other agents
```

---

## Implementation

### Phase 1: Core Engine (Weeks 1-4)

1. **OPA integration**: Embed OPA as Go library
2. **Base Rego libraries**: Tool access, rate limiting
3. **Policy compiler**: YAML → Rego generation
4. **SDK integration**: Middleware hooks for policy evaluation

### Phase 2: Advanced Policies (Weeks 5-8)

1. **Data flow policies**: Classification-aware routing
2. **HITL framework**: Approval workflows
3. **Capability bounds**: Scope enforcement
4. **Policy testing**: Automated test generation

### Phase 3: Policy Management (Weeks 9-12)

1. **Portal UI**: Policy editor with validation
2. **Bundle distribution**: OPA bundle server
3. **Policy versioning**: Git-backed policy history
4. **Audit logging**: Decision recording for compliance

---

## Consequences

### Positive

- ✅ **Industry-standard foundation**: OPA is proven at scale
- ✅ **Agent-specific abstractions**: High-level YAML hides Rego complexity
- ✅ **Performance**: Sub-millisecond policy evaluation
- ✅ **Testability**: OPA's policy testing reduces bugs
- ✅ **Extensibility**: Rego available for advanced needs
- ✅ **Audit trail**: All decisions recorded for compliance

### Negative

- ⚠️ **OPA dependency**: Tied to OPA project health (mitigated by large ecosystem)
- ⚠️ **Rego learning curve**: Advanced customization requires Rego knowledge
- ⚠️ **Abstraction maintenance**: Policy compiler adds development burden

### Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| OPA performance degrades with complex policies | Low | High | Pre-compile policies, benchmark in CI |
| Policy language doesn't cover edge case | Medium | Medium | Rego escape hatch for custom needs |
| Policy conflicts across agents | Medium | Medium | Conflict detection in compiler |
| Policy distribution delays | Low | Medium | Local policy cache with TTL |

---

## Appendix: Sample Rego Library

### Tool Access Library

```rego
# agentguard/policies/tool_access.rego

package agentguard.tool_access

import future.keywords.in

default allow = false

# Allow if tool is in agent's allowed list and parameters pass validation
allow {
    tool_allowed
    parameters_valid
    rate_limit_ok
}

tool_allowed {
    input.tool.name in data.policies.allowed_tools[input.agent.id]
}

tool_allowed {
    input.tool.category in data.policies.allowed_categories[input.agent.id]
}

parameters_valid {
    tool_config := data.policies.tool_config[input.tool.name]
    param_checks_pass(tool_config, input.tool.parameters)
}

param_checks_pass(config, params) {
    # Check forbidden patterns
    not contains_forbidden(config.forbidden_patterns, params)
    
    # Check required fields
    all_required_present(config.required_fields, params)
    
    # Check value limits
    within_limits(config.limits, params)
}

# Denial reasons for audit logging
denial_reasons[reason] {
    not tool_allowed
    reason := sprintf("Tool '%s' not allowed for agent '%s'", [input.tool.name, input.agent.id])
}

denial_reasons[reason] {
    not parameters_valid
    reason := sprintf("Invalid parameters for tool '%s'", [input.tool.name])
}
```

### Data Flow Library

```rego
# agentguard/policies/data_flow.rego

package agentguard.data_flow

import future.keywords.in

default allow_flow = false

# Allow data flow if classification allows destination
allow_flow {
    classification := input.data.classification
    destination := input.destination.type
    
    destination in data.policies.allowed_destinations[classification]
}

# Check if redaction is required
requires_redaction {
    input.data.classification == "PII"
    input.destination.type in data.policies.redact_destinations
}

# Generate redaction instructions
redaction_fields[field] {
    input.data.classification == "PII"
    field := data.policies.pii_fields[_]
    field in object.keys(input.data.content)
}
```

---

## References

- [Open Policy Agent Documentation](https://www.openpolicyagent.org/docs/latest/)
- [Rego Language Reference](https://www.openpolicyagent.org/docs/latest/policy-reference/)
- [OPA Performance Best Practices](https://www.openpolicyagent.org/docs/latest/policy-performance/)
- [Cedar Policy Language](https://www.cedarpolicy.com/en)
