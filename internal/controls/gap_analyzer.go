package controls

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/agentguard/agentguard/internal/models"
)

// GapAnalyzer provides gap analysis functionality.
type GapAnalyzer struct {
	service *Service
}

// NewGapAnalyzer creates a new gap analyzer.
func NewGapAnalyzer(dataDir string) (*GapAnalyzer, error) {
	svc, err := NewService(dataDir)
	if err != nil {
		return nil, err
	}
	return &GapAnalyzer{service: svc}, nil
}

// AnalysisInput represents input for gap analysis.
type AnalysisInput struct {
	TargetFramework     string   `json:"target_framework"`
	ImplementedControls []string `json:"implemented_controls"`
	SourceFramework     string   `json:"source_framework,omitempty"`
}

// AnalysisOutput represents the output of gap analysis.
type AnalysisOutput struct {
	Framework          string             `json:"framework"`
	FrameworkName      string             `json:"framework_name"`
	TotalControls      int                `json:"total_controls"`
	ImplementedCount   int                `json:"implemented_count"`
	GapCount           int                `json:"gap_count"`
	CoveragePercentage float64            `json:"coverage_percentage"`
	Gaps               []GapDetail        `json:"gaps"`
	Summary            GapSummaryOutput   `json:"summary"`
	Crosswalks         []CrosswalkSummary `json:"crosswalks,omitempty"`
}

// GapDetail provides details about a specific gap.
type GapDetail struct {
	ControlID          string   `json:"control_id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Priority           string   `json:"priority"`
	EstimatedEffort    string   `json:"estimated_effort"`
	RemediationOptions []string `json:"remediation_options"`
}

// GapSummaryOutput provides aggregate statistics.
type GapSummaryOutput struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

// CrosswalkSummary provides a summary of crosswalk mappings.
type CrosswalkSummary struct {
	SourceControl  string `json:"source_control"`
	TargetControls string `json:"target_controls"`
	MappingType    string `json:"mapping_type"`
	Confidence     string `json:"confidence"`
}

// RunAnalysis performs gap analysis based on the input.
func (g *GapAnalyzer) RunAnalysis(ctx context.Context, input *AnalysisInput) (*AnalysisOutput, error) {
	targetFW := FrameworkID(input.TargetFramework)

	fw, err := g.service.GetFramework(targetFW)
	if err != nil {
		return nil, fmt.Errorf("unknown framework: %s", input.TargetFramework)
	}

	analysis, err := g.service.AnalyzeGaps(ctx, targetFW, input.ImplementedControls)
	if err != nil {
		return nil, err
	}

	controls, err := g.service.GetControls(targetFW)
	if err != nil {
		return nil, fmt.Errorf("getting controls for %s: %w", input.TargetFramework, err)
	}
	controlMap := make(map[string]models.Control)
	for _, c := range controls {
		controlMap[strings.ToLower(c.ControlID)] = c
	}

	gaps := make([]GapDetail, 0, len(analysis.Gaps))
	summary := GapSummaryOutput{}

	for _, gap := range analysis.Gaps {
		ctrl := controlMap[strings.ToLower(gap.ControlID)]
		detail := GapDetail{
			ControlID:          gap.ControlID,
			Title:              ctrl.Title,
			Description:        ctrl.Description,
			Priority:           gap.Priority,
			EstimatedEffort:    gap.EstimatedEffort,
			RemediationOptions: gap.RemediationOptions,
		}
		gaps = append(gaps, detail)

		switch gap.Priority {
		case "critical":
			summary.Critical++
		case "high":
			summary.High++
		case "medium":
			summary.Medium++
		case "low":
			summary.Low++
		}
	}

	output := &AnalysisOutput{
		Framework:          input.TargetFramework,
		FrameworkName:      fw.Name,
		TotalControls:      analysis.Summary.TotalControls,
		ImplementedCount:   analysis.Summary.FullyCovered,
		GapCount:           len(gaps),
		CoveragePercentage: analysis.Summary.CoveragePercentage,
		Gaps:               gaps,
		Summary:            summary,
	}

	// Add crosswalk information if source framework specified
	if input.SourceFramework != "" {
		crosswalks, err := g.service.GetCrosswalks(
			FrameworkID(input.SourceFramework),
			targetFW,
		)
		if err == nil {
			xwSummaries := make([]CrosswalkSummary, 0, len(crosswalks))
			for _, xw := range crosswalks {
				xwSummaries = append(xwSummaries, CrosswalkSummary{
					SourceControl:  xw.SourceControlID,
					TargetControls: xw.TargetControlID,
					MappingType:    string(xw.MappingType),
					Confidence:     fmt.Sprintf("%.0f%%", xw.Confidence*100),
				})
			}
			output.Crosswalks = xwSummaries
		}
	}

	return output, nil
}

// PrintReport prints a formatted gap analysis report.
func (g *GapAnalyzer) PrintReport(w io.Writer, output *AnalysisOutput) {
	fmt.Fprintf(w, "\n╔══════════════════════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(w, "║                          GAP ANALYSIS REPORT                                 ║\n")
	fmt.Fprintf(w, "╚══════════════════════════════════════════════════════════════════════════════╝\n\n")

	fmt.Fprintf(w, "Framework: %s (%s)\n", output.FrameworkName, output.Framework)
	fmt.Fprintf(w, "═══════════════════════════════════════════════════════════════════════════════\n\n")

	fmt.Fprintf(w, "COVERAGE SUMMARY\n")
	fmt.Fprintf(w, "────────────────\n")
	fmt.Fprintf(w, "  Total Controls:      %d\n", output.TotalControls)
	fmt.Fprintf(w, "  Implemented:         %d\n", output.ImplementedCount)
	fmt.Fprintf(w, "  Gaps Identified:     %d\n", output.GapCount)
	fmt.Fprintf(w, "  Coverage:            %.1f%%\n\n", output.CoveragePercentage)

	fmt.Fprintf(w, "GAPS BY PRIORITY\n")
	fmt.Fprintf(w, "────────────────\n")
	fmt.Fprintf(w, "  Critical: %d\n", output.Summary.Critical)
	fmt.Fprintf(w, "  High:     %d\n", output.Summary.High)
	fmt.Fprintf(w, "  Medium:   %d\n", output.Summary.Medium)
	fmt.Fprintf(w, "  Low:      %d\n\n", output.Summary.Low)

	if len(output.Gaps) > 0 {
		fmt.Fprintf(w, "DETAILED GAPS\n")
		fmt.Fprintf(w, "═════════════\n\n")

		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "CONTROL ID\tTITLE\tPRIORITY\tEFFORT\n")
		fmt.Fprintf(tw, "──────────\t─────\t────────\t──────\n")

		for _, gap := range output.Gaps {
			title := gap.Title
			if len(title) > 40 {
				title = title[:37] + "..."
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
				gap.ControlID, title, gap.Priority, gap.EstimatedEffort)
		}
		tw.Flush()
	}

	if len(output.Crosswalks) > 0 {
		fmt.Fprintf(w, "\n\nCROSSWALK MAPPINGS\n")
		fmt.Fprintf(w, "══════════════════\n\n")

		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "SOURCE\tTARGET\tTYPE\tCONFIDENCE\n")
		fmt.Fprintf(tw, "──────\t──────\t────\t──────────\n")

		for _, xw := range output.Crosswalks {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
				xw.SourceControl, xw.TargetControls, xw.MappingType, xw.Confidence)
		}
		tw.Flush()
	}

	fmt.Fprintf(w, "\n")
}

// PrintJSON prints the analysis as JSON.
func (g *GapAnalyzer) PrintJSON(w io.Writer, output *AnalysisOutput) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

// ListFrameworks prints available frameworks.
func (g *GapAnalyzer) ListFrameworks(w io.Writer) {
	fmt.Fprintf(w, "\nAvailable Control Frameworks:\n")
	fmt.Fprintf(w, "══════════════════════════════\n\n")

	frameworks := g.service.ListFrameworks()
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tNAME\tVERSION\n")
	fmt.Fprintf(tw, "──\t────\t───────\n")

	for _, fw := range frameworks {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", fw.ID, fw.Name, fw.Version)
	}
	tw.Flush()
	fmt.Fprintf(w, "\n")
}

// GenerateCrosswalkReport generates a crosswalk report between two frameworks.
func (g *GapAnalyzer) GenerateCrosswalkReport(w io.Writer, source, target string) error {
	sourceFW := FrameworkID(source)
	targetFW := FrameworkID(target)

	sourceName, err := g.service.GetFramework(sourceFW)
	if err != nil {
		return fmt.Errorf("unknown source framework: %s", source)
	}

	targetName, err := g.service.GetFramework(targetFW)
	if err != nil {
		return fmt.Errorf("unknown target framework: %s", target)
	}

	crosswalks, err := g.service.GetCrosswalks(sourceFW, targetFW)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n╔══════════════════════════════════════════════════════════════════════════════╗\n")
	fmt.Fprintf(w, "║                          CROSSWALK REPORT                                    ║\n")
	fmt.Fprintf(w, "╚══════════════════════════════════════════════════════════════════════════════╝\n\n")

	fmt.Fprintf(w, "Source: %s (%s)\n", sourceName.Name, source)
	fmt.Fprintf(w, "Target: %s (%s)\n", targetName.Name, target)
	fmt.Fprintf(w, "══════════════════════════════════════════════════════════════════════════════\n\n")

	if len(crosswalks) == 0 {
		fmt.Fprintf(w, "No predefined crosswalks found between these frameworks.\n\n")
		return nil
	}

	fmt.Fprintf(w, "Found %d control mappings:\n\n", len(crosswalks))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "SOURCE CONTROL\tTARGET CONTROL\tMAPPING TYPE\tCONFIDENCE\n")
	fmt.Fprintf(tw, "──────────────\t──────────────\t────────────\t──────────\n")

	for _, xw := range crosswalks {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%.0f%%\n",
			xw.SourceControlID, xw.TargetControlID, xw.MappingType, xw.Confidence*100)
	}
	tw.Flush()

	fmt.Fprintf(w, "\nMapping Type Legend:\n")
	fmt.Fprintf(w, "  exact    - Controls are equivalent\n")
	fmt.Fprintf(w, "  partial  - Some overlap exists\n")
	fmt.Fprintf(w, "  superset - Source includes target\n")
	fmt.Fprintf(w, "  subset   - Target includes source\n")
	fmt.Fprintf(w, "  related  - Controls address similar topics\n\n")

	return nil
}

// LoadInputFromFile loads analysis input from a JSON file.
func LoadInputFromFile(path string) (*AnalysisInput, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var input AnalysisInput
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	return &input, nil
}
