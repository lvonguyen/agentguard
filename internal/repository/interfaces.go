// Package repository defines data access interfaces for AgentGuard.
package repository

import (
	"context"

	"github.com/agentguard/agentguard/internal/models"
	"github.com/google/uuid"
)

// ControlRepository defines operations for control framework data.
type ControlRepository interface {
	// Frameworks
	ListFrameworks(ctx context.Context) ([]models.Framework, error)
	GetFramework(ctx context.Context, id string) (*models.Framework, error)
	CreateFramework(ctx context.Context, f *models.Framework) error
	UpdateFramework(ctx context.Context, f *models.Framework) error
	DeleteFramework(ctx context.Context, id string) error

	// Controls
	ListControls(ctx context.Context, frameworkID string) ([]models.Control, error)
	GetControl(ctx context.Context, id string) (*models.Control, error)
	CreateControl(ctx context.Context, c *models.Control) error
	UpdateControl(ctx context.Context, c *models.Control) error
	DeleteControl(ctx context.Context, id string) error

	// Crosswalks
	GetCrosswalk(ctx context.Context, sourceFrameworkID, targetFrameworkID string) ([]models.Crosswalk, error)
	CreateCrosswalk(ctx context.Context, cw *models.Crosswalk) error
	DeleteCrosswalk(ctx context.Context, id string) error
}

// AgentRepository defines operations for agent registry data.
type AgentRepository interface {
	List(ctx context.Context, filters *AgentFilters) ([]models.Agent, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Agent, error)
	Create(ctx context.Context, a *models.Agent) error
	Update(ctx context.Context, a *models.Agent) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetPolicies(ctx context.Context, agentID uuid.UUID) ([]models.Policy, error)
	BindPolicies(ctx context.Context, agentID uuid.UUID, policyIDs []string) error
}

// AgentFilters defines filtering options for agent queries.
type AgentFilters struct {
	Status      *models.AgentStatus
	Environment *string
	Team        *string
	Framework   *string
	Offset      int
	Limit       int
}

// PolicyRepository defines operations for policy data.
type PolicyRepository interface {
	List(ctx context.Context, filters *PolicyFilters) ([]models.Policy, error)
	Get(ctx context.Context, id string) (*models.Policy, error)
	Create(ctx context.Context, p *models.Policy) error
	Update(ctx context.Context, p *models.Policy) error
	Delete(ctx context.Context, id string) error
	GetByType(ctx context.Context, policyType models.PolicyType) ([]models.Policy, error)
}

// PolicyFilters defines filtering options for policy queries.
type PolicyFilters struct {
	Type    *models.PolicyType
	Enabled *bool
	Offset  int
	Limit   int
}

// TraceRepository defines operations for observability trace data.
type TraceRepository interface {
	Create(ctx context.Context, t *models.AgentTrace) error
	Get(ctx context.Context, traceID string) (*models.AgentTrace, error)
	List(ctx context.Context, filters *TraceFilters) ([]models.AgentTrace, error)
	GetSpans(ctx context.Context, traceID string) ([]models.Span, error)
	ListSecuritySignals(ctx context.Context, filters *SignalFilters) ([]models.SecuritySignal, error)
}

// TraceFilters defines filtering options for trace queries.
type TraceFilters struct {
	AgentID   *uuid.UUID
	SessionID *string
	Status    *models.TraceStatus
	StartFrom *int64 // Unix timestamp
	StartTo   *int64
	Offset    int
	Limit     int
}

// SignalFilters defines filtering options for security signal queries.
type SignalFilters struct {
	TraceID  *string
	Type     *models.SignalType
	Severity *string
	Offset   int
	Limit    int
}

// ThreatModelRepository defines operations for threat model data.
type ThreatModelRepository interface {
	List(ctx context.Context) ([]models.ThreatModel, error)
	Get(ctx context.Context, id string) (*models.ThreatModel, error)
	Create(ctx context.Context, tm *models.ThreatModel) error
	Update(ctx context.Context, tm *models.ThreatModel) error
	Delete(ctx context.Context, id string) error
}

// MaturityRepository defines operations for maturity assessment data.
type MaturityRepository interface {
	ListAssessments(ctx context.Context, orgID string) ([]models.MaturityAssessment, error)
	GetAssessment(ctx context.Context, id string) (*models.MaturityAssessment, error)
	CreateAssessment(ctx context.Context, ma *models.MaturityAssessment) error
}

// GapAnalysisRepository defines operations for gap analysis data.
type GapAnalysisRepository interface {
	List(ctx context.Context, orgID string) ([]models.GapAnalysis, error)
	Get(ctx context.Context, id string) (*models.GapAnalysis, error)
	Create(ctx context.Context, ga *models.GapAnalysis) error
}
