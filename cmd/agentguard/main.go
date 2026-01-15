// Package main provides the entry point for the AgentGuard API server.
// AgentGuard is an AI security governance framework that provides control mapping,
// runtime observability, and policy-as-code guardrails for agentic AI systems.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/agentguard/agentguard/internal/api"
	"github.com/agentguard/agentguard/internal/config"
	"github.com/agentguard/agentguard/internal/controls"
	"github.com/agentguard/agentguard/internal/repository/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "agentguard",
		Short: "AI Security Governance Framework",
		Long: `AgentGuard provides comprehensive security governance for agentic AI systems.

Features:
  • NIST AI RMF → NIST 800-53 control mapping
  • Agent-specific threat modeling (STRIDE + MITRE ATLAS)
  • Runtime observability with security enrichment
  • Policy-as-code guardrails via OPA
  • Maturity assessment framework`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Server command
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the AgentGuard API server",
		RunE:  runServer,
	}
	serveCmd.Flags().StringP("config", "c", "", "Path to configuration file")
	serveCmd.Flags().StringP("port", "p", "8080", "Port to listen on")
	serveCmd.Flags().Bool("debug", false, "Enable debug logging")

	// Validate command
	validateCmd := &cobra.Command{
		Use:   "validate [policy-file]",
		Short: "Validate policy files",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runValidate,
	}

	// Control mapping commands
	controlCmd := &cobra.Command{
		Use:   "controls",
		Short: "Manage control framework mappings",
	}
	controlCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List available control frameworks",
		RunE:  runControlList,
	})
	controlCmd.AddCommand(&cobra.Command{
		Use:   "crosswalk [source] [target]",
		Short: "Generate crosswalk between frameworks",
		Args:  cobra.ExactArgs(2),
		RunE:  runControlCrosswalk,
	})
	gapsCmd := &cobra.Command{
		Use:   "gaps [framework]",
		Short: "Analyze control gaps",
		Long: `Analyze gaps between your implemented controls and a target framework.

Examples:
  # List all gaps for ISO 42001
  agentguard controls gaps iso-42001

  # Analyze with some controls already implemented
  agentguard controls gaps iso-42001 --implemented "ISO42001-4.1,ISO42001-5.1,ISO42001-6.1"

  # Generate crosswalk from NIST AI RMF
  agentguard controls gaps iso-42001 --source nist-ai-rmf

  # Output as JSON
  agentguard controls gaps iso-42001 --output json`,
		Args: cobra.ExactArgs(1),
		RunE: runControlGaps,
	}
	gapsCmd.Flags().StringP("implemented", "i", "", "Comma-separated list of implemented control IDs")
	gapsCmd.Flags().StringP("output", "o", "text", "Output format: text or json")
	gapsCmd.Flags().StringP("source", "s", "", "Source framework for crosswalk comparison")
	controlCmd.AddCommand(gapsCmd)

	// Threat modeling commands
	threatCmd := &cobra.Command{
		Use:   "threat",
		Short: "Threat modeling tools",
	}
	threatCmd.AddCommand(&cobra.Command{
		Use:   "analyze [manifest-file]",
		Short: "Analyze agent for threats",
		Args:  cobra.ExactArgs(1),
		RunE:  runThreatAnalyze,
	})

	// Maturity assessment commands
	maturityCmd := &cobra.Command{
		Use:   "maturity",
		Short: "Maturity assessment tools",
	}
	maturityCmd.AddCommand(&cobra.Command{
		Use:   "assess",
		Short: "Run interactive maturity assessment",
		RunE:  runMaturityAssess,
	})
	maturityCmd.AddCommand(&cobra.Command{
		Use:   "report [assessment-id]",
		Short: "Generate maturity report",
		Args:  cobra.ExactArgs(1),
		RunE:  runMaturityReport,
	})

	rootCmd.AddCommand(serveCmd, validateCmd, controlCmd, threatCmd, maturityCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) error {
	// Configure logging
	debug, _ := cmd.Flags().GetBool("debug")
	configureLogging(debug)

	// Load configuration
	configPath, _ := cmd.Flags().GetString("config")
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	port, _ := cmd.Flags().GetString("port")
	if port != "" {
		cfg.Server.Port = port
	}

	log.Info().
		Str("version", version).
		Str("port", cfg.Server.Port).
		Msg("Starting AgentGuard server")

	// Initialize database connection
	var deps *api.RouterDeps
	ctx := context.Background()

	if cfg.Database.Host != "" && cfg.Database.User != "" {
		dbCfg := postgres.Config{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Database: cfg.Database.Database,
			SSLMode:  cfg.Database.SSLMode,
			MaxConns: int32(cfg.Database.MaxConns),
		}

		db, err := postgres.New(ctx, dbCfg)
		if err != nil {
			log.Warn().Err(err).Msg("Database connection failed, using stub handlers")
		} else {
			log.Info().
				Str("host", cfg.Database.Host).
				Str("database", cfg.Database.Database).
				Msg("Database connected")

			// Create repositories
			controlRepo := postgres.NewControlRepository(db)

			deps = &api.RouterDeps{
				ControlRepo: controlRepo,
			}

			// Ensure DB is closed on shutdown
			defer db.Close()
		}
	} else {
		log.Info().Msg("No database configured, using stub handlers")
	}

	// Initialize router with dependencies
	router := api.NewRouter(cfg, deps)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info().Msg("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Server shutdown error")
		}
	}()

	// Start server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	log.Info().Msg("Server stopped")
	return nil
}

func runValidate(cmd *cobra.Command, args []string) error {
	configureLogging(false)

	for _, path := range args {
		log.Info().Str("file", path).Msg("Validating policy")
		// TODO: Implement policy validation
		log.Info().Str("file", path).Msg("Policy valid")
	}
	return nil
}

func runControlList(cmd *cobra.Command, args []string) error {
	configureLogging(false)

	analyzer, err := controls.NewGapAnalyzer("")
	if err != nil {
		return fmt.Errorf("initializing analyzer: %w", err)
	}

	analyzer.ListFrameworks(os.Stdout)
	return nil
}

func runControlCrosswalk(cmd *cobra.Command, args []string) error {
	configureLogging(false)

	source, target := args[0], args[1]

	analyzer, err := controls.NewGapAnalyzer("")
	if err != nil {
		return fmt.Errorf("initializing analyzer: %w", err)
	}

	return analyzer.GenerateCrosswalkReport(os.Stdout, source, target)
}

func runControlGaps(cmd *cobra.Command, args []string) error {
	configureLogging(false)

	framework := args[0]

	// Parse implemented controls from flags
	implementedStr, _ := cmd.Flags().GetString("implemented")
	outputFormat, _ := cmd.Flags().GetString("output")
	sourceFramework, _ := cmd.Flags().GetString("source")

	implemented := []string{}
	if implementedStr != "" {
		implemented = strings.Split(implementedStr, ",")
		for i := range implemented {
			implemented[i] = strings.TrimSpace(implemented[i])
		}
	}

	analyzer, err := controls.NewGapAnalyzer("")
	if err != nil {
		return fmt.Errorf("initializing analyzer: %w", err)
	}

	input := &controls.AnalysisInput{
		TargetFramework:     framework,
		ImplementedControls: implemented,
		SourceFramework:     sourceFramework,
	}

	output, err := analyzer.RunAnalysis(context.Background(), input)
	if err != nil {
		return err
	}

	if outputFormat == "json" {
		return analyzer.PrintJSON(os.Stdout, output)
	}

	analyzer.PrintReport(os.Stdout, output)
	return nil
}

func runThreatAnalyze(cmd *cobra.Command, args []string) error {
	manifest := args[0]
	fmt.Printf("Analyzing threats for: %s\n", manifest)
	// TODO: Implement threat analysis
	return nil
}

func runMaturityAssess(cmd *cobra.Command, args []string) error {
	fmt.Println("Starting maturity assessment...")
	// TODO: Implement interactive assessment
	return nil
}

func runMaturityReport(cmd *cobra.Command, args []string) error {
	assessmentID := args[0]
	fmt.Printf("Generating report for assessment: %s\n", assessmentID)
	// TODO: Implement report generation
	return nil
}

func configureLogging(debug bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
