package report

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

type Reporter struct {
}

func NewReporter() *Reporter {
	return &Reporter{}
}

// Generate creates a markdown report from findings (backward compatible)
func (r *Reporter) Generate(findings model.FindingsBundle) (string, error) {
	return r.GenerateWithInsights(findings, nil)
}

// GenerateWithInsights creates a markdown report with optional LLM insights
func (r *Reporter) GenerateWithInsights(findings model.FindingsBundle, insights *model.InsightsBundle) (string, error) {
	var sb strings.Builder

	// Header
	sb.WriteString("# Performance Triage Report\n\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	// Executive Summary
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString(fmt.Sprintf("- **Overall Score**: %d/100\n", findings.Summary.OverallScore))
	sb.WriteString(fmt.Sprintf("- **Top Issues**: %s\n", strings.Join(findings.Summary.TopIssueTags, ", ")))
	if len(findings.Summary.Notes) > 0 {
		sb.WriteString("- **Notes**:\n")
		for _, note := range findings.Summary.Notes {
			sb.WriteString(fmt.Sprintf("  - %s\n", note))
		}
	}

	// Add LLM insights to executive summary if available
	if insights != nil && insights.ExecutiveSummary.Overview != "" {
		sb.WriteString("\n### LLM Insights\n\n")
		sb.WriteString(fmt.Sprintf("**Overview**: %s\n", insights.ExecutiveSummary.Overview))
		sb.WriteString(fmt.Sprintf("**Overall Severity**: %s (Confidence: %d%%)\n", 
			insights.ExecutiveSummary.OverallSeverity, insights.ExecutiveSummary.Confidence))
		if len(insights.ExecutiveSummary.KeyThemes) > 0 {
			sb.WriteString("**Key Themes**: ")
			sb.WriteString(strings.Join(insights.ExecutiveSummary.KeyThemes, ", "))
			sb.WriteString("\n")
		}
		if insights.DisabledReason != "" {
			sb.WriteString(fmt.Sprintf("*LLM Status*: %s\n", insights.DisabledReason))
		}
	}

	sb.WriteString("\n")

	// Findings by category
	for _, finding := range findings.Findings {
		sb.WriteString(fmt.Sprintf("## %s: %s\n\n", strings.Title(finding.Category), finding.Title))
		sb.WriteString(fmt.Sprintf("- **Severity**: %s\n", strings.Title(finding.Severity)))
		sb.WriteString(fmt.Sprintf("- **Score**: %d\n", finding.Score))
		sb.WriteString("- **Evidence**:\n")
		sb.WriteString(fmt.Sprintf("  - Profile: %s\n", finding.Evidence.ProfileType))
		sb.WriteString(fmt.Sprintf("  - Artifact: %s\n", finding.Evidence.ArtifactPath))
		sb.WriteString("\n")

		if len(finding.Top) > 0 {
			sb.WriteString("### Top Hotspots\n\n")
			sb.WriteString("| Function | File | Line | Cumulative | Flat |\n")
			sb.WriteString("|----------|------|------|------------|------|\n")
			for _, frame := range finding.Top {
				sb.WriteString(fmt.Sprintf("| %s | %s | %d | %.2f | %.2f |\n",
					frame.Function, frame.File, frame.Line, frame.Cum, frame.Flat))
			}
			sb.WriteString("\n")
		}

		// Add callgraph visualization if available
		if len(finding.Callgraph) > 0 {
			sb.WriteString("### Callgraph Analysis (Depth 3)\n\n")
			sb.WriteString("```\n")
			for _, node := range finding.Callgraph {
				renderCallgraphNode(&node, 0, &sb)
			}
			sb.WriteString("```\n\n")
		}

		// Add regression analysis if available
		if finding.Regression != nil {
			sb.WriteString("### Regression Analysis\n\n")
			sb.WriteString(fmt.Sprintf("- **Baseline Score**: %d\n", finding.Regression.BaselineScore))
			sb.WriteString(fmt.Sprintf("- **Current Score**: %d\n", finding.Regression.CurrentScore))
			sb.WriteString(fmt.Sprintf("- **Delta**: %d (%.1f%%)\n", finding.Regression.Delta, finding.Regression.Percentage))
			sb.WriteString(fmt.Sprintf("- **Severity**: %s\n", strings.Title(finding.Regression.Severity)))
			sb.WriteString(fmt.Sprintf("- **Confidence**: %d%%\n", finding.Regression.Confidence))
			sb.WriteString("\n")

			if finding.Regression.Severity == "improved" {
				sb.WriteString("📈 **Performance Improvement Detected**\n")
				sb.WriteString("This profile shows significant improvement over the baseline.\n")
			} else if finding.Regression.Severity != "none" && finding.Regression.Severity != "low" {
				sb.WriteString("⚠️ **Potential Regression Detected**\n")
				sb.WriteString(fmt.Sprintf("This profile shows %s regression over the baseline.\n", finding.Regression.Severity))
			}
			sb.WriteString("\n")
		}

		// Add LLM insights for this finding if available
		if insights != nil && len(insights.PerFinding) > 0 {
			for _, insight := range insights.PerFinding {
				if insight.FindingID == finding.Category {
					sb.WriteString("### LLM Insights\n\n")
					sb.WriteString(fmt.Sprintf("**Narrative**: %s\n\n", insight.Narrative))
					if len(insight.LikelyRootCauses) > 0 {
						sb.WriteString("**Likely Root Causes**:\n")
						for _, cause := range insight.LikelyRootCauses {
							sb.WriteString(fmt.Sprintf("  - %s\n", cause))
						}
						sb.WriteString("\n")
					}
					if len(insight.Suggestions) > 0 {
						sb.WriteString("**Suggestions**:\n")
						for _, suggestion := range insight.Suggestions {
							sb.WriteString(fmt.Sprintf("  - %s\n", suggestion))
						}
						sb.WriteString("\n")
					}
					if len(insight.NextMeasurements) > 0 {
						sb.WriteString("**Next Measurements**:\n")
						for _, measurement := range insight.NextMeasurements {
							sb.WriteString(fmt.Sprintf("  - %s\n", measurement))
						}
						sb.WriteString("\n")
					}
					if len(insight.Caveats) > 0 {
						sb.WriteString("**Caveats**:\n")
						for _, caveat := range insight.Caveats {
							sb.WriteString(fmt.Sprintf("  - %s\n", caveat))
						}
						sb.WriteString("\n")
					}
					sb.WriteString(fmt.Sprintf("*Confidence: %d%%*\n\n", insight.Confidence))
					break
				}
			}
		}
	}

	// Footer
	sb.WriteString("---\n\n")
	sb.WriteString("*Generated by triageprof*\n")

	return sb.String(), nil
}

// GenerateJSON creates a structured JSON report
func (r *Reporter) GenerateJSON(findings model.FindingsBundle, insights *model.InsightsBundle, options model.JSONReportOptions) ([]byte, error) {
	report := model.JSONReport{
		SchemaVersion: "1.0",
		GeneratedAt:   time.Now(),
		Summary: model.ReportSummary{
			OverallScore: findings.Summary.OverallScore,
			TopIssueTags: findings.Summary.TopIssueTags,
			Severity:     determineSeverity(findings.Summary.OverallScore),
		},
		Findings: make([]model.ReportFinding, len(findings.Findings)),
		Insights: insights,
	}

	// Convert findings to report format
	for i, finding := range findings.Findings {
		report.Findings[i] = model.ReportFinding{
			ID:          fmt.Sprintf("finding-%d", i),
			Category:    finding.Category,
			Title:       finding.Title,
			Severity:    finding.Severity,
			Score:       finding.Score,
			TopHotspots: finding.Top,
			Evidence:    finding.Evidence,
		}
	}

	// Apply options
	if options.PrettyPrint {
		return json.MarshalIndent(report, "", "  ")
	}
	return json.Marshal(report)
}

// renderCallgraphNode recursively renders a callgraph node as ASCII tree
func renderCallgraphNode(node *model.CallgraphNode, indent int, sb *strings.Builder) {
	indentStr := strings.Repeat("  ", indent)
	sb.WriteString(fmt.Sprintf("%s%s (%.1f%% cum, %.1f%% flat)\n",
		indentStr, node.Function, node.Cum, node.Flat))

	for _, child := range node.Children {
		renderCallgraphNode(&child, indent+1, sb)
	}
}

func determineSeverity(score int) string {
	switch {
	case score >= 80:
		return "critical"
	case score >= 60:
		return "high"
	case score >= 40:
		return "medium"
	case score >= 20:
		return "low"
	default:
		return "info"
	}
}