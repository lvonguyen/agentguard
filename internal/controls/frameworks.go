// Package controls provides control framework management and gap analysis.
package controls

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentguard/agentguard/internal/models"
)

// FrameworkID identifies a control framework.
type FrameworkID string

const (
	FrameworkNISTAIRMF  FrameworkID = "nist-ai-rmf"
	FrameworkNIST80053  FrameworkID = "nist-800-53"
	FrameworkISO42001   FrameworkID = "iso-42001"
	FrameworkSOC2       FrameworkID = "soc2"
)

// Service provides control framework operations.
type Service struct {
	dataDir    string
	frameworks map[FrameworkID]*models.Framework
	controls   map[FrameworkID][]models.Control
	crosswalks []models.Crosswalk
}

// NewService creates a new control framework service.
func NewService(dataDir string) (*Service, error) {
	s := &Service{
		dataDir:    dataDir,
		frameworks: make(map[FrameworkID]*models.Framework),
		controls:   make(map[FrameworkID][]models.Control),
	}

	if err := s.loadFrameworks(); err != nil {
		return nil, fmt.Errorf("loading frameworks: %w", err)
	}

	return s, nil
}

// loadFrameworks loads all framework definitions from the data directory.
func (s *Service) loadFrameworks() error {
	// Load embedded framework definitions
	s.loadEmbeddedFrameworks()

	// Override with files from data directory if present
	if s.dataDir != "" {
		files, err := filepath.Glob(filepath.Join(s.dataDir, "frameworks", "*.json"))
		if err != nil {
			return err
		}
		for _, file := range files {
			if err := s.loadFrameworkFile(file); err != nil {
				return fmt.Errorf("loading %s: %w", file, err)
			}
		}
	}

	return nil
}

func (s *Service) loadFrameworkFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fw models.Framework
	if err := json.Unmarshal(data, &fw); err != nil {
		return err
	}

	s.frameworks[FrameworkID(fw.ID)] = &fw
	return nil
}

// loadEmbeddedFrameworks loads built-in framework definitions.
func (s *Service) loadEmbeddedFrameworks() {
	// NIST AI RMF
	s.frameworks[FrameworkNISTAIRMF] = &models.Framework{
		ID:          string(FrameworkNISTAIRMF),
		Name:        "NIST AI Risk Management Framework",
		Version:     "1.0",
		Publisher:   "NIST",
		Description: "Framework for managing risks associated with AI systems",
		URL:         "https://www.nist.gov/itl/ai-risk-management-framework",
	}
	s.controls[FrameworkNISTAIRMF] = getNISTAIRMFControls()

	// NIST 800-53
	s.frameworks[FrameworkNIST80053] = &models.Framework{
		ID:          string(FrameworkNIST80053),
		Name:        "NIST SP 800-53 Rev 5",
		Version:     "5.1",
		Publisher:   "NIST",
		Description: "Security and Privacy Controls for Information Systems",
		URL:         "https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final",
	}
	s.controls[FrameworkNIST80053] = getNIST80053Controls()

	// ISO 42001
	s.frameworks[FrameworkISO42001] = &models.Framework{
		ID:          string(FrameworkISO42001),
		Name:        "ISO/IEC 42001:2023",
		Version:     "2023",
		Publisher:   "ISO/IEC",
		Description: "AI Management System - Requirements with guidance for use",
		URL:         "https://www.iso.org/standard/81230.html",
	}
	s.controls[FrameworkISO42001] = getISO42001Controls()
}

// GetFramework returns a framework by ID.
func (s *Service) GetFramework(id FrameworkID) (*models.Framework, error) {
	fw, ok := s.frameworks[id]
	if !ok {
		return nil, fmt.Errorf("framework not found: %s", id)
	}
	return fw, nil
}

// ListFrameworks returns all available frameworks.
func (s *Service) ListFrameworks() []*models.Framework {
	result := make([]*models.Framework, 0, len(s.frameworks))
	for _, fw := range s.frameworks {
		result = append(result, fw)
	}
	return result
}

// GetControls returns controls for a framework.
func (s *Service) GetControls(id FrameworkID) ([]models.Control, error) {
	controls, ok := s.controls[id]
	if !ok {
		return nil, fmt.Errorf("controls not found for framework: %s", id)
	}
	return controls, nil
}

// GetCrosswalks returns mappings between two frameworks.
func (s *Service) GetCrosswalks(source, target FrameworkID) ([]models.Crosswalk, error) {
	result := []models.Crosswalk{}

	for _, xw := range s.crosswalks {
		if xw.SourceFrameworkID == string(source) && xw.TargetFrameworkID == string(target) {
			result = append(result, xw)
		}
	}

	// Generate crosswalks dynamically if not pre-loaded
	if len(result) == 0 {
		generated, err := s.generateCrosswalks(source, target)
		if err != nil {
			return nil, err
		}
		result = generated
	}

	return result, nil
}

// generateCrosswalks creates mappings between frameworks using known relationships.
func (s *Service) generateCrosswalks(source, target FrameworkID) ([]models.Crosswalk, error) {
	sourceControls, err := s.GetControls(source)
	if err != nil {
		return nil, err
	}

	targetControls, err := s.GetControls(target)
	if err != nil {
		return nil, err
	}

	crosswalks := []models.Crosswalk{}

	// Use predefined mapping tables based on framework pair
	mappings := getCrosswalkMappings(source, target)

	for _, sc := range sourceControls {
		for _, tc := range targetControls {
			if mapping, ok := mappings[sc.ControlID]; ok {
				for _, targetID := range mapping.TargetIDs {
					if tc.ControlID == targetID {
						crosswalks = append(crosswalks, models.Crosswalk{
							SourceFrameworkID: string(source),
							SourceControlID:   sc.ControlID,
							TargetFrameworkID: string(target),
							TargetControlID:   tc.ControlID,
							MappingType:       mapping.Type,
							Confidence:        mapping.Confidence,
							Rationale:         mapping.Rationale,
						})
					}
				}
			}
		}
	}

	return crosswalks, nil
}

// AnalyzeGaps performs gap analysis between current state and target framework.
func (s *Service) AnalyzeGaps(ctx context.Context, targetFramework FrameworkID, implementedControls []string) (*models.GapAnalysis, error) {
	controls, err := s.GetControls(targetFramework)
	if err != nil {
		return nil, err
	}

	implemented := make(map[string]bool)
	for _, c := range implementedControls {
		implemented[strings.ToLower(c)] = true
	}

	gaps := []models.ControlGap{}
	fullyCovered := 0
	partiallyCovered := 0

	for _, ctrl := range controls {
		ctrlID := strings.ToLower(ctrl.ControlID)
		if implemented[ctrlID] {
			fullyCovered++
			continue
		}

		// Check for partial coverage (parent or child implemented)
		partial := false
		if ctrl.ParentControlID != nil && implemented[strings.ToLower(*ctrl.ParentControlID)] {
			partial = true
			partiallyCovered++
		}

		if !partial {
			gap := models.ControlGap{
				ControlID:   ctrl.ControlID,
				GapType:     "not_implemented",
				Description: fmt.Sprintf("Control '%s' (%s) is not implemented", ctrl.ControlID, ctrl.Title),
				Priority:    determineGapPriority(ctrl),
				RemediationOptions: generateRemediationOptions(ctrl),
				EstimatedEffort:    estimateEffort(ctrl),
			}
			gaps = append(gaps, gap)
		}
	}

	totalControls := len(controls)
	notCovered := totalControls - fullyCovered - partiallyCovered
	coverage := float64(fullyCovered) / float64(totalControls) * 100

	gapsByPriority := make(map[string]int)
	for _, g := range gaps {
		gapsByPriority[g.Priority]++
	}

	return &models.GapAnalysis{
		TargetFrameworkID: string(targetFramework),
		Gaps:              gaps,
		Summary: models.GapSummary{
			TotalControls:      totalControls,
			FullyCovered:       fullyCovered,
			PartiallyCovered:   partiallyCovered,
			NotCovered:         notCovered,
			CoveragePercentage: coverage,
			GapsByPriority:     gapsByPriority,
		},
	}, nil
}

func determineGapPriority(ctrl models.Control) string {
	// Priority based on control characteristics
	for _, layer := range ctrl.ApplicableLayers {
		if layer == "governance" || layer == "risk_management" {
			return "high"
		}
	}
	if len(ctrl.EvidenceTypes) > 3 {
		return "high"
	}
	if len(ctrl.Activities) > 2 {
		return "medium"
	}
	return "low"
}

func generateRemediationOptions(ctrl models.Control) []string {
	options := []string{}
	for _, activity := range ctrl.Activities {
		options = append(options, fmt.Sprintf("Implement: %s", activity))
	}
	if len(options) == 0 {
		options = append(options, "Review control requirements and implement appropriate measures")
	}
	return options
}

func estimateEffort(ctrl models.Control) string {
	activityCount := len(ctrl.Activities)
	evidenceCount := len(ctrl.EvidenceTypes)

	total := activityCount + evidenceCount
	if total > 6 {
		return "large"
	}
	if total > 3 {
		return "medium"
	}
	return "small"
}
