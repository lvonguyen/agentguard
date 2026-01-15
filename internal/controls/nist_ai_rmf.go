package controls

import "github.com/agentguard/agentguard/internal/models"

// getNISTAIRMFControls returns NIST AI RMF control definitions.
func getNISTAIRMFControls() []models.Control {
	return []models.Control{
		// GOVERN Function
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "GOVERN-1",
			Title:       "Legal and Regulatory Requirements",
			Description: "Legal and regulatory requirements involving AI are understood, managed, and documented.",
			Objectives: []string{
				"Understand applicable laws and regulations",
				"Document compliance requirements",
				"Establish monitoring processes",
			},
			Activities: []string{
				"Identify applicable AI regulations",
				"Map regulatory requirements to AI systems",
				"Establish compliance monitoring",
			},
			EvidenceTypes: []string{
				"Regulatory mapping documentation",
				"Compliance assessment reports",
				"Legal review records",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "GOVERN-2",
			Title:       "Accountability Structures",
			Description: "Accountability structures are in place so that the appropriate teams and individuals are empowered, responsible, and trained for mapping, measuring, and managing AI risks.",
			Objectives: []string{
				"Define roles and responsibilities for AI risk management",
				"Establish accountability mechanisms",
				"Ensure adequate training",
			},
			Activities: []string{
				"Define AI risk management roles",
				"Assign responsibilities",
				"Develop training programs",
				"Establish escalation procedures",
			},
			EvidenceTypes: []string{
				"RACI matrices",
				"Job descriptions",
				"Training records",
				"Org charts",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "GOVERN-3",
			Title:       "Workforce Diversity and AI Impact",
			Description: "Workforce diversity, equity, inclusion, and accessibility processes are prioritized in the mapping, measuring, and managing of AI risks throughout the lifecycle.",
			Objectives: []string{
				"Promote diverse perspectives in AI development",
				"Assess AI impact on diverse populations",
			},
			Activities: []string{
				"Include diverse stakeholders in AI development",
				"Assess bias and fairness",
				"Document accessibility considerations",
			},
			EvidenceTypes: []string{
				"Diversity metrics",
				"Bias assessment reports",
				"Stakeholder engagement records",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "GOVERN-4",
			Title:       "Organizational Risk Tolerance",
			Description: "Organizational teams are committed to a culture that considers and communicates AI risk.",
			Objectives: []string{
				"Define AI risk tolerance",
				"Communicate risk culture",
			},
			Activities: []string{
				"Establish risk tolerance thresholds",
				"Communicate risk policies",
				"Monitor risk culture",
			},
			EvidenceTypes: []string{
				"Risk appetite statements",
				"Communication records",
				"Culture assessment results",
			},
			ApplicableLayers: []string{"governance", "risk_management"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "GOVERN-5",
			Title:       "Processes for Managing AI Risks",
			Description: "Processes are in place for robust engagement with relevant AI actors.",
			Objectives: []string{
				"Establish stakeholder engagement",
				"Define feedback mechanisms",
			},
			Activities: []string{
				"Identify relevant stakeholders",
				"Establish engagement processes",
				"Collect and act on feedback",
			},
			EvidenceTypes: []string{
				"Stakeholder maps",
				"Engagement records",
				"Feedback logs",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "GOVERN-6",
			Title:       "Policies and Procedures",
			Description: "Policies and procedures are in place to address AI risks and benefits arising from third-party software and data and from other supply chain risks.",
			Objectives: []string{
				"Manage third-party AI risks",
				"Establish supply chain security",
			},
			Activities: []string{
				"Assess third-party AI components",
				"Establish vendor requirements",
				"Monitor supply chain risks",
			},
			EvidenceTypes: []string{
				"Vendor assessments",
				"Supply chain policies",
				"Third-party audit reports",
			},
			ApplicableLayers: []string{"governance", "supply_chain"},
		},

		// MAP Function
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MAP-1",
			Title:       "Intended Purposes and Uses",
			Description: "Context is established and understood.",
			Objectives: []string{
				"Define AI system purpose",
				"Identify use cases and users",
				"Document operational context",
			},
			Activities: []string{
				"Document intended use cases",
				"Identify end users and affected parties",
				"Analyze operational environment",
			},
			EvidenceTypes: []string{
				"Use case documentation",
				"User personas",
				"Context analysis",
			},
			ApplicableLayers: []string{"system", "application"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MAP-2",
			Title:       "Categorization of AI System",
			Description: "Categorization of the AI system is performed.",
			Objectives: []string{
				"Classify AI system by risk level",
				"Determine applicable requirements",
			},
			Activities: []string{
				"Perform risk categorization",
				"Apply classification criteria",
				"Document categorization rationale",
			},
			EvidenceTypes: []string{
				"Risk categorization records",
				"Classification documentation",
			},
			ApplicableLayers: []string{"system", "risk_management"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MAP-3",
			Title:       "AI Capabilities and Limitations",
			Description: "AI capabilities, targeted usage, goals, and expected benefits and costs are understood.",
			Objectives: []string{
				"Document AI capabilities",
				"Identify limitations",
				"Assess benefits and costs",
			},
			Activities: []string{
				"Define system capabilities",
				"Document known limitations",
				"Perform cost-benefit analysis",
			},
			EvidenceTypes: []string{
				"Capability documentation",
				"Limitation disclosure",
				"Cost-benefit analysis",
			},
			ApplicableLayers: []string{"system", "application"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MAP-4",
			Title:       "Risks and Benefits",
			Description: "Risks and benefits are mapped for all components of the AI system including third-party software and data.",
			Objectives: []string{
				"Identify system-wide risks",
				"Map benefits to stakeholders",
			},
			Activities: []string{
				"Conduct risk assessment",
				"Map benefits by stakeholder",
				"Document trade-offs",
			},
			EvidenceTypes: []string{
				"Risk assessment reports",
				"Benefit mapping",
				"Trade-off analysis",
			},
			ApplicableLayers: []string{"system", "risk_management"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MAP-5",
			Title:       "Impacts to Individuals and Communities",
			Description: "Impacts to individuals, groups, communities, organizations, and society are characterized.",
			Objectives: []string{
				"Assess societal impact",
				"Identify affected communities",
			},
			Activities: []string{
				"Conduct impact assessment",
				"Engage affected communities",
				"Document potential harms",
			},
			EvidenceTypes: []string{
				"Impact assessment reports",
				"Community engagement records",
				"Harm documentation",
			},
			ApplicableLayers: []string{"society", "organization"},
		},

		// MEASURE Function
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MEASURE-1",
			Title:       "AI Risks Measured and Monitored",
			Description: "Appropriate methods and metrics are identified and applied.",
			Objectives: []string{
				"Define risk metrics",
				"Establish measurement methods",
			},
			Activities: []string{
				"Define KPIs and KRIs",
				"Implement measurement tools",
				"Establish baselines",
			},
			EvidenceTypes: []string{
				"Metrics definitions",
				"Measurement procedures",
				"Baseline documentation",
			},
			ApplicableLayers: []string{"system", "risk_management"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MEASURE-2",
			Title:       "AI Systems Evaluated",
			Description: "AI systems are evaluated for trustworthy characteristics.",
			Objectives: []string{
				"Evaluate trustworthiness",
				"Assess reliability and safety",
			},
			Activities: []string{
				"Conduct trustworthiness evaluation",
				"Test for bias and fairness",
				"Assess explainability",
			},
			EvidenceTypes: []string{
				"Evaluation reports",
				"Bias testing results",
				"Explainability assessments",
			},
			ApplicableLayers: []string{"system", "testing"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MEASURE-3",
			Title:       "Mechanisms for Tracking Metrics",
			Description: "Mechanisms for tracking identified AI risks over time are in place.",
			Objectives: []string{
				"Implement risk tracking",
				"Enable trend analysis",
			},
			Activities: []string{
				"Deploy monitoring systems",
				"Configure alerting",
				"Establish reporting cadence",
			},
			EvidenceTypes: []string{
				"Monitoring dashboards",
				"Alert configurations",
				"Trend reports",
			},
			ApplicableLayers: []string{"system", "operations"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MEASURE-4",
			Title:       "Feedback Incorporated",
			Description: "Feedback about efficacy of measurement is gathered and assessed.",
			Objectives: []string{
				"Collect measurement feedback",
				"Improve measurement processes",
			},
			Activities: []string{
				"Gather stakeholder feedback",
				"Review measurement effectiveness",
				"Update metrics as needed",
			},
			EvidenceTypes: []string{
				"Feedback records",
				"Review documentation",
				"Metrics update logs",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},

		// MANAGE Function
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MANAGE-1",
			Title:       "AI Risk Prioritization",
			Description: "AI risks based on assessments and other analytical output from the MAP and MEASURE functions are prioritized, responded to, and managed.",
			Objectives: []string{
				"Prioritize identified risks",
				"Develop response strategies",
			},
			Activities: []string{
				"Rank risks by severity",
				"Develop mitigation plans",
				"Allocate resources",
			},
			EvidenceTypes: []string{
				"Risk prioritization matrix",
				"Mitigation plans",
				"Resource allocation records",
			},
			ApplicableLayers: []string{"governance", "risk_management"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MANAGE-2",
			Title:       "Strategies to Maximize AI Benefits",
			Description: "Strategies to maximize AI benefits and minimize negative impacts are planned, prepared, implemented, documented, and informed by input from relevant AI actors.",
			Objectives: []string{
				"Maximize AI benefits",
				"Minimize negative impacts",
			},
			Activities: []string{
				"Develop optimization strategies",
				"Implement harm reduction measures",
				"Document outcomes",
			},
			EvidenceTypes: []string{
				"Strategy documentation",
				"Implementation records",
				"Outcome assessments",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MANAGE-3",
			Title:       "AI Risk Management Integrated",
			Description: "AI risk management is integrated into broader organizational risk management.",
			Objectives: []string{
				"Integrate with enterprise risk management",
				"Align with organizational processes",
			},
			Activities: []string{
				"Map AI risks to ERM framework",
				"Establish reporting integration",
				"Align governance structures",
			},
			EvidenceTypes: []string{
				"ERM integration documentation",
				"Reporting workflows",
				"Governance alignment records",
			},
			ApplicableLayers: []string{"governance", "organization"},
		},
		{
			FrameworkID: string(FrameworkNISTAIRMF),
			ControlID:   "MANAGE-4",
			Title:       "Risk Treatments Documented",
			Description: "Risk treatments, including response and recovery, and communication plans for the identified and measured AI risks are documented and monitored regularly.",
			Objectives: []string{
				"Document risk treatments",
				"Establish recovery plans",
			},
			Activities: []string{
				"Develop treatment plans",
				"Create recovery procedures",
				"Establish communication protocols",
			},
			EvidenceTypes: []string{
				"Treatment documentation",
				"Recovery plans",
				"Communication procedures",
			},
			ApplicableLayers: []string{"governance", "operations"},
		},
	}
}
