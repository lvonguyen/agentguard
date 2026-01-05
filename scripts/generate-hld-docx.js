const fs = require('fs');
const { Document, Packer, Paragraph, TextRun, Table, TableRow, TableCell, Header, Footer, 
        AlignmentType, PageNumber, HeadingLevel, BorderStyle, WidthType, ShadingType,
        LevelFormat, TableOfContents, PageBreak } = require('docx');

// AgentGuard HLD Document - Professional formatting with Georgia font

const tableBorder = { style: BorderStyle.SINGLE, size: 1, color: "CCCCCC" };
const cellBorders = { top: tableBorder, bottom: tableBorder, left: tableBorder, right: tableBorder };
const headerShading = { fill: "1F4E79", type: ShadingType.CLEAR };

const doc = new Document({
  styles: {
    default: {
      document: {
        run: { font: "Georgia", size: 22 } // 11pt default
      }
    },
    paragraphStyles: [
      {
        id: "Title",
        name: "Title",
        basedOn: "Normal",
        run: { size: 56, bold: true, color: "1F4E79", font: "Georgia" },
        paragraph: { spacing: { before: 0, after: 200 }, alignment: AlignmentType.CENTER }
      },
      {
        id: "Heading1",
        name: "Heading 1",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 32, bold: true, color: "1F4E79", font: "Georgia" },
        paragraph: { spacing: { before: 400, after: 200 }, outlineLevel: 0 }
      },
      {
        id: "Heading2",
        name: "Heading 2",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 26, bold: true, color: "2E75B6", font: "Georgia" },
        paragraph: { spacing: { before: 300, after: 150 }, outlineLevel: 1 }
      },
      {
        id: "Heading3",
        name: "Heading 3",
        basedOn: "Normal",
        next: "Normal",
        quickFormat: true,
        run: { size: 24, bold: true, color: "404040", font: "Georgia" },
        paragraph: { spacing: { before: 200, after: 100 }, outlineLevel: 2 }
      }
    ]
  },
  numbering: {
    config: [
      {
        reference: "main-bullets",
        levels: [{
          level: 0,
          format: LevelFormat.BULLET,
          text: "•",
          alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 720, hanging: 360 } } }
        }]
      },
      {
        reference: "sub-bullets",
        levels: [{
          level: 0,
          format: LevelFormat.BULLET,
          text: "○",
          alignment: AlignmentType.LEFT,
          style: { paragraph: { indent: { left: 1080, hanging: 360 } } }
        }]
      }
    ]
  },
  sections: [{
    properties: {
      page: {
        margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 }
      }
    },
    headers: {
      default: new Header({
        children: [
          new Paragraph({
            alignment: AlignmentType.RIGHT,
            children: [
              new TextRun({ text: "AgentGuard High-Level Design", font: "Georgia", size: 18, color: "808080" })
            ]
          })
        ]
      })
    },
    footers: {
      default: new Footer({
        children: [
          new Paragraph({
            alignment: AlignmentType.CENTER,
            children: [
              new TextRun({ text: "Page ", font: "Georgia", size: 18 }),
              new TextRun({ children: [PageNumber.CURRENT], font: "Georgia", size: 18 }),
              new TextRun({ text: " of ", font: "Georgia", size: 18 }),
              new TextRun({ children: [PageNumber.TOTAL_PAGES], font: "Georgia", size: 18 }),
              new TextRun({ text: "  |  CONFIDENTIAL", font: "Georgia", size: 18, color: "808080" })
            ]
          })
        ]
      })
    },
    children: [
      // Title Page
      new Paragraph({ children: [] }),
      new Paragraph({ children: [] }),
      new Paragraph({ children: [] }),
      new Paragraph({
        heading: HeadingLevel.TITLE,
        children: [new TextRun({ text: "AgentGuard", font: "Georgia" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        spacing: { after: 400 },
        children: [new TextRun({ text: "High-Level Design Document", size: 36, font: "Georgia", color: "404040" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "AI Security Governance Framework", size: 28, font: "Georgia", italics: true, color: "666666" })]
      }),
      new Paragraph({ children: [] }),
      new Paragraph({ children: [] }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "Version 1.0", size: 24, font: "Georgia" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "January 2026", size: 24, font: "Georgia" })]
      }),
      new Paragraph({ children: [] }),
      new Paragraph({ children: [] }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "Author: Liem Vo-Nguyen", size: 22, font: "Georgia", color: "666666" })]
      }),
      new Paragraph({
        alignment: AlignmentType.CENTER,
        children: [new TextRun({ text: "Status: Draft", size: 22, font: "Georgia", color: "666666" })]
      }),
      
      // Page break before TOC
      new Paragraph({ children: [new PageBreak()] }),
      
      // Table of Contents
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: "Table of Contents", font: "Georgia" })]
      }),
      new TableOfContents("Table of Contents", { hyperlink: true, headingStyleRange: "1-3" }),
      
      // Page break before content
      new Paragraph({ children: [new PageBreak()] }),
      
      // Executive Summary
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: "1. Executive Summary", font: "Georgia" })]
      }),
      new Paragraph({
        spacing: { after: 200 },
        children: [new TextRun({ 
          text: "AgentGuard is an AI security governance framework that addresses the unique risks of agentic AI systems in enterprise environments. It provides control mapping to established compliance frameworks, runtime observability for agent execution chains, and policy-as-code guardrails—capabilities that do not exist in current vendor offerings.",
          font: "Georgia"
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "1.1 Key Differentiators", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "First comprehensive NIST AI RMF → NIST 800-53 crosswalk for FedRAMP-aligned organizations", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Agent-specific threat modeling beyond content safety (tool abuse, privilege escalation, data exfiltration)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Unified observability + policy enforcement in single platform", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "SDK middleware for major agent frameworks (LangChain, CrewAI, AutoGen)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Enterprise GRC integration for AI risk acceptance workflows", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "1.2 Target Users", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Security Architects evaluating AI deployment risks", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Compliance Teams mapping AI controls to existing frameworks", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Platform Engineers building internal AI platforms", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Auditors assessing AI governance posture", font: "Georgia" })]
      }),
      
      // Market Analysis
      new Paragraph({ children: [new PageBreak()] }),
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: "2. Market Analysis & Build/Buy Rationale", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "2.1 Vendor Landscape: LLM Observability", font: "Georgia" })]
      }),
      
      // Observability Vendor Table
      new Table({
        columnWidths: [1800, 1500, 2500, 2500, 1800],
        rows: [
          new TableRow({
            tableHeader: true,
            children: [
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Vendor", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Pricing", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Strengths", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Gaps", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Verdict", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
            ]
          }),
          new TableRow({
            children: [
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "LangSmith", bold: true, font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "$39/seat + usage", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Best LangChain integration, excellent trace UI", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Closed ecosystem, no self-host, no security signals", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Integrate for LC users", font: "Georgia", size: 20 })] })] }),
            ]
          }),
          new TableRow({
            children: [
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Langfuse", bold: true, font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "OSS / $59/seat", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Self-hostable, open source, cost tracking", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Security features immature, no compliance mapping", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: { fill: "E2EFDA", type: ShadingType.CLEAR }, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "PRIMARY", bold: true, font: "Georgia", size: 20, color: "375623" })] })] }),
            ]
          }),
          new TableRow({
            children: [
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Arize Phoenix", bold: true, font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "OSS / Enterprise", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Strong ML heritage, embedding drift detection", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "LLM features still maturing, complex setup", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Drift detection only", font: "Georgia", size: 20 })] })] }),
            ]
          }),
        ]
      }),
      
      new Paragraph({
        spacing: { before: 300, after: 200 },
        children: [new TextRun({ 
          text: "Decision: Integrate with Langfuse as primary observability backend. It's open source, self-hostable (important for compliance), and has clean APIs for extension. AgentGuard adds security-specific spans and enrichment.",
          font: "Georgia", italics: true
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "2.2 Vendor Landscape: AI Guardrails", font: "Georgia" })]
      }),
      
      // Guardrails Vendor Table
      new Table({
        columnWidths: [1800, 1500, 2500, 2500, 1800],
        rows: [
          new TableRow({
            tableHeader: true,
            children: [
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Vendor", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Pricing", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Strengths", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Gaps", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Verdict", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
            ]
          }),
          new TableRow({
            children: [
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Lakera Guard", bold: true, font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "$0.001/request", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Best prompt injection detection (<50ms, >95% accuracy)", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "No agent awareness, no tool-use controls", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: { fill: "E2EFDA", type: ShadingType.CLEAR }, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "INTEGRATE", bold: true, font: "Georgia", size: 20, color: "375623" })] })] }),
            ]
          }),
          new TableRow({
            children: [
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "AWS Guardrails", bold: true, font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "$0.75/1K units", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "Native AWS, content filters", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "AWS-only, no customization, no agent features", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, width: { size: 1800, type: WidthType.DXA },
                children: [new Paragraph({ children: [new TextRun({ text: "AWS deployments", font: "Georgia", size: 20 })] })] }),
            ]
          }),
        ]
      }),
      
      new Paragraph({
        spacing: { before: 300, after: 200 },
        children: [new TextRun({ 
          text: "Decision: Integrate Lakera Guard for prompt injection detection as a preprocessing layer. Build custom agent-specific policy engine on OPA for tool access control, as no vendor addresses this.",
          font: "Georgia", italics: true
        })]
      }),
      
      // Build vs Buy Summary
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "2.3 Build vs. Buy Summary", font: "Georgia" })]
      }),
      
      new Table({
        columnWidths: [2500, 2000, 5600],
        rows: [
          new TableRow({
            tableHeader: true,
            children: [
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Capability", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2000, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Decision", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 5600, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Rationale", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
            ]
          }),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Trace collection/storage", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "FFF2CC", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "BUY (Langfuse)", font: "Georgia", size: 20, bold: true })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Commodity capability, mature OSS solution", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Prompt injection detection", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "FFF2CC", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "BUY (Lakera)", font: "Georgia", size: 20, bold: true })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Specialized ML models, not core competency", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Control framework mapping", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "E2EFDA", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "BUILD", font: "Georgia", size: 20, bold: true, color: "375623" })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "No vendor solution exists, core differentiator", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Agent policy engine", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "E2EFDA", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "BUILD", font: "Georgia", size: 20, bold: true, color: "375623" })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "No vendor addresses agent-specific risks", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Threat modeling templates", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "E2EFDA", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "BUILD", font: "Georgia", size: 20, bold: true, color: "375623" })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Novel frameworks needed for agentic systems", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Maturity model", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "E2EFDA", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "BUILD", font: "Georgia", size: 20, bold: true, color: "375623" })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "AI-specific assessment methodology", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "GRC integration", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, shading: { fill: "DEEBF7", type: ShadingType.CLEAR }, children: [new Paragraph({ children: [new TextRun({ text: "EXTEND", font: "Georgia", size: 20, bold: true, color: "1F4E79" })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "APIs exist, need AI-specific workflows", font: "Georgia", size: 20 })] })] }),
          ]}),
        ]
      }),
      
      // Architecture Overview
      new Paragraph({ children: [new PageBreak()] }),
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: "3. Architecture Overview", font: "Georgia" })]
      }),
      
      new Paragraph({
        spacing: { after: 200 },
        children: [new TextRun({ 
          text: "AgentGuard follows a microservices architecture with five core services that communicate through a shared API gateway. The platform integrates with agent frameworks via SDK middleware that intercepts execution flows.",
          font: "Georgia"
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "3.1 Core Services", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_3,
        children: [new TextRun({ text: "Control Mapping Service", font: "Georgia" })]
      }),
      new Paragraph({
        spacing: { after: 150 },
        children: [new TextRun({ 
          text: "Manages control framework definitions (NIST AI RMF, 800-53, ISO 42001) and crosswalk mappings between frameworks. Provides gap analysis and compliance reporting.",
          font: "Georgia"
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_3,
        children: [new TextRun({ text: "Observability Service", font: "Georgia" })]
      }),
      new Paragraph({
        spacing: { after: 150 },
        children: [new TextRun({ 
          text: "Ingests traces from agent executions via OpenTelemetry. Enriches spans with security signals, detects anomalies, and stores time-series data in ClickHouse.",
          font: "Georgia"
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_3,
        children: [new TextRun({ text: "Policy Service", font: "Georgia" })]
      }),
      new Paragraph({
        spacing: { after: 150 },
        children: [new TextRun({ 
          text: "Evaluates security policies using OPA (Open Policy Agent). Supports tool access control, data flow restrictions, rate limiting, and human-in-the-loop approval gates.",
          font: "Georgia"
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_3,
        children: [new TextRun({ text: "Threat Modeling Service", font: "Georgia" })]
      }),
      new Paragraph({
        spacing: { after: 150 },
        children: [new TextRun({ 
          text: "Provides STRIDE-based threat analysis with MITRE ATLAS mapping for AI-specific attack techniques. Generates attack trees and mitigation recommendations.",
          font: "Georgia"
        })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_3,
        children: [new TextRun({ text: "Maturity Assessment Service", font: "Georgia" })]
      }),
      new Paragraph({
        spacing: { after: 150 },
        children: [new TextRun({ 
          text: "Implements a 5-level maturity model across governance, risk management, security controls, and operations domains. Generates improvement roadmaps.",
          font: "Georgia"
        })]
      }),
      
      // Technology Stack
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "3.2 Technology Stack", font: "Georgia" })]
      }),
      
      new Table({
        columnWidths: [2500, 3500, 4100],
        rows: [
          new TableRow({
            tableHeader: true,
            children: [
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 2500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Layer", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 3500, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Technology", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
              new TableCell({ borders: cellBorders, shading: headerShading, width: { size: 4100, type: WidthType.DXA },
                children: [new Paragraph({ alignment: AlignmentType.CENTER, children: [new TextRun({ text: "Rationale", bold: true, color: "FFFFFF", font: "Georgia", size: 20 })] })] }),
            ]
          }),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "API Server", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Go 1.22+ with Gin", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Performance, OPA native integration", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Portal UI", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "React / Next.js", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Enterprise UI patterns, SSR", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Policy Engine", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "OPA / Rego", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Industry standard, sub-ms evaluation", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Observability", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "OpenTelemetry + ClickHouse", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Vendor-neutral, high-volume time-series", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Persistence", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "PostgreSQL + Redis", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Reliable ACID, policy caching", font: "Georgia", size: 20 })] })] }),
          ]}),
          new TableRow({ children: [
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Infrastructure", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Kubernetes / Helm", font: "Georgia", size: 20 })] })] }),
            new TableCell({ borders: cellBorders, children: [new Paragraph({ children: [new TextRun({ text: "Cloud-agnostic, enterprise standard", font: "Georgia", size: 20 })] })] }),
          ]}),
        ]
      }),
      
      // Roadmap
      new Paragraph({ children: [new PageBreak()] }),
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: "4. Implementation Roadmap", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "Phase 1: Framework Foundation (Weeks 1-4)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Core Go API server scaffold with health/ready endpoints", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "PostgreSQL schema for controls, frameworks, crosswalks", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "NIST AI RMF control definitions (72 subcategories)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "NIST 800-53 crosswalk mappings", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "Phase 2: Observability (Weeks 5-8)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "OpenTelemetry trace ingestion", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Security signal enrichment pipeline", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "ClickHouse schema for time-series data", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Python/Go SDK middleware for LangChain", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "Phase 3: Policy Engine (Weeks 9-12)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "OPA integration with base Rego libraries", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Tool access control policies", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Data flow restriction policies", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Policy YAML → Rego compiler", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "Phase 4: Threat Modeling (Weeks 13-16)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "STRIDE threat templates for agentic systems", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "MITRE ATLAS technique mapping", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Attack tree generator", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Mitigation recommendation engine", font: "Georgia" })]
      }),
      
      new Paragraph({
        heading: HeadingLevel.HEADING_2,
        children: [new TextRun({ text: "Phase 5: Portal & Assessment (Weeks 17-20)", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "React governance portal", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Maturity self-assessment wizard", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Compliance report generation", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "GRC integration APIs (ServiceNow, Archer)", font: "Georgia" })]
      }),
      
      // Appendix
      new Paragraph({ children: [new PageBreak()] }),
      new Paragraph({
        heading: HeadingLevel.HEADING_1,
        children: [new TextRun({ text: "Appendix A: References", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "NIST AI Risk Management Framework: https://www.nist.gov/itl/ai-risk-management-framework", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "NIST SP 800-53 Rev 5: https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "ISO/IEC 42001:2023: https://www.iso.org/standard/81230.html", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "MITRE ATLAS: https://atlas.mitre.org", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Open Policy Agent: https://www.openpolicyagent.org", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Langfuse: https://langfuse.com", font: "Georgia" })]
      }),
      new Paragraph({
        numbering: { reference: "main-bullets", level: 0 },
        children: [new TextRun({ text: "Lakera Guard: https://www.lakera.ai", font: "Georgia" })]
      }),
    ]
  }]
});

// Save the document
Packer.toBuffer(doc).then(buffer => {
  fs.writeFileSync('/home/claude/agentguard/docs/AgentGuard-HLD.docx', buffer);
  console.log('Document created successfully: AgentGuard-HLD.docx');
});
