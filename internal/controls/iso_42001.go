package controls

import "github.com/agentguard/agentguard/internal/models"

// getISO42001Controls returns ISO/IEC 42001:2023 AI Management System controls.
// This standard specifies requirements for establishing, implementing, maintaining,
// and continually improving an AI management system within organizations.
func getISO42001Controls() []models.Control {
	return []models.Control{
		// Clause 4: Context of the Organization
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-4.1",
			Title:       "Understanding the Organization and Its Context",
			Description: "Determine external and internal issues relevant to AI management system purpose and strategic direction.",
			Objectives: []string{
				"Identify external issues affecting AI systems",
				"Identify internal issues affecting AI systems",
				"Understand organizational context for AI",
			},
			Activities: []string{
				"Conduct environmental analysis",
				"Assess organizational AI maturity",
				"Document context factors",
			},
			EvidenceTypes: []string{
				"Context analysis documentation",
				"PESTLE analysis",
				"Internal capability assessment",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-4.2",
			Title:       "Understanding Stakeholder Needs",
			Description: "Determine interested parties relevant to the AI management system and their requirements.",
			Objectives: []string{
				"Identify AI stakeholders",
				"Understand stakeholder requirements",
				"Document stakeholder expectations",
			},
			Activities: []string{
				"Map AI stakeholders",
				"Gather stakeholder requirements",
				"Prioritize stakeholder needs",
			},
			EvidenceTypes: []string{
				"Stakeholder register",
				"Requirements documentation",
				"Stakeholder communication records",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-4.3",
			Title:       "Scope of the AI Management System",
			Description: "Determine boundaries and applicability of the AI management system.",
			Objectives: []string{
				"Define AIMS boundaries",
				"Identify applicable AI systems",
			},
			Activities: []string{
				"Define system boundaries",
				"Document included AI systems",
				"Maintain scope documentation",
			},
			EvidenceTypes: []string{
				"Scope statement",
				"AI system inventory",
				"Boundary documentation",
			},
			ApplicableLayers: []string{"governance", "system"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-4.4",
			Title:       "AI Management System",
			Description: "Establish, implement, maintain, and continually improve an AI management system.",
			Objectives: []string{
				"Establish AIMS processes",
				"Implement management system",
				"Enable continuous improvement",
			},
			Activities: []string{
				"Design AIMS architecture",
				"Implement required processes",
				"Establish improvement mechanisms",
			},
			EvidenceTypes: []string{
				"AIMS documentation",
				"Process documentation",
				"Improvement records",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},

		// Clause 5: Leadership
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-5.1",
			Title:       "Leadership and Commitment",
			Description: "Top management demonstrates leadership and commitment to the AI management system.",
			Objectives: []string{
				"Demonstrate management commitment",
				"Allocate resources",
				"Support AI governance",
			},
			Activities: []string{
				"Establish AI governance structure",
				"Communicate importance of AIMS",
				"Allocate adequate resources",
				"Review AIMS effectiveness",
			},
			EvidenceTypes: []string{
				"Management commitment statements",
				"Resource allocation records",
				"Management review minutes",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-5.2",
			Title:       "AI Policy",
			Description: "Establish an AI policy appropriate to the organization's purpose.",
			Objectives: []string{
				"Define AI policy",
				"Communicate policy",
				"Ensure policy availability",
			},
			Activities: []string{
				"Develop AI policy",
				"Obtain management approval",
				"Communicate to stakeholders",
				"Review and update periodically",
			},
			EvidenceTypes: []string{
				"AI policy document",
				"Policy approval records",
				"Communication records",
			},
			ApplicableLayers: []string{"governance"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-5.3",
			Title:       "Organizational Roles and Responsibilities",
			Description: "Assign responsibilities and authorities for relevant roles.",
			Objectives: []string{
				"Define AI-related roles",
				"Assign responsibilities",
				"Communicate authorities",
			},
			Activities: []string{
				"Define AIMS roles",
				"Document responsibilities",
				"Communicate role assignments",
			},
			EvidenceTypes: []string{
				"Role definitions",
				"RACI matrix",
				"Organization charts",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},

		// Clause 6: Planning
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-6.1",
			Title:       "Actions to Address Risks and Opportunities",
			Description: "Plan actions to address AI-related risks and opportunities.",
			Objectives: []string{
				"Identify AI risks",
				"Identify AI opportunities",
				"Plan risk treatments",
			},
			Activities: []string{
				"Conduct AI risk assessment",
				"Identify improvement opportunities",
				"Develop risk treatment plans",
				"Integrate actions into AIMS",
			},
			EvidenceTypes: []string{
				"Risk assessment reports",
				"Risk treatment plans",
				"Opportunity assessment",
			},
			ApplicableLayers: []string{"risk_management", "governance"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-6.2",
			Title:       "AI Objectives and Planning",
			Description: "Establish AI objectives at relevant functions and levels.",
			Objectives: []string{
				"Define AI objectives",
				"Align with AI policy",
				"Plan objective achievement",
			},
			Activities: []string{
				"Establish measurable objectives",
				"Define action plans",
				"Allocate resources",
				"Monitor progress",
			},
			EvidenceTypes: []string{
				"AI objectives documentation",
				"Action plans",
				"Progress reports",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-6.3",
			Title:       "Planning of Changes",
			Description: "Plan changes to the AI management system in a systematic manner.",
			Objectives: []string{
				"Manage AIMS changes",
				"Assess change impacts",
				"Maintain system integrity",
			},
			Activities: []string{
				"Establish change process",
				"Assess change requests",
				"Implement changes systematically",
			},
			EvidenceTypes: []string{
				"Change management process",
				"Change request records",
				"Impact assessments",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},

		// Clause 7: Support
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-7.1",
			Title:       "Resources",
			Description: "Determine and provide resources needed for the AI management system.",
			Objectives: []string{
				"Identify resource needs",
				"Provide adequate resources",
				"Maintain resource availability",
			},
			Activities: []string{
				"Assess resource requirements",
				"Allocate human resources",
				"Provide technical infrastructure",
				"Budget for AI activities",
			},
			EvidenceTypes: []string{
				"Resource plans",
				"Budget allocations",
				"Infrastructure documentation",
			},
			ApplicableLayers: []string{"organization", "operations"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-7.2",
			Title:       "Competence",
			Description: "Ensure persons doing work affecting AI system performance are competent.",
			Objectives: []string{
				"Define competency requirements",
				"Ensure staff competence",
				"Address competency gaps",
			},
			Activities: []string{
				"Define competency requirements",
				"Assess current competencies",
				"Provide training",
				"Evaluate training effectiveness",
			},
			EvidenceTypes: []string{
				"Competency framework",
				"Training records",
				"Competency assessments",
			},
			ApplicableLayers: []string{"organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-7.3",
			Title:       "Awareness",
			Description: "Ensure persons are aware of AI policy and their contribution to AIMS.",
			Objectives: []string{
				"Promote AI policy awareness",
				"Communicate individual roles",
				"Ensure understanding of consequences",
			},
			Activities: []string{
				"Conduct awareness campaigns",
				"Communicate role expectations",
				"Reinforce policy understanding",
			},
			EvidenceTypes: []string{
				"Awareness program documentation",
				"Communication records",
				"Awareness assessment results",
			},
			ApplicableLayers: []string{"organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-7.4",
			Title:       "Communication",
			Description: "Determine internal and external communications relevant to AIMS.",
			Objectives: []string{
				"Plan AI communications",
				"Enable effective information flow",
				"Ensure transparency",
			},
			Activities: []string{
				"Define communication requirements",
				"Establish communication channels",
				"Document communications",
			},
			EvidenceTypes: []string{
				"Communication plan",
				"Communication records",
				"Stakeholder correspondence",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-7.5",
			Title:       "Documented Information",
			Description: "Manage documented information required by the AI management system.",
			Objectives: []string{
				"Maintain required documentation",
				"Control document creation and updates",
				"Ensure document availability",
			},
			Activities: []string{
				"Create required documents",
				"Control document versions",
				"Distribute documents appropriately",
				"Retain records",
			},
			EvidenceTypes: []string{
				"Document control procedure",
				"Document register",
				"Retention schedules",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},

		// Clause 8: Operation
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-8.1",
			Title:       "Operational Planning and Control",
			Description: "Plan, implement, and control processes for AI operations.",
			Objectives: []string{
				"Plan AI operations",
				"Control AI processes",
				"Manage operational changes",
			},
			Activities: []string{
				"Define operational criteria",
				"Implement process controls",
				"Monitor operations",
				"Control outsourced processes",
			},
			EvidenceTypes: []string{
				"Operational procedures",
				"Process documentation",
				"Control records",
			},
			ApplicableLayers: []string{"operations", "system"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-8.2",
			Title:       "AI System Impact Assessment",
			Description: "Conduct impact assessment before AI system deployment.",
			Objectives: []string{
				"Assess AI system impacts",
				"Identify potential harms",
				"Document assessment results",
			},
			Activities: []string{
				"Conduct impact assessments",
				"Evaluate risks and benefits",
				"Document findings",
				"Implement mitigations",
			},
			EvidenceTypes: []string{
				"Impact assessment reports",
				"Risk-benefit analysis",
				"Mitigation plans",
			},
			ApplicableLayers: []string{"system", "risk_management"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-8.3",
			Title:       "AI System Lifecycle",
			Description: "Manage AI systems throughout their lifecycle.",
			Objectives: []string{
				"Manage full AI lifecycle",
				"Control development and deployment",
				"Manage decommissioning",
			},
			Activities: []string{
				"Define lifecycle stages",
				"Control development processes",
				"Manage deployment",
				"Plan decommissioning",
			},
			EvidenceTypes: []string{
				"Lifecycle documentation",
				"Development records",
				"Deployment records",
				"Decommissioning plans",
			},
			ApplicableLayers: []string{"system", "operations"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-8.4",
			Title:       "AI System Documentation",
			Description: "Maintain documentation for AI systems.",
			Objectives: []string{
				"Document AI systems",
				"Maintain technical records",
				"Enable system understanding",
			},
			Activities: []string{
				"Document system design",
				"Record training data sources",
				"Document model decisions",
				"Maintain version history",
			},
			EvidenceTypes: []string{
				"System documentation",
				"Training data records",
				"Model cards",
				"Version control records",
			},
			ApplicableLayers: []string{"system"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-8.5",
			Title:       "Data for AI Systems",
			Description: "Manage data used in AI systems appropriately.",
			Objectives: []string{
				"Manage AI data quality",
				"Ensure data appropriateness",
				"Control data usage",
			},
			Activities: []string{
				"Define data requirements",
				"Validate data quality",
				"Control data access",
				"Document data lineage",
			},
			EvidenceTypes: []string{
				"Data quality reports",
				"Data validation records",
				"Access control documentation",
				"Lineage documentation",
			},
			ApplicableLayers: []string{"data", "system"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-8.6",
			Title:       "Third-Party Considerations",
			Description: "Manage third-party AI components and services.",
			Objectives: []string{
				"Assess third-party AI",
				"Control vendor relationships",
				"Manage supply chain risks",
			},
			Activities: []string{
				"Evaluate third-party AI",
				"Establish vendor requirements",
				"Monitor vendor compliance",
				"Manage contracts",
			},
			EvidenceTypes: []string{
				"Vendor assessments",
				"Contract documentation",
				"Compliance monitoring records",
			},
			ApplicableLayers: []string{"supply_chain", "governance"},
		},

		// Clause 9: Performance Evaluation
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-9.1",
			Title:       "Monitoring, Measurement, Analysis and Evaluation",
			Description: "Determine what needs to be monitored and measured.",
			Objectives: []string{
				"Define monitoring requirements",
				"Measure AI performance",
				"Analyze results",
			},
			Activities: []string{
				"Define metrics and KPIs",
				"Implement monitoring",
				"Analyze performance data",
				"Report results",
			},
			EvidenceTypes: []string{
				"Monitoring procedures",
				"Measurement records",
				"Analysis reports",
			},
			ApplicableLayers: []string{"operations", "system"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-9.2",
			Title:       "Internal Audit",
			Description: "Conduct internal audits at planned intervals.",
			Objectives: []string{
				"Verify AIMS conformance",
				"Assess effectiveness",
				"Identify improvements",
			},
			Activities: []string{
				"Plan audit program",
				"Conduct audits",
				"Report findings",
				"Track corrective actions",
			},
			EvidenceTypes: []string{
				"Audit program",
				"Audit reports",
				"Corrective action records",
			},
			ApplicableLayers: []string{"governance"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-9.3",
			Title:       "Management Review",
			Description: "Review the AI management system at planned intervals.",
			Objectives: []string{
				"Review AIMS performance",
				"Ensure continuing suitability",
				"Drive improvement",
			},
			Activities: []string{
				"Schedule management reviews",
				"Prepare review inputs",
				"Conduct reviews",
				"Document decisions",
			},
			EvidenceTypes: []string{
				"Review schedule",
				"Review inputs",
				"Meeting minutes",
				"Action items",
			},
			ApplicableLayers: []string{"governance"},
		},

		// Clause 10: Improvement
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-10.1",
			Title:       "Nonconformity and Corrective Action",
			Description: "React to nonconformities and take corrective action.",
			Objectives: []string{
				"Address nonconformities",
				"Prevent recurrence",
				"Implement corrections",
			},
			Activities: []string{
				"Identify nonconformities",
				"Determine root causes",
				"Implement corrective actions",
				"Evaluate effectiveness",
			},
			EvidenceTypes: []string{
				"Nonconformity records",
				"Root cause analysis",
				"Corrective action records",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-10.2",
			Title:       "Continual Improvement",
			Description: "Continually improve the AI management system.",
			Objectives: []string{
				"Improve AIMS effectiveness",
				"Enhance AI governance",
				"Optimize AI operations",
			},
			Activities: []string{
				"Identify improvement opportunities",
				"Prioritize improvements",
				"Implement enhancements",
				"Measure improvement",
			},
			EvidenceTypes: []string{
				"Improvement plans",
				"Implementation records",
				"Improvement metrics",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},

		// Annex A Controls (Selected Key Controls)
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.2.2",
			Title:       "AI System Transparency",
			Description: "Provide appropriate transparency about AI system capabilities and limitations.",
			Objectives: []string{
				"Enable system understanding",
				"Disclose limitations",
				"Support informed decisions",
			},
			Activities: []string{
				"Document system capabilities",
				"Communicate limitations",
				"Provide user guidance",
			},
			EvidenceTypes: []string{
				"Transparency documentation",
				"User communications",
				"Limitation disclosures",
			},
			ApplicableLayers: []string{"system", "organization"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.2.3",
			Title:       "AI System Explainability",
			Description: "Provide explanations of AI system outputs where appropriate.",
			Objectives: []string{
				"Enable output understanding",
				"Support accountability",
				"Enable contestability",
			},
			Activities: []string{
				"Implement explanation mechanisms",
				"Document explanation approaches",
				"Validate explanation quality",
			},
			EvidenceTypes: []string{
				"Explanation mechanisms",
				"Documentation",
				"Validation records",
			},
			ApplicableLayers: []string{"system"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.3.2",
			Title:       "Bias Assessment and Mitigation",
			Description: "Assess and mitigate bias in AI systems.",
			Objectives: []string{
				"Identify bias sources",
				"Assess bias impacts",
				"Mitigate identified bias",
			},
			Activities: []string{
				"Conduct bias assessments",
				"Implement mitigation measures",
				"Monitor for emerging bias",
			},
			EvidenceTypes: []string{
				"Bias assessment reports",
				"Mitigation documentation",
				"Monitoring records",
			},
			ApplicableLayers: []string{"system", "data"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.4.4",
			Title:       "AI System Security",
			Description: "Implement security controls for AI systems.",
			Objectives: []string{
				"Protect AI systems",
				"Secure AI data",
				"Prevent adversarial attacks",
			},
			Activities: []string{
				"Implement access controls",
				"Protect training data",
				"Monitor for attacks",
				"Test system robustness",
			},
			EvidenceTypes: []string{
				"Security controls documentation",
				"Access logs",
				"Penetration test results",
				"Monitoring records",
			},
			ApplicableLayers: []string{"system", "security"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.5.2",
			Title:       "Human Oversight",
			Description: "Implement appropriate human oversight of AI systems.",
			Objectives: []string{
				"Enable human control",
				"Support intervention",
				"Maintain accountability",
			},
			Activities: []string{
				"Define oversight mechanisms",
				"Implement intervention capabilities",
				"Train oversight personnel",
			},
			EvidenceTypes: []string{
				"Oversight procedures",
				"Intervention logs",
				"Training records",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.6.2",
			Title:       "AI System Reliability",
			Description: "Ensure AI systems perform reliably as intended.",
			Objectives: []string{
				"Ensure consistent performance",
				"Manage failures gracefully",
				"Maintain availability",
			},
			Activities: []string{
				"Define reliability requirements",
				"Test system reliability",
				"Implement fallback mechanisms",
			},
			EvidenceTypes: []string{
				"Reliability requirements",
				"Test results",
				"Fallback documentation",
			},
			ApplicableLayers: []string{"system", "operations"},
		},
		{
			FrameworkID: string(FrameworkISO42001),
			ControlID:   "ISO42001-A.7.3",
			Title:       "Privacy Protection",
			Description: "Protect privacy in AI system development and operation.",
			Objectives: []string{
				"Protect personal data",
				"Implement privacy by design",
				"Enable data subject rights",
			},
			Activities: []string{
				"Conduct privacy impact assessments",
				"Implement privacy controls",
				"Manage data subject requests",
			},
			EvidenceTypes: []string{
				"Privacy impact assessments",
				"Privacy controls documentation",
				"Data subject request logs",
			},
			ApplicableLayers: []string{"data", "governance"},
		},
	}
}
