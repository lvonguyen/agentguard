package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agentguard/agentguard/internal/models"
	"github.com/agentguard/agentguard/internal/repository"
)

// mockControlRepo implements repository.ControlRepository for unit testing
// without a live database connection.
type mockControlRepo struct {
	frameworks []models.Framework
	controls   []models.Control
	crosswalks []models.Crosswalk

	listControlsErr   error
	getControlErr     error
	createControlErr  error
	updateControlErr  error
	deleteControlErr  error
	getCrosswalkErr   error
}

func (m *mockControlRepo) ListFrameworks(_ context.Context) ([]models.Framework, error) {
	return m.frameworks, nil
}

func (m *mockControlRepo) GetFramework(_ context.Context, id string) (*models.Framework, error) {
	for i := range m.frameworks {
		if m.frameworks[i].ID == id {
			return &m.frameworks[i], nil
		}
	}
	return nil, nil
}

func (m *mockControlRepo) CreateFramework(_ context.Context, f *models.Framework) error {
	m.frameworks = append(m.frameworks, *f)
	return nil
}

func (m *mockControlRepo) UpdateFramework(_ context.Context, f *models.Framework) error {
	for i := range m.frameworks {
		if m.frameworks[i].ID == f.ID {
			m.frameworks[i] = *f
			return nil
		}
	}
	return nil
}

func (m *mockControlRepo) DeleteFramework(_ context.Context, id string) error {
	for i := range m.frameworks {
		if m.frameworks[i].ID == id {
			m.frameworks = append(m.frameworks[:i], m.frameworks[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockControlRepo) ListControls(_ context.Context, frameworkID string) ([]models.Control, error) {
	if m.listControlsErr != nil {
		return nil, m.listControlsErr
	}
	var result []models.Control
	for _, c := range m.controls {
		if c.FrameworkID == frameworkID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *mockControlRepo) GetControl(_ context.Context, id string) (*models.Control, error) {
	if m.getControlErr != nil {
		return nil, m.getControlErr
	}
	for i := range m.controls {
		if m.controls[i].ID == id {
			return &m.controls[i], nil
		}
	}
	return nil, nil
}

func (m *mockControlRepo) CreateControl(_ context.Context, c *models.Control) error {
	if m.createControlErr != nil {
		return m.createControlErr
	}
	m.controls = append(m.controls, *c)
	return nil
}

func (m *mockControlRepo) UpdateControl(_ context.Context, c *models.Control) error {
	if m.updateControlErr != nil {
		return m.updateControlErr
	}
	for i := range m.controls {
		if m.controls[i].ID == c.ID {
			m.controls[i] = *c
			return nil
		}
	}
	return nil
}

func (m *mockControlRepo) DeleteControl(_ context.Context, id string) error {
	if m.deleteControlErr != nil {
		return m.deleteControlErr
	}
	for i := range m.controls {
		if m.controls[i].ID == id {
			m.controls = append(m.controls[:i], m.controls[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *mockControlRepo) GetCrosswalk(_ context.Context, sourceFrameworkID, targetFrameworkID string) ([]models.Crosswalk, error) {
	if m.getCrosswalkErr != nil {
		return nil, m.getCrosswalkErr
	}
	var result []models.Crosswalk
	for _, cw := range m.crosswalks {
		if cw.SourceFrameworkID == sourceFrameworkID && cw.TargetFrameworkID == targetFrameworkID {
			result = append(result, cw)
		}
	}
	return result, nil
}

func (m *mockControlRepo) CreateCrosswalk(_ context.Context, cw *models.Crosswalk) error {
	m.crosswalks = append(m.crosswalks, *cw)
	return nil
}

func (m *mockControlRepo) DeleteCrosswalk(_ context.Context, id string) error {
	for i := range m.crosswalks {
		if m.crosswalks[i].ID == id {
			m.crosswalks = append(m.crosswalks[:i], m.crosswalks[i+1:]...)
			return nil
		}
	}
	return nil
}

// Compile-time interface check
var _ repository.ControlRepository = (*mockControlRepo)(nil)

func seedControls() []models.Control {
	return []models.Control{
		{
			ID:               "ctrl-001",
			FrameworkID:      "nist-ai-rmf",
			ControlID:        "GOVERN-1",
			Title:            "AI Risk Management Policy",
			Description:      "Establish AI risk management policies.",
			Objectives:       []string{"Define risk tolerance"},
			Activities:       []string{"Draft policy document"},
			EvidenceTypes:    []string{"policy_document"},
			ApplicableLayers: []string{"governance"},
		},
		{
			ID:               "ctrl-002",
			FrameworkID:      "nist-ai-rmf",
			ControlID:        "GOVERN-2",
			Title:            "AI Governance Structure",
			Description:      "Establish governance structure for AI.",
			Objectives:       []string{"Define roles and responsibilities"},
			Activities:       []string{"Appoint AI governance board"},
			EvidenceTypes:    []string{"org_chart"},
			ApplicableLayers: []string{"governance"},
		},
		{
			ID:               "ctrl-003",
			FrameworkID:      "nist-800-53",
			ControlID:        "AC-1",
			Title:            "Access Control Policy",
			Description:      "Develop access control policy.",
			Objectives:       []string{"Define access control requirements"},
			Activities:       []string{"Document access control procedures"},
			EvidenceTypes:    []string{"policy_document"},
			ApplicableLayers: []string{"platform"},
		},
	}
}

func TestListControls(t *testing.T) {
	tests := []struct {
		name        string
		frameworkID string
		seed        []models.Control
		wantCount   int
		wantErr     bool
		setupErr    error
	}{
		{
			name:        "returns controls for matching framework",
			frameworkID: "nist-ai-rmf",
			seed:        seedControls(),
			wantCount:   2,
		},
		{
			name:        "returns empty for unknown framework",
			frameworkID: "nonexistent",
			seed:        seedControls(),
			wantCount:   0,
		},
		{
			name:        "returns empty when no controls exist",
			frameworkID: "nist-ai-rmf",
			seed:        nil,
			wantCount:   0,
		},
		{
			name:        "propagates repository error",
			frameworkID: "nist-ai-rmf",
			seed:        seedControls(),
			wantErr:     true,
			setupErr:    fmt.Errorf("connection refused"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockControlRepo{
				controls:        tt.seed,
				listControlsErr: tt.setupErr,
			}

			got, err := repo.ListControls(context.Background(), tt.frameworkID)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("got %d controls, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestGetControl(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		seed     []models.Control
		wantNil  bool
		wantID   string
		wantErr  bool
		setupErr error
	}{
		{
			name:   "returns existing control",
			id:     "ctrl-001",
			seed:   seedControls(),
			wantID: "ctrl-001",
		},
		{
			name:    "returns nil for unknown ID",
			id:      "ctrl-999",
			seed:    seedControls(),
			wantNil: true,
		},
		{
			name:     "propagates repository error",
			id:       "ctrl-001",
			seed:     seedControls(),
			wantErr:  true,
			setupErr: fmt.Errorf("timeout"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockControlRepo{
				controls:      tt.seed,
				getControlErr: tt.setupErr,
			}

			got, err := repo.GetControl(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected control, got nil")
			}
			if got.ID != tt.wantID {
				t.Errorf("got ID %q, want %q", got.ID, tt.wantID)
			}
		})
	}
}

func TestCreateControl(t *testing.T) {
	tests := []struct {
		name     string
		control  models.Control
		wantErr  bool
		setupErr error
	}{
		{
			name: "creates control successfully",
			control: models.Control{
				ID:          "ctrl-new",
				FrameworkID: "nist-ai-rmf",
				ControlID:   "MAP-1",
				Title:       "Context Mapping",
			},
		},
		{
			name: "propagates repository error",
			control: models.Control{
				ID:          "ctrl-fail",
				FrameworkID: "nist-ai-rmf",
				ControlID:   "MAP-2",
				Title:       "Stakeholder Engagement",
			},
			wantErr:  true,
			setupErr: fmt.Errorf("unique constraint violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockControlRepo{
				createControlErr: tt.setupErr,
			}

			err := repo.CreateControl(context.Background(), &tt.control)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, _ := repo.GetControl(context.Background(), tt.control.ID)
			if got == nil {
				t.Fatal("control not found after create")
			}
			if got.Title != tt.control.Title {
				t.Errorf("got title %q, want %q", got.Title, tt.control.Title)
			}
		})
	}
}

func TestDeleteControl(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		seed     []models.Control
		wantErr  bool
		setupErr error
	}{
		{
			name: "deletes existing control",
			id:   "ctrl-001",
			seed: seedControls(),
		},
		{
			name: "no-op for unknown ID",
			id:   "ctrl-999",
			seed: seedControls(),
		},
		{
			name:     "propagates repository error",
			id:       "ctrl-001",
			seed:     seedControls(),
			wantErr:  true,
			setupErr: fmt.Errorf("foreign key violation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockControlRepo{
				controls:         tt.seed,
				deleteControlErr: tt.setupErr,
			}

			err := repo.DeleteControl(context.Background(), tt.id)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, _ := repo.GetControl(context.Background(), tt.id)
			if got != nil {
				t.Errorf("control %q still exists after delete", tt.id)
			}
		})
	}
}

func TestGetCrosswalk(t *testing.T) {
	crosswalks := []models.Crosswalk{
		{
			ID:                "xw-001",
			SourceFrameworkID: "nist-ai-rmf",
			SourceControlID:   "GOVERN-1",
			TargetFrameworkID: "nist-800-53",
			TargetControlID:   "PM-9",
			MappingType:       models.MappingPartial,
			Confidence:        0.85,
			Rationale:         "AI RMF governance maps to 800-53 risk management.",
		},
		{
			ID:                "xw-002",
			SourceFrameworkID: "nist-ai-rmf",
			SourceControlID:   "GOVERN-2",
			TargetFrameworkID: "nist-800-53",
			TargetControlID:   "PM-1",
			MappingType:       models.MappingExact,
			Confidence:        0.95,
			Rationale:         "Direct governance structure mapping.",
		},
	}

	tests := []struct {
		name      string
		source    string
		target    string
		wantCount int
		wantErr   bool
		setupErr  error
	}{
		{
			name:      "returns crosswalks for valid pair",
			source:    "nist-ai-rmf",
			target:    "nist-800-53",
			wantCount: 2,
		},
		{
			name:      "returns empty for unknown pair",
			source:    "nist-ai-rmf",
			target:    "iso-42001",
			wantCount: 0,
		},
		{
			name:     "propagates repository error",
			source:   "nist-ai-rmf",
			target:   "nist-800-53",
			wantErr:  true,
			setupErr: fmt.Errorf("connection lost"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockControlRepo{
				crosswalks:      crosswalks,
				getCrosswalkErr: tt.setupErr,
			}

			got, err := repo.GetCrosswalk(context.Background(), tt.source, tt.target)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Errorf("got %d crosswalks, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestUpdateControl(t *testing.T) {
	tests := []struct {
		name     string
		update   models.Control
		seed     []models.Control
		wantErr  bool
		setupErr error
	}{
		{
			name: "updates existing control title",
			update: models.Control{
				ID:          "ctrl-001",
				FrameworkID: "nist-ai-rmf",
				ControlID:   "GOVERN-1",
				Title:       "Updated AI Risk Policy",
			},
			seed: seedControls(),
		},
		{
			name: "propagates repository error",
			update: models.Control{
				ID:    "ctrl-001",
				Title: "Should Fail",
			},
			seed:     seedControls(),
			wantErr:  true,
			setupErr: fmt.Errorf("serialization failure"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockControlRepo{
				controls:         tt.seed,
				updateControlErr: tt.setupErr,
			}

			err := repo.UpdateControl(context.Background(), &tt.update)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got, _ := repo.GetControl(context.Background(), tt.update.ID)
			if got == nil {
				t.Fatal("control not found after update")
			}
			if got.Title != tt.update.Title {
				t.Errorf("got title %q, want %q", got.Title, tt.update.Title)
			}
		})
	}
}
