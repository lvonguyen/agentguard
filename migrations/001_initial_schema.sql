-- AgentGuard Initial Schema
-- Migration: 001_initial_schema
-- Description: Create core tables for control frameworks, agents, policies, and observability

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- -----------------------------------------------------------------------------
-- Control Framework Tables
-- -----------------------------------------------------------------------------

CREATE TABLE frameworks (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(32) NOT NULL,
    publisher VARCHAR(255) NOT NULL,
    description TEXT,
    url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(name, version)
);

CREATE INDEX idx_frameworks_name ON frameworks(name);

CREATE TABLE controls (
    id VARCHAR(128) PRIMARY KEY,
    framework_id VARCHAR(64) NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
    control_id VARCHAR(64) NOT NULL,  -- Human-readable ID (e.g., "AC-1")
    title VARCHAR(512) NOT NULL,
    description TEXT,
    objectives JSONB DEFAULT '[]',
    activities JSONB DEFAULT '[]',
    evidence_types JSONB DEFAULT '[]',
    applicable_layers JSONB DEFAULT '[]',
    parent_control_id VARCHAR(128) REFERENCES controls(id) ON DELETE SET NULL,
    UNIQUE(framework_id, control_id)
);

CREATE INDEX idx_controls_framework ON controls(framework_id);
CREATE INDEX idx_controls_parent ON controls(parent_control_id);
CREATE INDEX idx_controls_control_id ON controls(control_id);

CREATE TABLE crosswalks (
    id VARCHAR(128) PRIMARY KEY,
    source_framework_id VARCHAR(64) NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
    source_control_id VARCHAR(128) NOT NULL REFERENCES controls(id) ON DELETE CASCADE,
    target_framework_id VARCHAR(64) NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
    target_control_id VARCHAR(128) NOT NULL REFERENCES controls(id) ON DELETE CASCADE,
    mapping_type VARCHAR(32) NOT NULL CHECK (mapping_type IN ('exact', 'partial', 'superset', 'subset', 'related')),
    confidence NUMERIC(3,2) NOT NULL DEFAULT 0.0 CHECK (confidence >= 0 AND confidence <= 1),
    rationale TEXT,
    gaps JSONB DEFAULT '[]',
    supplements JSONB DEFAULT '[]',
    evidence_mapping JSONB DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_crosswalks_source ON crosswalks(source_framework_id, source_control_id);
CREATE INDEX idx_crosswalks_target ON crosswalks(target_framework_id, target_control_id);

CREATE TABLE gap_analyses (
    id VARCHAR(64) PRIMARY KEY,
    organization_id VARCHAR(64) NOT NULL,
    source_framework_id VARCHAR(64) NOT NULL REFERENCES frameworks(id),
    target_framework_id VARCHAR(64) NOT NULL REFERENCES frameworks(id),
    analysis_date TIMESTAMP NOT NULL DEFAULT NOW(),
    gaps JSONB NOT NULL DEFAULT '[]',
    summary JSONB NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_gap_analyses_org ON gap_analyses(organization_id);

-- -----------------------------------------------------------------------------
-- Agent Registry Tables
-- -----------------------------------------------------------------------------

CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    framework VARCHAR(64) NOT NULL,  -- langchain, crewai, autogen
    version VARCHAR(32),
    owner VARCHAR(255),
    team VARCHAR(255),
    environment VARCHAR(32) NOT NULL CHECK (environment IN ('dev', 'staging', 'prod')),
    capabilities JSONB DEFAULT '[]',
    tools JSONB DEFAULT '[]',
    policies JSONB DEFAULT '[]',  -- Array of policy IDs
    risk_level VARCHAR(32) CHECK (risk_level IN ('low', 'medium', 'high', 'critical')),
    status VARCHAR(32) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended', 'deprecated')),
    last_active_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_environment ON agents(environment);
CREATE INDEX idx_agents_team ON agents(team);
CREATE INDEX idx_agents_updated ON agents(updated_at DESC);

-- -----------------------------------------------------------------------------
-- Policy Tables
-- -----------------------------------------------------------------------------

CREATE TABLE policies (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(32) NOT NULL CHECK (type IN ('tool_access', 'data_flow', 'human_in_loop', 'rate_limit', 'capability')),
    version VARCHAR(32),
    scope JSONB NOT NULL DEFAULT '{}',
    rules JSONB NOT NULL DEFAULT '[]',
    enabled BOOLEAN NOT NULL DEFAULT true,
    priority INTEGER NOT NULL DEFAULT 0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_policies_enabled ON policies(enabled);
CREATE INDEX idx_policies_type ON policies(type);

-- Junction table for agent-policy bindings
CREATE TABLE agent_policies (
    agent_id UUID NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    policy_id VARCHAR(64) NOT NULL REFERENCES policies(id) ON DELETE CASCADE,
    bound_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (agent_id, policy_id)
);

-- -----------------------------------------------------------------------------
-- Observability Tables (Time-Series)
-- -----------------------------------------------------------------------------

CREATE TABLE agent_traces (
    trace_id VARCHAR(64) PRIMARY KEY,
    agent_id UUID NOT NULL REFERENCES agents(id),
    session_id VARCHAR(64),
    user_id VARCHAR(255),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_ms BIGINT,
    status VARCHAR(32) NOT NULL CHECK (status IN ('running', 'completed', 'failed', 'blocked')),
    spans JSONB NOT NULL DEFAULT '[]',
    security_signals JSONB NOT NULL DEFAULT '[]',
    metrics JSONB NOT NULL DEFAULT '{}',
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_traces_agent_time ON agent_traces(agent_id, start_time DESC);
CREATE INDEX idx_traces_status ON agent_traces(status);
CREATE INDEX idx_traces_session ON agent_traces(session_id);

-- Separate table for security signals (for easier querying)
CREATE TABLE security_signals (
    id VARCHAR(64) PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL REFERENCES agent_traces(trace_id) ON DELETE CASCADE,
    span_id VARCHAR(64),
    type VARCHAR(64) NOT NULL,
    severity VARCHAR(32) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    title VARCHAR(512) NOT NULL,
    description TEXT,
    evidence JSONB DEFAULT '{}',
    timestamp TIMESTAMP NOT NULL,
    mitigated BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_signals_trace ON security_signals(trace_id);
CREATE INDEX idx_signals_type_severity ON security_signals(type, severity);
CREATE INDEX idx_signals_timestamp ON security_signals(timestamp DESC);

-- -----------------------------------------------------------------------------
-- Threat Model Tables
-- -----------------------------------------------------------------------------

CREATE TABLE threat_models (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    target_agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    scope TEXT,
    trust_boundaries JSONB DEFAULT '[]',
    threats JSONB DEFAULT '[]',
    mitigations JSONB DEFAULT '[]',
    risk_summary JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_threat_models_agent ON threat_models(target_agent_id);

-- -----------------------------------------------------------------------------
-- Maturity Assessment Tables
-- -----------------------------------------------------------------------------

CREATE TABLE maturity_assessments (
    id VARCHAR(64) PRIMARY KEY,
    organization_id VARCHAR(64) NOT NULL,
    assessor_id VARCHAR(255),
    assessment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    domains JSONB NOT NULL DEFAULT '[]',
    overall_score NUMERIC(4,2),
    overall_level INTEGER CHECK (overall_level >= 1 AND overall_level <= 5),
    recommendations JSONB DEFAULT '[]',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_maturity_org ON maturity_assessments(organization_id);
CREATE INDEX idx_maturity_date ON maturity_assessments(assessment_date DESC);

-- -----------------------------------------------------------------------------
-- Updated timestamp trigger
-- -----------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_frameworks_updated BEFORE UPDATE ON frameworks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_crosswalks_updated BEFORE UPDATE ON crosswalks
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_agents_updated BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_policies_updated BEFORE UPDATE ON policies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tr_threat_models_updated BEFORE UPDATE ON threat_models
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
