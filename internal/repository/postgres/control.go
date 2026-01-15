package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/agentguard/agentguard/internal/models"
	"github.com/jackc/pgx/v5"
)

// ControlRepository implements repository.ControlRepository for PostgreSQL.
type ControlRepository struct {
	db *DB
}

// NewControlRepository creates a new ControlRepository.
func NewControlRepository(db *DB) *ControlRepository {
	return &ControlRepository{db: db}
}

// -----------------------------------------------------------------------------
// Framework Operations
// -----------------------------------------------------------------------------

// ListFrameworks returns all frameworks.
func (r *ControlRepository) ListFrameworks(ctx context.Context) ([]models.Framework, error) {
	query := `
		SELECT id, name, version, publisher, description, url, created_at, updated_at
		FROM frameworks
		ORDER BY name, version`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("querying frameworks: %w", err)
	}
	defer rows.Close()

	var frameworks []models.Framework
	for rows.Next() {
		var f models.Framework
		if err := rows.Scan(
			&f.ID, &f.Name, &f.Version, &f.Publisher,
			&f.Description, &f.URL, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning framework: %w", err)
		}
		frameworks = append(frameworks, f)
	}

	return frameworks, rows.Err()
}

// GetFramework returns a framework by ID.
func (r *ControlRepository) GetFramework(ctx context.Context, id string) (*models.Framework, error) {
	query := `
		SELECT id, name, version, publisher, description, url, created_at, updated_at
		FROM frameworks
		WHERE id = $1`

	var f models.Framework
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&f.ID, &f.Name, &f.Version, &f.Publisher,
		&f.Description, &f.URL, &f.CreatedAt, &f.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting framework %s: %w", id, err)
	}

	return &f, nil
}

// CreateFramework creates a new framework.
func (r *ControlRepository) CreateFramework(ctx context.Context, f *models.Framework) error {
	query := `
		INSERT INTO frameworks (id, name, version, publisher, description, url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := r.db.Pool.Exec(ctx, query,
		f.ID, f.Name, f.Version, f.Publisher, f.Description, f.URL,
	)
	if err != nil {
		return fmt.Errorf("creating framework: %w", err)
	}

	return nil
}

// UpdateFramework updates an existing framework.
func (r *ControlRepository) UpdateFramework(ctx context.Context, f *models.Framework) error {
	query := `
		UPDATE frameworks
		SET name = $2, version = $3, publisher = $4, description = $5, url = $6
		WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query,
		f.ID, f.Name, f.Version, f.Publisher, f.Description, f.URL,
	)
	if err != nil {
		return fmt.Errorf("updating framework: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("framework %s not found", f.ID)
	}

	return nil
}

// DeleteFramework deletes a framework by ID.
func (r *ControlRepository) DeleteFramework(ctx context.Context, id string) error {
	query := `DELETE FROM frameworks WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting framework: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("framework %s not found", id)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Control Operations
// -----------------------------------------------------------------------------

// ListControls returns all controls for a framework.
func (r *ControlRepository) ListControls(ctx context.Context, frameworkID string) ([]models.Control, error) {
	query := `
		SELECT id, framework_id, control_id, title, description,
		       objectives, activities, evidence_types, applicable_layers, parent_control_id
		FROM controls
		WHERE framework_id = $1
		ORDER BY control_id`

	rows, err := r.db.Pool.Query(ctx, query, frameworkID)
	if err != nil {
		return nil, fmt.Errorf("querying controls: %w", err)
	}
	defer rows.Close()

	var controls []models.Control
	for rows.Next() {
		var c models.Control
		var objectives, activities, evidenceTypes, applicableLayers []byte

		if err := rows.Scan(
			&c.ID, &c.FrameworkID, &c.ControlID, &c.Title, &c.Description,
			&objectives, &activities, &evidenceTypes, &applicableLayers, &c.ParentControlID,
		); err != nil {
			return nil, fmt.Errorf("scanning control: %w", err)
		}

		// Unmarshal JSONB arrays
		if err := json.Unmarshal(objectives, &c.Objectives); err != nil {
			c.Objectives = []string{}
		}
		if err := json.Unmarshal(activities, &c.Activities); err != nil {
			c.Activities = []string{}
		}
		if err := json.Unmarshal(evidenceTypes, &c.EvidenceTypes); err != nil {
			c.EvidenceTypes = []string{}
		}
		if err := json.Unmarshal(applicableLayers, &c.ApplicableLayers); err != nil {
			c.ApplicableLayers = []string{}
		}

		controls = append(controls, c)
	}

	return controls, rows.Err()
}

// GetControl returns a control by ID.
func (r *ControlRepository) GetControl(ctx context.Context, id string) (*models.Control, error) {
	query := `
		SELECT id, framework_id, control_id, title, description,
		       objectives, activities, evidence_types, applicable_layers, parent_control_id
		FROM controls
		WHERE id = $1`

	var c models.Control
	var objectives, activities, evidenceTypes, applicableLayers []byte

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.FrameworkID, &c.ControlID, &c.Title, &c.Description,
		&objectives, &activities, &evidenceTypes, &applicableLayers, &c.ParentControlID,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting control %s: %w", id, err)
	}

	// Unmarshal JSONB arrays
	json.Unmarshal(objectives, &c.Objectives)
	json.Unmarshal(activities, &c.Activities)
	json.Unmarshal(evidenceTypes, &c.EvidenceTypes)
	json.Unmarshal(applicableLayers, &c.ApplicableLayers)

	return &c, nil
}

// CreateControl creates a new control.
func (r *ControlRepository) CreateControl(ctx context.Context, c *models.Control) error {
	objectives, _ := json.Marshal(c.Objectives)
	activities, _ := json.Marshal(c.Activities)
	evidenceTypes, _ := json.Marshal(c.EvidenceTypes)
	applicableLayers, _ := json.Marshal(c.ApplicableLayers)

	query := `
		INSERT INTO controls (id, framework_id, control_id, title, description,
		                      objectives, activities, evidence_types, applicable_layers, parent_control_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Pool.Exec(ctx, query,
		c.ID, c.FrameworkID, c.ControlID, c.Title, c.Description,
		objectives, activities, evidenceTypes, applicableLayers, c.ParentControlID,
	)
	if err != nil {
		return fmt.Errorf("creating control: %w", err)
	}

	return nil
}

// UpdateControl updates an existing control.
func (r *ControlRepository) UpdateControl(ctx context.Context, c *models.Control) error {
	objectives, _ := json.Marshal(c.Objectives)
	activities, _ := json.Marshal(c.Activities)
	evidenceTypes, _ := json.Marshal(c.EvidenceTypes)
	applicableLayers, _ := json.Marshal(c.ApplicableLayers)

	query := `
		UPDATE controls
		SET framework_id = $2, control_id = $3, title = $4, description = $5,
		    objectives = $6, activities = $7, evidence_types = $8, applicable_layers = $9, parent_control_id = $10
		WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query,
		c.ID, c.FrameworkID, c.ControlID, c.Title, c.Description,
		objectives, activities, evidenceTypes, applicableLayers, c.ParentControlID,
	)
	if err != nil {
		return fmt.Errorf("updating control: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("control %s not found", c.ID)
	}

	return nil
}

// DeleteControl deletes a control by ID.
func (r *ControlRepository) DeleteControl(ctx context.Context, id string) error {
	query := `DELETE FROM controls WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting control: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("control %s not found", id)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Crosswalk Operations
// -----------------------------------------------------------------------------

// GetCrosswalk returns crosswalks between two frameworks.
func (r *ControlRepository) GetCrosswalk(ctx context.Context, sourceFrameworkID, targetFrameworkID string) ([]models.Crosswalk, error) {
	query := `
		SELECT id, source_framework_id, source_control_id, target_framework_id, target_control_id,
		       mapping_type, confidence, rationale, gaps, supplements, evidence_mapping, created_at, updated_at
		FROM crosswalks
		WHERE source_framework_id = $1 AND target_framework_id = $2
		ORDER BY source_control_id`

	rows, err := r.db.Pool.Query(ctx, query, sourceFrameworkID, targetFrameworkID)
	if err != nil {
		return nil, fmt.Errorf("querying crosswalks: %w", err)
	}
	defer rows.Close()

	var crosswalks []models.Crosswalk
	for rows.Next() {
		var cw models.Crosswalk
		var gaps, supplements, evidenceMapping []byte

		if err := rows.Scan(
			&cw.ID, &cw.SourceFrameworkID, &cw.SourceControlID,
			&cw.TargetFrameworkID, &cw.TargetControlID,
			&cw.MappingType, &cw.Confidence, &cw.Rationale,
			&gaps, &supplements, &evidenceMapping,
			&cw.CreatedAt, &cw.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning crosswalk: %w", err)
		}

		json.Unmarshal(gaps, &cw.Gaps)
		json.Unmarshal(supplements, &cw.Supplements)
		json.Unmarshal(evidenceMapping, &cw.EvidenceMapping)

		crosswalks = append(crosswalks, cw)
	}

	return crosswalks, rows.Err()
}

// CreateCrosswalk creates a new crosswalk.
func (r *ControlRepository) CreateCrosswalk(ctx context.Context, cw *models.Crosswalk) error {
	gaps, _ := json.Marshal(cw.Gaps)
	supplements, _ := json.Marshal(cw.Supplements)
	evidenceMapping, _ := json.Marshal(cw.EvidenceMapping)

	query := `
		INSERT INTO crosswalks (id, source_framework_id, source_control_id, target_framework_id, target_control_id,
		                        mapping_type, confidence, rationale, gaps, supplements, evidence_mapping)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.Pool.Exec(ctx, query,
		cw.ID, cw.SourceFrameworkID, cw.SourceControlID,
		cw.TargetFrameworkID, cw.TargetControlID,
		cw.MappingType, cw.Confidence, cw.Rationale,
		gaps, supplements, evidenceMapping,
	)
	if err != nil {
		return fmt.Errorf("creating crosswalk: %w", err)
	}

	return nil
}

// DeleteCrosswalk deletes a crosswalk by ID.
func (r *ControlRepository) DeleteCrosswalk(ctx context.Context, id string) error {
	query := `DELETE FROM crosswalks WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting crosswalk: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("crosswalk %s not found", id)
	}

	return nil
}
