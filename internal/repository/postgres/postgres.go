// Package postgres implements PostgreSQL repositories for AgentGuard.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Config holds PostgreSQL connection configuration.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	MaxConns int32
}

// DB wraps the PostgreSQL connection pool.
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new PostgreSQL connection pool.
// Uses struct-based config to avoid embedding credentials in the DSN string,
// which would leak passwords in error messages and log output.
func New(ctx context.Context, cfg Config) (*DB, error) {
	// Build DSN without password — set password via struct field to keep it
	// out of error-path string representations.
	dsn := fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)

	if cfg.MaxConns == 0 {
		cfg.MaxConns = 25
	}

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing connection config: %w", err)
	}

	// Set password via struct field — never appears in DSN string or error messages.
	poolCfg.ConnConfig.Password = cfg.Password

	// Connection pool settings
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	log.Info().
		Str("host", cfg.Host).
		Int("port", cfg.Port).
		Str("database", cfg.Database).
		Msg("PostgreSQL connection established")

	return &DB{Pool: pool}, nil
}

// Close closes the connection pool.
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Info().Msg("PostgreSQL connection closed")
	}
}

// Health checks if the database connection is healthy.
func (db *DB) Health(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database pool not initialized")
	}
	return db.Pool.Ping(ctx)
}

// Ping is an alias for Health for interface compatibility.
func (db *DB) Ping(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database pool not initialized")
	}
	return db.Pool.Ping(ctx)
}

// WithTx executes a function within a transaction.
func (db *DB) WithTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			log.Error().Err(rbErr).Msg("failed to rollback transaction")
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			log.Error().Err(rbErr).Msg("failed to rollback after commit failure")
		}
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
