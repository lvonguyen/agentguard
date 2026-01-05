// Package api provides the HTTP API for AgentGuard.
package api

import (
	"net/http"
	"time"

	"github.com/agentguard/agentguard/internal/config"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the HTTP router.
func NewRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg.Server.CORSOrigins))
	r.Use(requestLogger())

	// Health check
	r.GET("/health", healthCheck)
	r.GET("/ready", readinessCheck)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Control Framework endpoints
		controls := v1.Group("/controls")
		{
			controls.GET("/frameworks", listFrameworks)
			controls.GET("/frameworks/:id", getFramework)
			controls.GET("/frameworks/:id/controls", listControls)
			controls.GET("/controls/:id", getControl)
			controls.GET("/crosswalk", getCrosswalk)
			controls.POST("/gaps/analyze", analyzeGaps)
		}

		// Agent Registry endpoints
		agents := v1.Group("/agents")
		{
			agents.GET("", listAgents)
			agents.POST("", registerAgent)
			agents.GET("/:id", getAgent)
			agents.PUT("/:id", updateAgent)
			agents.DELETE("/:id", deleteAgent)
			agents.GET("/:id/policies", getAgentPolicies)
			agents.PUT("/:id/policies", bindAgentPolicies)
		}

		// Observability endpoints
		observe := v1.Group("/observe")
		{
			observe.POST("/traces", ingestTrace)
			observe.GET("/traces", queryTraces)
			observe.GET("/traces/:id", getTrace)
			observe.GET("/traces/:id/spans", getTraceSpans)
			observe.GET("/signals", querySecuritySignals)
			observe.GET("/anomalies", getAnomalies)
			observe.GET("/metrics", getMetrics)
		}

		// Policy endpoints
		policies := v1.Group("/policies")
		{
			policies.GET("", listPolicies)
			policies.POST("", createPolicy)
			policies.GET("/:id", getPolicy)
			policies.PUT("/:id", updatePolicy)
			policies.DELETE("/:id", deletePolicy)
			policies.POST("/validate", validatePolicy)
			policies.POST("/evaluate", evaluatePolicy)
		}

		// Threat Model endpoints
		threats := v1.Group("/threats")
		{
			threats.GET("/models", listThreatModels)
			threats.POST("/models", createThreatModel)
			threats.GET("/models/:id", getThreatModel)
			threats.PUT("/models/:id", updateThreatModel)
			threats.POST("/analyze", analyzeThreat)
			threats.GET("/atlas", getATLASMapping)
		}

		// Maturity Assessment endpoints
		maturity := v1.Group("/maturity")
		{
			maturity.GET("/assessments", listAssessments)
			maturity.POST("/assessments", createAssessment)
			maturity.GET("/assessments/:id", getAssessment)
			maturity.GET("/assessments/:id/report", getAssessmentReport)
			maturity.GET("/model", getMaturityModel)
			maturity.GET("/benchmarks", getBenchmarks)
		}

		// SDK webhook endpoints (for agent middleware callbacks)
		sdk := v1.Group("/sdk")
		{
			sdk.POST("/pre-invoke", preInvokeHook)
			sdk.POST("/post-invoke", postInvokeHook)
			sdk.POST("/error", errorHook)
		}
	}

	return r
}

// Middleware

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		
		// Log request (in production, use structured logging)
		_ = time.Since(start)
	}
}

// Health endpoints

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
	})
}

func readinessCheck(c *gin.Context) {
	// TODO: Check database, redis, OPA connectivity
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC(),
	})
}

// Control Framework handlers

func listFrameworks(c *gin.Context) {
	frameworks := []map[string]any{
		{"id": "nist-ai-rmf", "name": "NIST AI Risk Management Framework", "version": "1.0"},
		{"id": "nist-800-53", "name": "NIST SP 800-53 Rev 5", "version": "5.0"},
		{"id": "iso-42001", "name": "ISO/IEC 42001:2023", "version": "2023"},
	}
	c.JSON(http.StatusOK, gin.H{"frameworks": frameworks})
}

func getFramework(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "status": "not_implemented"})
}

func listControls(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"controls": []any{}, "status": "not_implemented"})
}

func getControl(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func getCrosswalk(c *gin.Context) {
	source := c.Query("source")
	target := c.Query("target")
	c.JSON(http.StatusOK, gin.H{
		"source": source,
		"target": target,
		"mappings": []any{},
		"status": "not_implemented",
	})
}

func analyzeGaps(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

// Agent Registry handlers

func listAgents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"agents": []any{}, "status": "not_implemented"})
}

func registerAgent(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "not_implemented"})
}

func getAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func updateAgent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func deleteAgent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func getAgentPolicies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"policies": []any{}, "status": "not_implemented"})
}

func bindAgentPolicies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

// Observability handlers

func ingestTrace(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "not_implemented"})
}

func queryTraces(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"traces": []any{}, "status": "not_implemented"})
}

func getTrace(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func getTraceSpans(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"spans": []any{}, "status": "not_implemented"})
}

func querySecuritySignals(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"signals": []any{}, "status": "not_implemented"})
}

func getAnomalies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"anomalies": []any{}, "status": "not_implemented"})
}

func getMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"metrics": map[string]any{}, "status": "not_implemented"})
}

// Policy handlers

func listPolicies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"policies": []any{}, "status": "not_implemented"})
}

func createPolicy(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "not_implemented"})
}

func getPolicy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func updatePolicy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func deletePolicy(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func validatePolicy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"valid": true, "status": "not_implemented"})
}

func evaluatePolicy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"decision": "allow", "status": "not_implemented"})
}

// Threat Model handlers

func listThreatModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"models": []any{}, "status": "not_implemented"})
}

func createThreatModel(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "not_implemented"})
}

func getThreatModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func updateThreatModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func analyzeThreat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func getATLASMapping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"techniques": []any{}, "status": "not_implemented"})
}

// Maturity Assessment handlers

func listAssessments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"assessments": []any{}, "status": "not_implemented"})
}

func createAssessment(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"status": "not_implemented"})
}

func getAssessment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func getAssessmentReport(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "not_implemented"})
}

func getMaturityModel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"domains": []any{}, "status": "not_implemented"})
}

func getBenchmarks(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"benchmarks": []any{}, "status": "not_implemented"})
}

// SDK webhook handlers

func preInvokeHook(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"allow": true,
		"decisions": []any{},
	})
}

func postInvokeHook(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "acknowledged"})
}

func errorHook(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "acknowledged"})
}
