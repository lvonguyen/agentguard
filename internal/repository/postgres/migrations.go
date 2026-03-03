package postgres

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

const schemaVersion = 1

var migrations = []struct {
	version     int
	description string
	sql         string
}{
	{
		version:     1,
		description: "initial schema: frameworks, controls, crosswalks, assessments",
		sql: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version     INT PRIMARY KEY,
				description TEXT NOT NULL,
				applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS frameworks (
				id          TEXT PRIMARY KEY,
				name        TEXT NOT NULL,
				version     TEXT NOT NULL,
				publisher   TEXT NOT NULL DEFAULT '',
				description TEXT NOT NULL DEFAULT '',
				url         TEXT NOT NULL DEFAULT '',
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS controls (
				id                TEXT PRIMARY KEY,
				framework_id      TEXT NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
				control_id        TEXT NOT NULL,
				title             TEXT NOT NULL,
				description       TEXT NOT NULL DEFAULT '',
				objectives        JSONB NOT NULL DEFAULT '[]',
				activities        JSONB NOT NULL DEFAULT '[]',
				evidence_types    JSONB NOT NULL DEFAULT '[]',
				applicable_layers JSONB NOT NULL DEFAULT '[]',
				parent_control_id TEXT,
				created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				UNIQUE (framework_id, control_id)
			);

			CREATE INDEX IF NOT EXISTS idx_controls_framework_id ON controls(framework_id);
			CREATE INDEX IF NOT EXISTS idx_controls_control_id ON controls(control_id);

			CREATE TABLE IF NOT EXISTS crosswalks (
				id                    TEXT PRIMARY KEY,
				source_framework_id   TEXT NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
				source_control_id     TEXT NOT NULL,
				target_framework_id   TEXT NOT NULL REFERENCES frameworks(id) ON DELETE CASCADE,
				target_control_id     TEXT NOT NULL,
				mapping_type          TEXT NOT NULL DEFAULT 'related',
				confidence            DOUBLE PRECISION NOT NULL DEFAULT 0.0,
				rationale             TEXT NOT NULL DEFAULT '',
				gaps                  JSONB NOT NULL DEFAULT '[]',
				supplements           JSONB NOT NULL DEFAULT '[]',
				evidence_mapping      JSONB NOT NULL DEFAULT '[]',
				created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE INDEX IF NOT EXISTS idx_crosswalks_source ON crosswalks(source_framework_id, source_control_id);
			CREATE INDEX IF NOT EXISTS idx_crosswalks_target ON crosswalks(target_framework_id, target_control_id);

			CREATE TABLE IF NOT EXISTS assessments (
				id              TEXT PRIMARY KEY,
				organization_id TEXT NOT NULL,
				assessor_id     TEXT NOT NULL DEFAULT '',
				assessment_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				domains         JSONB NOT NULL DEFAULT '[]',
				overall_score   DOUBLE PRECISION NOT NULL DEFAULT 0.0,
				overall_level   INT NOT NULL DEFAULT 1,
				recommendations JSONB NOT NULL DEFAULT '[]',
				created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE INDEX IF NOT EXISTS idx_assessments_org ON assessments(organization_id);

			INSERT INTO schema_migrations (version, description)
			VALUES (1, 'initial schema: frameworks, controls, crosswalks, assessments')
			ON CONFLICT (version) DO NOTHING;
		`,
	},
}

// RunMigrations applies all pending database migrations in order.
func (db *DB) RunMigrations(ctx context.Context) error {
	current, err := db.currentVersion(ctx)
	if err != nil {
		return fmt.Errorf("checking migration version: %w", err)
	}

	applied := 0
	for _, m := range migrations {
		if m.version <= current {
			continue
		}

		log.Info().
			Int("version", m.version).
			Str("description", m.description).
			Msg("applying migration")

		if _, err := db.Pool.Exec(ctx, m.sql); err != nil {
			return fmt.Errorf("applying migration v%d (%s): %w", m.version, m.description, err)
		}

		applied++
	}

	if applied == 0 {
		log.Info().Int("current_version", current).Msg("database schema up to date")
	} else {
		log.Info().Int("applied", applied).Int("target_version", schemaVersion).Msg("migrations complete")
	}

	return nil
}

func (db *DB) currentVersion(ctx context.Context) (int, error) {
	// schema_migrations table may not exist on first run
	var exists bool
	err := db.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'schema_migrations'
		)
	`).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("checking schema_migrations table: %w", err)
	}

	if !exists {
		return 0, nil
	}

	var version int
	err = db.Pool.QueryRow(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("querying current version: %w", err)
	}

	return version, nil
}
