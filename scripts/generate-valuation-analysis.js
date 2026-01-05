const { Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell, 
        AlignmentType, HeadingLevel, BorderStyle, WidthType, ShadingType,
        Header, Footer, PageNumber, LevelFormat, PageBreak } = require('docx');
const fs = require('fs');
const path = require('path');

// Style configuration
const styles = {
  default: { document: { run: { font: "Arial", size: 22 } } }, // 11pt default
  paragraphStyles: [
    { id: "Title", name: "Title", basedOn: "Normal",
      run: { size: 48, bold: true, color: "1a365d", font: "Georgia" },
      paragraph: { spacing: { before: 0, after: 240 }, alignment: AlignmentType.CENTER } },
    { id: "Heading1", name: "Heading 1", basedOn: "Normal", next: "Normal", quickFormat: true,
      run: { size: 32, bold: true, color: "1a365d", font: "Georgia" },
      paragraph: { spacing: { before: 360, after: 180 }, outlineLevel: 0 } },
    { id: "Heading2", name: "Heading 2", basedOn: "Normal", next: "Normal", quickFormat: true,
      run: { size: 26, bold: true, color: "2d5a87", font: "Georgia" },
      paragraph: { spacing: { before: 280, after: 140 }, outlineLevel: 1 } },
    { id: "Heading3", name: "Heading 3", basedOn: "Normal", next: "Normal", quickFormat: true,
      run: { size: 24, bold: true, color: "3d7db3", font: "Arial" },
      paragraph: { spacing: { before: 200, after: 100 }, outlineLevel: 2 } }
  ]
};

// Numbering configuration
const numbering = {
  config: [
    { reference: "bullet-list",
      levels: [{ level: 0, format: LevelFormat.BULLET, text: "•", alignment: AlignmentType.LEFT,
        style: { paragraph: { indent: { left: 720, hanging: 360 } } } }] },
    { reference: "numbered-1",
      levels: [{ level: 0, format: LevelFormat.DECIMAL, text: "%1.", alignment: AlignmentType.LEFT,
        style: { paragraph: { indent: { left: 720, hanging: 360 } } } }] },
    { reference: "numbered-2",
      levels: [{ level: 0, format: LevelFormat.DECIMAL, text: "%1.", alignment: AlignmentType.LEFT,
        style: { paragraph: { indent: { left: 720, hanging: 360 } } } }] }
  ]
};

// Table helpers
const tableBorder = { style: BorderStyle.SINGLE, size: 1, color: "CCCCCC" };
const cellBorders = { top: tableBorder, bottom: tableBorder, left: tableBorder, right: tableBorder };
const headerShading = { fill: "1a365d", type: ShadingType.CLEAR };
const altRowShading = { fill: "f0f4f8", type: ShadingType.CLEAR };

function createHeaderCell(text, width) {
  return new TableCell({
    borders: cellBorders,
    width: { size: width, type: WidthType.DXA },
    shading: headerShading,
    children: [new Paragraph({ 
      alignment: AlignmentType.CENTER,
      children: [new TextRun({ text, bold: true, color: "FFFFFF", size: 20 })]
    })]
  });
}

function createDataCell(text, width, isAlt = false, align = AlignmentType.LEFT) {
  return new TableCell({
    borders: cellBorders,
    width: { size: width, type: WidthType.DXA },
    shading: isAlt ? altRowShading : undefined,
    children: [new Paragraph({ 
      alignment: align,
      children: [new TextRun({ text, size: 20 })]
    })]
  });
}

// Document content
const children = [
  // Title
  new Paragraph({ heading: HeadingLevel.TITLE, children: [new TextRun("AgentGuard")] }),
  new Paragraph({ 
    alignment: AlignmentType.CENTER,
    spacing: { after: 480 },
    children: [new TextRun({ text: "SaaS Valuation & Investment Analysis", size: 28, color: "666666" })]
  }),
  new Paragraph({ 
    alignment: AlignmentType.CENTER,
    spacing: { after: 480 },
    children: [new TextRun({ text: "Series A Investment Memorandum • January 2026", size: 22, color: "888888", italics: true })]
  }),

  // Executive Summary
  new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("Executive Summary")] }),
  new Paragraph({ children: [new TextRun("AgentGuard is an AI security governance platform purpose-built for enterprise agentic AI deployments. As organizations rapidly adopt autonomous AI agents that browse, execute code, and make decisions, they face an unguarded security frontier. AgentGuard provides comprehensive control mapping, threat modeling, and policy enforcement—with the market's only NIST AI RMF to 800-53 crosswalk.")] }),
  new Paragraph({ spacing: { before: 200 }, children: [new TextRun({ text: "Investment Highlights:", bold: true })] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun("First-mover advantage in agentic AI security governance")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun("$4.3B TAM growing 34% CAGR through 2028")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun("FedRAMP-ready architecture targeting $890M federal market")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun("Enterprise SaaS model with $85K target ACV")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun("Seeking $12M Series A at $48M pre-money valuation")] }),

  // Market Analysis
  new Paragraph({ children: [new PageBreak()] }),
  new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("Market Analysis")] }),
  
  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Total Addressable Market")] }),
  new Paragraph({ children: [new TextRun("The AI security market is experiencing explosive growth as enterprises adopt generative AI and agentic systems. Market sizing reflects the convergence of AI security, GRC compliance, and observability platforms.")] }),
  
  new Table({
    columnWidths: [2340, 2340, 2340, 2340],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Segment", 2340),
        createHeaderCell("2025 Size", 2340),
        createHeaderCell("2028 Size", 2340),
        createHeaderCell("CAGR", 2340)
      ]}),
      new TableRow({ children: [
        createDataCell("TAM - AI Security", 2340),
        createDataCell("$2.1B", 2340, false, AlignmentType.RIGHT),
        createDataCell("$4.3B", 2340, false, AlignmentType.RIGHT),
        createDataCell("34%", 2340, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("SAM - Enterprise AI Governance", 2340, true),
        createDataCell("$680M", 2340, true, AlignmentType.RIGHT),
        createDataCell("$1.4B", 2340, true, AlignmentType.RIGHT),
        createDataCell("38%", 2340, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("SOM - Agentic AI Security", 2340),
        createDataCell("$180M", 2340, false, AlignmentType.RIGHT),
        createDataCell("$420M", 2340, false, AlignmentType.RIGHT),
        createDataCell("42%", 2340, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Federal/FedRAMP Segment", 2340, true),
        createDataCell("$340M", 2340, true, AlignmentType.RIGHT),
        createDataCell("$890M", 2340, true, AlignmentType.RIGHT),
        createDataCell("48%", 2340, true, AlignmentType.RIGHT)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 400 }, children: [new TextRun("Market Drivers")] }),
  new Paragraph({ numbering: { reference: "numbered-1", level: 0 }, children: [new TextRun({ text: "Regulatory Pressure: ", bold: true }), new TextRun("EU AI Act, NIST AI RMF, and emerging state laws mandate AI governance controls")] }),
  new Paragraph({ numbering: { reference: "numbered-1", level: 0 }, children: [new TextRun({ text: "Enterprise Adoption: ", bold: true }), new TextRun("77% of enterprises deploying agentic AI lack adequate security controls")] }),
  new Paragraph({ numbering: { reference: "numbered-1", level: 0 }, children: [new TextRun({ text: "Breach Costs: ", bold: true }), new TextRun("Average AI-related breach costs $4.2M, 23% higher than traditional breaches")] }),
  new Paragraph({ numbering: { reference: "numbered-1", level: 0 }, children: [new TextRun({ text: "FedRAMP Demand: ", bold: true }), new TextRun("Federal AI modernization creating $890M opportunity by 2028")] }),

  // Financial Projections
  new Paragraph({ children: [new PageBreak()] }),
  new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("Financial Projections")] }),
  
  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("36-Month Revenue Forecast")] }),
  new Table({
    columnWidths: [1872, 1872, 1872, 1872, 1872],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Metric", 1872),
        createHeaderCell("Y1 (2026)", 1872),
        createHeaderCell("Y2 (2027)", 1872),
        createHeaderCell("Y3 (2028)", 1872),
        createHeaderCell("Y3 Exit", 1872)
      ]}),
      new TableRow({ children: [
        createDataCell("ARR", 1872),
        createDataCell("$1.8M", 1872, false, AlignmentType.RIGHT),
        createDataCell("$7.2M", 1872, false, AlignmentType.RIGHT),
        createDataCell("$21.6M", 1872, false, AlignmentType.RIGHT),
        createDataCell("$28M", 1872, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Customers", 1872, true),
        createDataCell("15", 1872, true, AlignmentType.RIGHT),
        createDataCell("52", 1872, true, AlignmentType.RIGHT),
        createDataCell("145", 1872, true, AlignmentType.RIGHT),
        createDataCell("180+", 1872, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("ACV", 1872),
        createDataCell("$72K", 1872, false, AlignmentType.RIGHT),
        createDataCell("$85K", 1872, false, AlignmentType.RIGHT),
        createDataCell("$95K", 1872, false, AlignmentType.RIGHT),
        createDataCell("$100K+", 1872, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("NRR", 1872, true),
        createDataCell("115%", 1872, true, AlignmentType.RIGHT),
        createDataCell("125%", 1872, true, AlignmentType.RIGHT),
        createDataCell("130%", 1872, true, AlignmentType.RIGHT),
        createDataCell("135%", 1872, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Gross Margin", 1872),
        createDataCell("72%", 1872, false, AlignmentType.RIGHT),
        createDataCell("78%", 1872, false, AlignmentType.RIGHT),
        createDataCell("82%", 1872, false, AlignmentType.RIGHT),
        createDataCell("85%", 1872, false, AlignmentType.RIGHT)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 400 }, children: [new TextRun("Unit Economics")] }),
  new Table({
    columnWidths: [4680, 4680],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Metric", 4680),
        createHeaderCell("Target", 4680)
      ]}),
      new TableRow({ children: [
        createDataCell("Customer Acquisition Cost (CAC)", 4680),
        createDataCell("$45,000", 4680, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Annual Contract Value (ACV)", 4680, true),
        createDataCell("$85,000", 4680, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Lifetime Value (LTV)", 4680),
        createDataCell("$340,000", 4680, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("LTV:CAC Ratio", 4680, true),
        createDataCell("7.6x", 4680, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("CAC Payback", 4680),
        createDataCell("18 months", 4680, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Annual Churn", 4680, true),
        createDataCell("<5%", 4680, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Magic Number", 4680),
        createDataCell("1.2", 4680, false, AlignmentType.RIGHT)
      ]})
    ]
  }),

  // Competitive Positioning
  new Paragraph({ children: [new PageBreak()] }),
  new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("Competitive Positioning")] }),
  
  new Paragraph({ children: [new TextRun("AgentGuard operates at the intersection of AI security, GRC compliance, and LLM observability—a unique position that existing players cannot easily replicate.")] }),
  
  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Competitive Matrix")] }),
  new Table({
    columnWidths: [2808, 1638, 1638, 1638, 1638],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Capability", 2808),
        createHeaderCell("Lakera", 1638),
        createHeaderCell("LangSmith", 1638),
        createHeaderCell("Arize", 1638),
        createHeaderCell("AgentGuard", 1638)
      ]}),
      new TableRow({ children: [
        createDataCell("NIST AI RMF Crosswalk", 2808),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("✓", 1638, false, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("Tool Access Policies", 2808, true),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("Partial", 1638, true, AlignmentType.CENTER),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("Data Flow Governance", 2808),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("✓", 1638, false, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("Human-in-the-Loop Gates", 2808, true),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("Agent Threat Modeling", 2808),
        createDataCell("Partial", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("✓", 1638, false, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("LLM Observability", 2808, true),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("Content Safety", 2808),
        createDataCell("✓", 1638, false, AlignmentType.CENTER),
        createDataCell("—", 1638, false, AlignmentType.CENTER),
        createDataCell("Partial", 1638, false, AlignmentType.CENTER),
        createDataCell("✓", 1638, false, AlignmentType.CENTER)
      ]}),
      new TableRow({ children: [
        createDataCell("FedRAMP Ready", 2808, true),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("—", 1638, true, AlignmentType.CENTER),
        createDataCell("✓", 1638, true, AlignmentType.CENTER)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 400 }, children: [new TextRun("Sustainable Moats")] }),
  new Paragraph({ numbering: { reference: "numbered-2", level: 0 }, children: [new TextRun({ text: "Compliance IP: ", bold: true }), new TextRun("Only vendor with complete NIST AI RMF to 800-53 control mapping")] }),
  new Paragraph({ numbering: { reference: "numbered-2", level: 0 }, children: [new TextRun({ text: "Agent-Native Architecture: ", bold: true }), new TextRun("Built for autonomous agents, not retrofitted from LLM guardrails")] }),
  new Paragraph({ numbering: { reference: "numbered-2", level: 0 }, children: [new TextRun({ text: "FedRAMP Head Start: ", bold: true }), new TextRun("18-24 month lead time for federal market access")] }),
  new Paragraph({ numbering: { reference: "numbered-2", level: 0 }, children: [new TextRun({ text: "Data Network Effects: ", bold: true }), new TextRun("Threat intelligence improves with scale across customer base")] }),

  // Valuation Model
  new Paragraph({ children: [new PageBreak()] }),
  new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("Valuation Model")] }),
  
  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Series A Terms")] }),
  new Table({
    columnWidths: [4680, 4680],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Term", 4680),
        createHeaderCell("Value", 4680)
      ]}),
      new TableRow({ children: [
        createDataCell("Raise Amount", 4680),
        createDataCell("$12,000,000", 4680, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Pre-Money Valuation", 4680, true),
        createDataCell("$48,000,000", 4680, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Post-Money Valuation", 4680),
        createDataCell("$60,000,000", 4680, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Investor Ownership", 4680, true),
        createDataCell("20.0%", 4680, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Implied ARR Multiple (at Y1 exit)", 4680),
        createDataCell("26.7x", 4680, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Option Pool Increase", 4680, true),
        createDataCell("10%", 4680, true, AlignmentType.RIGHT)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 400 }, children: [new TextRun("Use of Funds")] }),
  new Table({
    columnWidths: [4680, 2340, 2340],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Category", 4680),
        createHeaderCell("Amount", 2340),
        createHeaderCell("% of Raise", 2340)
      ]}),
      new TableRow({ children: [
        createDataCell("Engineering & Product", 4680),
        createDataCell("$5,400,000", 2340, false, AlignmentType.RIGHT),
        createDataCell("45%", 2340, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Sales & Marketing", 4680, true),
        createDataCell("$3,600,000", 2340, true, AlignmentType.RIGHT),
        createDataCell("30%", 2340, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("FedRAMP Certification", 4680),
        createDataCell("$1,200,000", 2340, false, AlignmentType.RIGHT),
        createDataCell("10%", 2340, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("G&A / Operations", 4680, true),
        createDataCell("$1,080,000", 2340, true, AlignmentType.RIGHT),
        createDataCell("9%", 2340, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Working Capital Reserve", 4680),
        createDataCell("$720,000", 2340, false, AlignmentType.RIGHT),
        createDataCell("6%", 2340, false, AlignmentType.RIGHT)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 400 }, children: [new TextRun("Exit Scenarios")] }),
  new Table({
    columnWidths: [2340, 2340, 2340, 2340],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Scenario", 2340),
        createHeaderCell("Y3 ARR", 2340),
        createHeaderCell("Exit Multiple", 2340),
        createHeaderCell("Exit Value", 2340)
      ]}),
      new TableRow({ children: [
        createDataCell("Base Case", 2340),
        createDataCell("$21.6M", 2340, false, AlignmentType.RIGHT),
        createDataCell("12x ARR", 2340, false, AlignmentType.RIGHT),
        createDataCell("$259M", 2340, false, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Bull Case", 2340, true),
        createDataCell("$28M", 2340, true, AlignmentType.RIGHT),
        createDataCell("15x ARR", 2340, true, AlignmentType.RIGHT),
        createDataCell("$420M", 2340, true, AlignmentType.RIGHT)
      ]}),
      new TableRow({ children: [
        createDataCell("Bear Case", 2340),
        createDataCell("$14M", 2340, false, AlignmentType.RIGHT),
        createDataCell("8x ARR", 2340, false, AlignmentType.RIGHT),
        createDataCell("$112M", 2340, false, AlignmentType.RIGHT)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_3, spacing: { before: 300 }, children: [new TextRun("Return Analysis (Base Case)")] }),
  new Paragraph({ children: [new TextRun("At 20% ownership and $259M exit value, Series A investors would receive $51.8M—a 4.3x return on $12M invested over 36 months (62% IRR).")] }),

  // Investment Thesis
  new Paragraph({ children: [new PageBreak()] }),
  new Paragraph({ heading: HeadingLevel.HEADING_1, children: [new TextRun("Investment Thesis")] }),
  
  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Why Now")] }),
  new Paragraph({ children: [new TextRun("The agentic AI security market is at an inflection point. Enterprise adoption is accelerating while security controls remain nascent. Regulatory frameworks (NIST AI RMF, EU AI Act) are creating compliance mandates that existing vendors cannot address. AgentGuard is positioned to capture this market at the precise moment demand crystallizes.")] }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Why AgentGuard")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun({ text: "Technical Differentiation: ", bold: true }), new TextRun("Agent-native architecture with unified governance, observability, and enforcement")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun({ text: "Compliance Moat: ", bold: true }), new TextRun("Only vendor with complete NIST AI RMF to NIST 800-53 control crosswalk")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun({ text: "Federal Opportunity: ", bold: true }), new TextRun("FedRAMP-ready architecture targeting $890M government market")] }),
  new Paragraph({ numbering: { reference: "bullet-list", level: 0 }, children: [new TextRun({ text: "Enterprise DNA: ", bold: true }), new TextRun("Team with deep security, compliance, and enterprise SaaS experience")] }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, children: [new TextRun("Key Risks & Mitigations")] }),
  new Table({
    columnWidths: [4680, 4680],
    rows: [
      new TableRow({ tableHeader: true, children: [
        createHeaderCell("Risk", 4680),
        createHeaderCell("Mitigation", 4680)
      ]}),
      new TableRow({ children: [
        createDataCell("Market timing (agentic AI adoption)", 4680),
        createDataCell("Multi-modal product supports current LLM + future agent workloads", 4680)
      ]}),
      new TableRow({ children: [
        createDataCell("Platform vendor competition (AWS, Azure)", 4680, true),
        createDataCell("Multi-cloud neutrality; deeper compliance expertise than platform teams", 4680, true)
      ]}),
      new TableRow({ children: [
        createDataCell("FedRAMP timeline uncertainty", 4680),
        createDataCell("Commercial traction funds operations during 18-24 month certification", 4680)
      ]}),
      new TableRow({ children: [
        createDataCell("Team scaling in competitive market", 4680, true),
        createDataCell("Remote-first culture; equity compensation; mission-driven positioning", 4680, true)
      ]})
    ]
  }),

  new Paragraph({ heading: HeadingLevel.HEADING_2, spacing: { before: 400 }, children: [new TextRun("Conclusion")] }),
  new Paragraph({ children: [new TextRun("AgentGuard represents a compelling Series A opportunity in the rapidly emerging AI security governance market. With differentiated technology, clear regulatory tailwinds, and a capital-efficient go-to-market strategy, the company is positioned to achieve market leadership and deliver strong investor returns.")] }),
  
  new Paragraph({ spacing: { before: 400 }, alignment: AlignmentType.CENTER, children: [
    new TextRun({ text: "Contact: ", bold: true }),
    new TextRun("founders@agentguard.ai")
  ]})
];

// Create document
const doc = new Document({
  styles,
  numbering,
  sections: [{
    properties: {
      page: { margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 } }
    },
    headers: {
      default: new Header({ children: [new Paragraph({ 
        alignment: AlignmentType.RIGHT,
        children: [new TextRun({ text: "AgentGuard • Confidential", size: 18, color: "888888" })]
      })] })
    },
    footers: {
      default: new Footer({ children: [new Paragraph({ 
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "Page ", size: 18, color: "888888" }), 
                   new TextRun({ children: [PageNumber.CURRENT], size: 18, color: "888888" })]
      })] })
    },
    children
  }]
});

// Generate and save
const outputPath = path.join(__dirname, '..', 'docs', 'AgentGuard-Valuation-Analysis.docx');
Packer.toBuffer(doc).then(buffer => {
  fs.writeFileSync(outputPath, buffer);
  console.log(`Valuation analysis saved to: ${outputPath}`);
}).catch(err => {
  console.error('Error:', err);
  process.exit(1);
});
