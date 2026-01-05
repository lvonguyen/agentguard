const pptxgen = require("pptxgenjs");
const { html2pptx } = require("./html2pptx");
const path = require("path");

async function generatePitchDeck() {
  const pptx = new pptxgen();
  
  // Set presentation properties
  pptx.layout = "LAYOUT_16x9";
  pptx.title = "AgentGuard - AI Security Governance Platform";
  pptx.author = "AgentGuard Inc.";
  pptx.company = "AgentGuard";
  pptx.subject = "Series A Investment Opportunity";
  
  // Define slide order
  const slides = [
    "slide01-title.html",
    "slide02-problem.html", 
    "slide03-solution.html",
    "slide04-market.html",
    "slide05-architecture.html",
    "slide06-differentiation.html",
    "slide07-business.html",
    "slide08-roadmap.html",
    "slide09-ask.html"
  ];
  
  // Process each slide
  for (const slideFile of slides) {
    const slidePath = path.join(__dirname, "slides", slideFile);
    console.log(`Processing: ${slideFile}`);
    await html2pptx(slidePath, pptx);
  }
  
  // Save presentation
  const outputPath = path.join(__dirname, "..", "docs", "AgentGuard-PitchDeck.pptx");
  await pptx.writeFile(outputPath);
  console.log(`\nPresentation saved to: ${outputPath}`);
}

generatePitchDeck().catch(err => {
  console.error("Error generating presentation:", err);
  process.exit(1);
});
