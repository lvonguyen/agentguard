package controls

import "github.com/agentguard/agentguard/internal/models"

// CrosswalkMapping defines a mapping between controls.
type CrosswalkMapping struct {
	TargetIDs  []string
	Type       models.MappingType
	Confidence float64
	Rationale  string
}

// getCrosswalkMappings returns predefined crosswalk mappings between frameworks.
func getCrosswalkMappings(source, target FrameworkID) map[string]CrosswalkMapping {
	key := string(source) + "->" + string(target)

	mappings := map[string]map[string]CrosswalkMapping{
		// NIST AI RMF -> NIST 800-53
		string(FrameworkNISTAIRMF) + "->" + string(FrameworkNIST80053): {
			"GOVERN-1": {
				TargetIDs:  []string{"PL-1", "PL-2"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Legal/regulatory requirements map to security planning controls",
			},
			"GOVERN-2": {
				TargetIDs:  []string{"AC-1", "AU-1", "CM-1"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Accountability structures map to policy controls across families",
			},
			"GOVERN-4": {
				TargetIDs:  []string{"RA-1", "RA-3"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Risk tolerance maps to risk assessment controls",
			},
			"GOVERN-6": {
				TargetIDs:  []string{"CM-3", "CM-4"},
				Type:       models.MappingPartial,
				Confidence: 0.6,
				Rationale:  "Third-party/supply chain maps to configuration management",
			},
			"MAP-1": {
				TargetIDs:  []string{"PL-2"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Context establishment maps to system planning",
			},
			"MAP-2": {
				TargetIDs:  []string{"RA-3"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "AI categorization maps to risk assessment",
			},
			"MAP-4": {
				TargetIDs:  []string{"RA-3", "CM-4"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Risk/benefit mapping relates to risk and impact analysis",
			},
			"MEASURE-1": {
				TargetIDs:  []string{"AU-2", "AU-3", "SI-4"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "AI measurement maps to logging and monitoring controls",
			},
			"MEASURE-2": {
				TargetIDs:  []string{"RA-5"},
				Type:       models.MappingPartial,
				Confidence: 0.6,
				Rationale:  "AI evaluation maps to vulnerability monitoring",
			},
			"MEASURE-3": {
				TargetIDs:  []string{"AU-6", "SI-4"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Tracking mechanisms map to audit review and monitoring",
			},
			"MANAGE-1": {
				TargetIDs:  []string{"RA-3"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Risk prioritization directly maps to risk assessment",
			},
			"MANAGE-3": {
				TargetIDs:  []string{"PL-1", "RA-1"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "ERM integration maps to planning and risk policy",
			},
			"MANAGE-4": {
				TargetIDs:  []string{"CP-1", "CP-2"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Risk treatment and recovery maps to contingency planning",
			},
		},

		// NIST AI RMF -> ISO 42001
		string(FrameworkNISTAIRMF) + "->" + string(FrameworkISO42001): {
			"GOVERN-1": {
				TargetIDs:  []string{"ISO42001-4.1", "ISO42001-4.2"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Legal requirements align with organizational context",
			},
			"GOVERN-2": {
				TargetIDs:  []string{"ISO42001-5.3"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Accountability structures map directly to roles/responsibilities",
			},
			"GOVERN-3": {
				TargetIDs:  []string{"ISO42001-A.3.2"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Workforce diversity relates to bias assessment",
			},
			"GOVERN-4": {
				TargetIDs:  []string{"ISO42001-6.1"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Risk tolerance maps to risk/opportunity planning",
			},
			"GOVERN-5": {
				TargetIDs:  []string{"ISO42001-7.4"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Stakeholder engagement maps to communication",
			},
			"GOVERN-6": {
				TargetIDs:  []string{"ISO42001-8.6"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Third-party policies map directly to third-party considerations",
			},
			"MAP-1": {
				TargetIDs:  []string{"ISO42001-4.3", "ISO42001-8.4"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Context establishment maps to scope and documentation",
			},
			"MAP-2": {
				TargetIDs:  []string{"ISO42001-8.2"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Categorization maps to impact assessment",
			},
			"MAP-3": {
				TargetIDs:  []string{"ISO42001-A.2.2"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Capabilities/limitations maps to transparency",
			},
			"MAP-4": {
				TargetIDs:  []string{"ISO42001-6.1", "ISO42001-8.2"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Risk/benefit mapping aligns with risk planning and impact assessment",
			},
			"MAP-5": {
				TargetIDs:  []string{"ISO42001-8.2", "ISO42001-A.7.3"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Impact to communities maps to impact assessment and privacy",
			},
			"MEASURE-1": {
				TargetIDs:  []string{"ISO42001-9.1"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Measurement methods map directly to monitoring",
			},
			"MEASURE-2": {
				TargetIDs:  []string{"ISO42001-A.6.2"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Trustworthiness evaluation relates to reliability",
			},
			"MEASURE-3": {
				TargetIDs:  []string{"ISO42001-9.1"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Tracking mechanisms map to monitoring",
			},
			"MEASURE-4": {
				TargetIDs:  []string{"ISO42001-10.2"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Feedback incorporation maps to continual improvement",
			},
			"MANAGE-1": {
				TargetIDs:  []string{"ISO42001-6.1"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Risk prioritization maps to risk planning",
			},
			"MANAGE-2": {
				TargetIDs:  []string{"ISO42001-6.2"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Maximizing benefits maps to AI objectives",
			},
			"MANAGE-3": {
				TargetIDs:  []string{"ISO42001-4.4"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "ERM integration maps to AIMS establishment",
			},
			"MANAGE-4": {
				TargetIDs:  []string{"ISO42001-10.1"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Risk treatment maps to corrective action",
			},
		},

		// ISO 42001 -> NIST 800-53
		string(FrameworkISO42001) + "->" + string(FrameworkNIST80053): {
			"ISO42001-4.1": {
				TargetIDs:  []string{"PL-1", "PL-2"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Organization context maps to planning controls",
			},
			"ISO42001-5.1": {
				TargetIDs:  []string{"PL-1"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Leadership commitment maps to policy controls",
			},
			"ISO42001-5.2": {
				TargetIDs:  []string{"AC-1", "AU-1", "CM-1", "PL-1"},
				Type:       models.MappingPartial,
				Confidence: 0.6,
				Rationale:  "AI policy maps to various policy controls",
			},
			"ISO42001-6.1": {
				TargetIDs:  []string{"RA-1", "RA-3"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Risk planning maps directly to risk assessment",
			},
			"ISO42001-7.2": {
				TargetIDs:  []string{"AC-2"},
				Type:       models.MappingPartial,
				Confidence: 0.5,
				Rationale:  "Competence requirements relate to account management",
			},
			"ISO42001-7.5": {
				TargetIDs:  []string{"SI-12"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Documentation requirements map to information management",
			},
			"ISO42001-8.1": {
				TargetIDs:  []string{"CM-2", "CM-3"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Operational planning maps to configuration management",
			},
			"ISO42001-8.2": {
				TargetIDs:  []string{"RA-3", "CM-4"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "Impact assessment maps to risk and impact analysis",
			},
			"ISO42001-9.1": {
				TargetIDs:  []string{"AU-2", "AU-6", "SI-4"},
				Type:       models.MappingPartial,
				Confidence: 0.8,
				Rationale:  "Monitoring maps to audit and system monitoring",
			},
			"ISO42001-9.2": {
				TargetIDs:  []string{"AU-6"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Internal audit maps to audit review",
			},
			"ISO42001-10.1": {
				TargetIDs:  []string{"CP-2"},
				Type:       models.MappingPartial,
				Confidence: 0.6,
				Rationale:  "Corrective action maps to contingency planning",
			},
			"ISO42001-A.4.4": {
				TargetIDs:  []string{"AC-3", "AC-6", "SC-7", "SC-8"},
				Type:       models.MappingExact,
				Confidence: 0.9,
				Rationale:  "AI security maps directly to access and communications controls",
			},
			"ISO42001-A.5.2": {
				TargetIDs:  []string{"AC-2", "AC-6"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Human oversight relates to access controls",
			},
			"ISO42001-A.7.3": {
				TargetIDs:  []string{"SI-12"},
				Type:       models.MappingPartial,
				Confidence: 0.7,
				Rationale:  "Privacy protection maps to information management",
			},
		},
	}

	if m, ok := mappings[key]; ok {
		return m
	}
	return map[string]CrosswalkMapping{}
}
