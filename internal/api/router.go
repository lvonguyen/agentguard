// Package api provides the HTTP API for AgentGuard.
package api

import (
	"crypto/subtle"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/agentguard/agentguard/internal/config"
	"github.com/agentguard/agentguard/internal/controls"
	"github.com/agentguard/agentguard/internal/repository"
	"github.com/agentguard/agentguard/pkg/opa"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// scopeKey is the gin context key for storing JWT scopes.
const scopeKey = "auth_scopes"

// RouterDeps holds dependencies for router initialization.
type RouterDeps struct {
	ControlRepo  repository.ControlRepository
	GapAnalyzer  *controls.GapAnalyzer
	PolicyEngine *opa.Engine
	// StopRateLimiter is set by NewRouter. Call it during graceful shutdown to stop
	// the rate limiter's background cleanup goroutine.
	StopRateLimiter func()
}

// NewRouter creates and configures the HTTP router.
func NewRouter(cfg *config.Config, deps *RouterDeps) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	// Safe default: do not trust any proxy headers (X-Forwarded-For, etc.)
	// Production should configure trusted proxy CIDRs explicitly.
	r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())
	r.Use(securityHeadersMiddleware())
	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1MB
		c.Next()
	})
	r.Use(corsMiddleware(cfg.Server.CORSOrigins))

	// Create handlers with dependencies
	var h *Handlers
	if deps != nil && deps.ControlRepo != nil {
		h = NewHandlers(deps.ControlRepo, deps.GapAnalyzer)
	}

	// Health check
	r.GET("/health", healthCheck)
	r.GET("/ready", makeReadinessCheck(deps))

	// API v1
	rl := newRateLimiter(100, time.Minute)
	// Wire Stop() into deps so callers can halt the cleanup goroutine on shutdown.
	if deps != nil {
		deps.StopRateLimiter = rl.Stop
	}
	v1 := r.Group("/api/v1")
	// Middleware order: Auth → Rate Limiting so that:
	// 1. Unauthenticated requests are rejected before consuming rate limit budget.
	// 2. Rate limits key on bearer identity rather than IP (set after auth validates token).
	v1.Use(bearerTokenMiddleware(cfg.Auth.BearerToken))
	v1.Use(rateLimitMiddleware(rl))
	{
		// Control Framework endpoints
		controls := v1.Group("/controls")
		{
			if h != nil {
				// Use repository-backed handlers
				controls.GET("/frameworks", h.ListFrameworks)
				controls.GET("/frameworks/:id", h.GetFramework)
				controls.GET("/frameworks/:id/controls", h.ListControls)
				controls.GET("/controls/:id", h.GetControl)
				controls.GET("/crosswalk", h.GetCrosswalk)
				writeScope := requireScope(cfg.Auth.Provider, "write:controls")
				controls.POST("/frameworks", writeScope, h.CreateFramework)
				controls.POST("/controls", writeScope, h.CreateControl)
				controls.POST("/gaps/analyze", writeScope, h.AnalyzeGaps)
			} else {
				// Fallback to stub handlers (for testing without DB)
				controls.GET("/frameworks", listFrameworks)
				controls.GET("/frameworks/:id", getFramework)
				controls.GET("/frameworks/:id/controls", listControls)
				controls.GET("/controls/:id", getControl)
				controls.GET("/crosswalk", getCrosswalk)
				controls.POST("/gaps/analyze", requireScope(cfg.Auth.Provider, "write:controls"), analyzeGaps)
			}
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
			sdk.POST("/pre-invoke", makePreInvokeHook(deps))
			sdk.POST("/post-invoke", postInvokeHook)
			sdk.POST("/error", errorHook)
		}
	}

	return r
}

// rateLimiter implements a simple in-memory sliding window rate limiter per IP.
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string][]time.Time
	limit    int
	window   time.Duration
	done     chan struct{}
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		done:     make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

// Stop terminates the cleanup goroutine.
func (rl *rateLimiter) Stop() {
	close(rl.done)
}

func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	timestamps := rl.visitors[key]
	valid := make([]time.Time, 0, len(timestamps))
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= rl.limit {
		rl.visitors[key] = valid
		return false
	}

	rl.visitors[key] = append(valid, now)
	return true
}

func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for {
		select {
		case <-rl.done:
			return
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			cutoff := now.Add(-rl.window)
			for key, timestamps := range rl.visitors {
				valid := make([]time.Time, 0, len(timestamps))
				for _, ts := range timestamps {
					if ts.After(cutoff) {
						valid = append(valid, ts)
					}
				}
				if len(valid) == 0 {
					delete(rl.visitors, key)
				} else {
					rl.visitors[key] = valid
				}
			}
			rl.mu.Unlock()
		}
	}
}

// securityHeadersMiddleware adds security response headers to all responses.
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func rateLimitMiddleware(rl *rateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Key on bearer token identity when present — more accurate for authenticated APIs
		// and allows per-identity rate limits rather than per-IP (which breaks behind NAT).
		key := c.ClientIP()
		if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			if len(token) >= 8 {
				// Use last 8 chars as key suffix to avoid storing full tokens in memory.
				key = "bearer:" + token[len(token)-8:]
			}
		}

		if !rl.allow(key) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

// Middleware

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowed := false
		wildcard := false
		for _, o := range allowedOrigins {
			if o == "*" {
				allowed = true
				wildcard = true
				break
			}
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if wildcard {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
				c.Header("Vary", "Origin")
			}
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func bearerTokenMiddleware(token string) gin.HandlerFunc {
	if token == "" {
		log.Warn().Msg("AUTH_BEARER_TOKEN is not configured — all API requests will be rejected")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		}
	}
	if len(token) < 32 {
		log.Warn().Int("token_len", len(token)).
			Msg("AUTH_BEARER_TOKEN is shorter than 32 chars — consider using a stronger token")
	}
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		provided := strings.TrimPrefix(authHeader, "Bearer ")
		if subtle.ConstantTimeCompare([]byte(provided), []byte(token)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		// Bearer token grants full read+write access — store synthetic scope set.
		c.Set(scopeKey, []string{"read:controls", "write:controls"})
		c.Next()
	}
}

// requireScope returns middleware that enforces the presence of a required scope
// in the request context. In dev mode (auth.provider == "none"), scope checks
// are bypassed. Scopes are populated by the auth middleware upstream.
func requireScope(provider, scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dev mode: skip scope enforcement.
		if strings.EqualFold(provider, "none") {
			c.Next()
			return
		}

		raw, exists := c.Get(scopeKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "missing auth scopes"})
			return
		}

		scopes, ok := raw.([]string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid auth scopes"})
			return
		}

		for _, s := range scopes {
			if s == scope {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error":    "insufficient scope",
			"required": scope,
		})
	}
}

// Health endpoints

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
	})
}

func makeReadinessCheck(deps *RouterDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := gin.H{}
		ready := true

		if deps == nil || deps.ControlRepo == nil {
			checks["database"] = "unavailable"
			ready = false
		} else {
			checks["database"] = "ok"
		}

		if deps == nil || deps.PolicyEngine == nil {
			checks["policy_engine"] = "unavailable"
			ready = false
		} else if !deps.PolicyEngine.Ready() {
			checks["policy_engine"] = "no_policies_loaded"
			ready = false
		} else {
			checks["policy_engine"] = "ok"
		}

		status := http.StatusOK
		statusStr := "ready"
		if !ready {
			status = http.StatusServiceUnavailable
			statusStr = "degraded"
		}

		c.JSON(status, gin.H{
			"status":    statusStr,
			"checks":    checks,
			"timestamp": time.Now().UTC(),
		})
	}
}

// Control Framework handlers

func listFrameworks(c *gin.Context) {
	// TODO: implement — requires database connection
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func getFramework(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented", "id": id})
}

func listControls(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"controls": []any{}, "status": "not_implemented"})
}

func getControl(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getCrosswalk(c *gin.Context) {
	// TODO: implement — requires database connection
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func analyzeGaps(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

// Agent Registry handlers

func listAgents(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"agents": []any{}, "status": "not_implemented"})
}

func registerAgent(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getAgent(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func updateAgent(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func deleteAgent(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getAgentPolicies(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"policies": []any{}, "status": "not_implemented"})
}

func bindAgentPolicies(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

// Observability handlers

func ingestTrace(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func queryTraces(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"traces": []any{}, "status": "not_implemented"})
}

func getTrace(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getTraceSpans(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"spans": []any{}, "status": "not_implemented"})
}

func querySecuritySignals(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"signals": []any{}, "status": "not_implemented"})
}

func getAnomalies(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"anomalies": []any{}, "status": "not_implemented"})
}

func getMetrics(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"metrics": map[string]any{}, "status": "not_implemented"})
}

// Policy handlers

func listPolicies(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"policies": []any{}, "status": "not_implemented"})
}

func createPolicy(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getPolicy(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func updatePolicy(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func deletePolicy(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func validatePolicy(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"valid": false, "status": "not_implemented"})
}

func evaluatePolicy(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"decision": "deny", "status": "not_implemented"})
}

// Threat Model handlers

func listThreatModels(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"models": []any{}, "status": "not_implemented"})
}

func createThreatModel(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getThreatModel(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func updateThreatModel(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func analyzeThreat(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getATLASMapping(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"techniques": []any{}, "status": "not_implemented"})
}

// Maturity Assessment handlers

func listAssessments(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"assessments": []any{}, "status": "not_implemented"})
}

func createAssessment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getAssessment(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getAssessmentReport(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"status": "not_implemented"})
}

func getMaturityModel(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"domains": []any{}, "status": "not_implemented"})
}

func getBenchmarks(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"benchmarks": []any{}, "status": "not_implemented"})
}

// SDK webhook handlers

// makePreInvokeHook returns a handler that evaluates the request against OPA policies.
// Fail-closed: if no policy engine is configured, all requests are denied.
func makePreInvokeHook(deps *RouterDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Fail-closed if policy engine not available
		if deps == nil || deps.PolicyEngine == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"allow":   false,
				"reasons": []string{"policy engine not configured — denying by default"},
			})
			return
		}

		// Limit request body to 1MB to prevent memory exhaustion via large payloads
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)

		// Parse the SDK pre-invoke request body
		var input opa.EvaluationInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"allow":   false,
				"reasons": []string{"invalid request body"},
			})
			return
		}

		// Evaluate against OPA policies
		decision, err := deps.PolicyEngine.Evaluate(c.Request.Context(), "default", &input)
		if err != nil {
			log.Error().Err(err).Msg("policy evaluation failed")
			c.JSON(http.StatusForbidden, gin.H{
				"allow":   false,
				"reasons": []string{"policy evaluation failed — denying by default"},
			})
			return
		}

		c.JSON(http.StatusOK, decision)
	}
}

func postInvokeHook(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "acknowledged"})
}

func errorHook(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"status": "acknowledged"})
}
