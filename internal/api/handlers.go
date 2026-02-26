package api

import (
	"net/http"
	"regexp"

	"github.com/agentguard/agentguard/internal/controls"
	"github.com/agentguard/agentguard/internal/models"
	"github.com/agentguard/agentguard/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var validFrameworkID = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,62}[a-z0-9]$`)

// validID matches a UUID or a slug: non-empty, max 64 chars, safe chars only.
var validID = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,62}[a-zA-Z0-9]$`)

// validateID returns true if id is a valid UUID or slug.
func validateID(id string) bool {
	if _, err := uuid.Parse(id); err == nil {
		return true
	}
	return validID.MatchString(id)
}

// Handlers holds all API handlers with their dependencies.
type Handlers struct {
	ControlRepo repository.ControlRepository
	GapAnalyzer *controls.GapAnalyzer
	// AgentRepo   repository.AgentRepository  // TODO: implement
	// PolicyRepo  repository.PolicyRepository // TODO: implement
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(controlRepo repository.ControlRepository, gapAnalyzer *controls.GapAnalyzer) *Handlers {
	return &Handlers{
		ControlRepo: controlRepo,
		GapAnalyzer: gapAnalyzer,
	}
}

// -----------------------------------------------------------------------------
// Control Framework Handlers
// -----------------------------------------------------------------------------

// ListFrameworks returns all compliance frameworks.
func (h *Handlers) ListFrameworks(c *gin.Context) {
	ctx := c.Request.Context()

	frameworks, err := h.ControlRepo.ListFrameworks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to list frameworks")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list frameworks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"frameworks": frameworks})
}

// GetFramework returns a single framework by ID.
func (h *Handlers) GetFramework(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if !validateID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid framework ID format"})
		return
	}

	framework, err := h.ControlRepo.GetFramework(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get framework")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get framework"})
		return
	}

	if framework == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "framework not found"})
		return
	}

	c.JSON(http.StatusOK, framework)
}

// ListControls returns all controls for a framework.
func (h *Handlers) ListControls(c *gin.Context) {
	ctx := c.Request.Context()
	frameworkID := c.Param("id")

	if !validateID(frameworkID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid framework ID format"})
		return
	}

	controls, err := h.ControlRepo.ListControls(ctx, frameworkID)
	if err != nil {
		log.Error().Err(err).Str("framework_id", frameworkID).Msg("failed to list controls")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list controls"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"framework_id": frameworkID,
		"controls":     controls,
		"count":        len(controls),
	})
}

// GetControl returns a single control by ID.
func (h *Handlers) GetControl(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if !validateID(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid control ID format"})
		return
	}

	control, err := h.ControlRepo.GetControl(ctx, id)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("failed to get control")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get control"})
		return
	}

	if control == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "control not found"})
		return
	}

	c.JSON(http.StatusOK, control)
}

// GetCrosswalk returns crosswalks between two frameworks.
func (h *Handlers) GetCrosswalk(c *gin.Context) {
	ctx := c.Request.Context()
	source := c.Query("source")
	target := c.Query("target")

	if source == "" || target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source and target query parameters required"})
		return
	}

	if !validFrameworkID.MatchString(source) || !validFrameworkID.MatchString(target) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid framework ID format"})
		return
	}

	crosswalks, err := h.ControlRepo.GetCrosswalk(ctx, source, target)
	if err != nil {
		log.Error().Err(err).
			Str("source", source).
			Str("target", target).
			Msg("failed to get crosswalk")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get crosswalk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"source":   source,
		"target":   target,
		"mappings": crosswalks,
		"count":    len(crosswalks),
	})
}

// CreateFramework creates a new framework.
func (h *Handlers) CreateFramework(c *gin.Context) {
	ctx := c.Request.Context()

	var framework models.Framework
	if err := c.ShouldBindJSON(&framework); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if framework.ID == "" {
		framework.ID = uuid.New().String()
	} else if !validFrameworkID.MatchString(framework.ID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid framework ID: must be 2-64 lowercase alphanumeric chars, hyphens, or underscores",
		})
		return
	}

	if err := h.ControlRepo.CreateFramework(ctx, &framework); err != nil {
		log.Error().Err(err).Msg("failed to create framework")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create framework"})
		return
	}

	c.JSON(http.StatusCreated, framework)
}

// CreateControl creates a new control.
func (h *Handlers) CreateControl(c *gin.Context) {
	ctx := c.Request.Context()

	var control models.Control
	if err := c.ShouldBindJSON(&control); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate required fields
	if control.FrameworkID == "" || control.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "framework_id and title are required"})
		return
	}

	if err := h.ControlRepo.CreateControl(ctx, &control); err != nil {
		log.Error().Err(err).Msg("failed to create control")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create control"})
		return
	}

	c.JSON(http.StatusCreated, control)
}

// GapAnalysisRequest represents a gap analysis request.
type GapAnalysisRequest struct {
	TargetFramework     string   `json:"target_framework" binding:"required"`
	ImplementedControls []string `json:"implemented_controls"`
	SourceFramework     string   `json:"source_framework,omitempty"`
}

// AnalyzeGaps analyzes gaps between frameworks.
func (h *Handlers) AnalyzeGaps(c *gin.Context) {
	if h.GapAnalyzer == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "gap analyzer not initialized"})
		return
	}

	var req GapAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	input := &controls.AnalysisInput{
		TargetFramework:     req.TargetFramework,
		ImplementedControls: req.ImplementedControls,
		SourceFramework:     req.SourceFramework,
	}

	output, err := h.GapAnalyzer.RunAnalysis(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Str("framework", req.TargetFramework).Msg("gap analysis failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "analysis failed"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetGapAnalysisSummary returns a summary of gaps for a framework.
func (h *Handlers) GetGapAnalysisSummary(c *gin.Context) {
	if h.GapAnalyzer == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "gap analyzer not initialized"})
		return
	}

	frameworkID := c.Query("framework")
	if frameworkID == "" {
		frameworkID = "nist-ai-rmf" // default
	}

	// Run analysis with no implemented controls to get full gap list
	input := &controls.AnalysisInput{
		TargetFramework:     frameworkID,
		ImplementedControls: []string{},
	}

	output, err := h.GapAnalyzer.RunAnalysis(c.Request.Context(), input)
	if err != nil {
		log.Error().Err(err).Str("framework", frameworkID).Msg("gap summary failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "analysis failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"framework":      output.Framework,
		"framework_name": output.FrameworkName,
		"total_controls": output.TotalControls,
		"summary":        output.Summary,
	})
}
