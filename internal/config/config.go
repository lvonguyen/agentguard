// Package config handles application configuration.
package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration.
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	OPA        OPAConfig        `mapstructure:"opa"`
	OTEL       OTELConfig       `mapstructure:"otel"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Observability ObservabilityConfig `mapstructure:"observability"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            string `mapstructure:"port"`
	Host            string `mapstructure:"host"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
	CORSOrigins     []string `mapstructure:"cors_origins"`
}

// DatabaseConfig holds PostgreSQL configuration.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxConns int    `mapstructure:"max_conns"`
}

// RedisConfig holds Redis configuration.
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// OPAConfig holds Open Policy Agent configuration.
type OPAConfig struct {
	BundlePath    string `mapstructure:"bundle_path"`
	BundleURL     string `mapstructure:"bundle_url"`
	DecisionPath  string `mapstructure:"decision_path"`
	EnableMetrics bool   `mapstructure:"enable_metrics"`
}

// OTELConfig holds OpenTelemetry configuration.
type OTELConfig struct {
	Enabled        bool   `mapstructure:"enabled"`
	Endpoint       string `mapstructure:"endpoint"`
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	SamplingRate   float64 `mapstructure:"sampling_rate"`
}

// AuthConfig holds authentication configuration.
type AuthConfig struct {
	Provider     string   `mapstructure:"provider"` // okta, azure, none
	Issuer       string   `mapstructure:"issuer"`
	Audience     string   `mapstructure:"audience"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	AllowedRoles []string `mapstructure:"allowed_roles"`
}

// ObservabilityConfig holds observability backend configuration.
type ObservabilityConfig struct {
	Langfuse     LangfuseConfig     `mapstructure:"langfuse"`
	ClickHouse   ClickHouseConfig   `mapstructure:"clickhouse"`
}

// LangfuseConfig holds Langfuse integration configuration.
type LangfuseConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	PublicKey string `mapstructure:"public_key"`
	SecretKey string `mapstructure:"secret_key"`
	Host      string `mapstructure:"host"`
}

// ClickHouseConfig holds ClickHouse configuration for time-series data.
type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

// Load reads configuration from file and environment.
func Load(path string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read from config file if provided
	if path != "" {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	} else {
		// Look for config in standard locations
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/agentguard")
		v.AddConfigPath("$HOME/.agentguard")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config: %w", err)
			}
			// Config file not found - continue with defaults and env vars
		}
	}

	// Bind environment variables
	v.SetEnvPrefix("AGENTGUARD")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Override with explicit environment variables
	bindEnvVars(v)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.read_timeout", 15)
	v.SetDefault("server.write_timeout", 15)
	v.SetDefault("server.shutdown_timeout", 30)
	v.SetDefault("server.cors_origins", []string{"*"})

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.database", "agentguard")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_conns", 25)

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)

	// OPA defaults
	v.SetDefault("opa.bundle_path", "./policies/bundle.tar.gz")
	v.SetDefault("opa.decision_path", "agentguard/allow")
	v.SetDefault("opa.enable_metrics", true)

	// OTEL defaults
	v.SetDefault("otel.enabled", true)
	v.SetDefault("otel.service_name", "agentguard")
	v.SetDefault("otel.sampling_rate", 1.0)

	// Auth defaults
	v.SetDefault("auth.provider", "none")

	// Observability defaults
	v.SetDefault("observability.langfuse.enabled", false)
	v.SetDefault("observability.clickhouse.host", "localhost")
	v.SetDefault("observability.clickhouse.port", 9000)
	v.SetDefault("observability.clickhouse.database", "agentguard")
}

func bindEnvVars(v *viper.Viper) {
	// Database credentials from env
	if val := os.Getenv("DATABASE_URL"); val != "" {
		// Parse DATABASE_URL if provided
		v.Set("database.url", val)
	}
	if val := os.Getenv("POSTGRES_USER"); val != "" {
		v.Set("database.user", val)
	}
	if val := os.Getenv("POSTGRES_PASSWORD"); val != "" {
		v.Set("database.password", val)
	}

	// Redis from env
	if val := os.Getenv("REDIS_URL"); val != "" {
		v.Set("redis.url", val)
	}

	// Auth from env
	if val := os.Getenv("OIDC_ISSUER"); val != "" {
		v.Set("auth.issuer", val)
	}
	if val := os.Getenv("OIDC_CLIENT_ID"); val != "" {
		v.Set("auth.client_id", val)
	}
	if val := os.Getenv("OIDC_CLIENT_SECRET"); val != "" {
		v.Set("auth.client_secret", val)
	}
}

// DSN returns the PostgreSQL connection string.
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
	)
}
