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
		sb.WriteString("\n### 🤖 LLM Insights\n\n")
		sb.WriteString(fmt.Sprintf("**Overview**: %s\n", insights.ExecutiveSummary.Overview))
		sb.WriteString(fmt.Sprintf("**Overall Severity**: %s (Confidence: %d%%)",
			insights.ExecutiveSummary.OverallSeverity, insights.ExecutiveSummary.Confidence))
		if len(insights.ExecutiveSummary.KeyThemes) > 0 {
			sb.WriteString("\n**Key Themes**: ")
			sb.WriteString(strings.Join(insights.ExecutiveSummary.KeyThemes, ", "))
		}
		if insights.DisabledReason != "" {
			sb.WriteString(fmt.Sprintf("\n*LLM Status*: %s", insights.DisabledReason))
		}
		
		// Add performance categories if available
		if len(insights.PerformanceCategories) > 0 {
			sb.WriteString("\n\n#### 📊 Performance Categories\n\n")
			for category, count := range insights.PerformanceCategories {
				sb.WriteString(fmt.Sprintf("   - **%s**: %d findings\n", category, count))
			}
		}
		
		// Add top risks if available
		if len(insights.TopRisks) > 0 {
			sb.WriteString("\n\n#### 🚨 Top Risks\n\n")
			for i, risk := range insights.TopRisks {
				if i >= 3 {
					break
				}
				sb.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, risk.Description))
				sb.WriteString(fmt.Sprintf("   - Severity: %s\n", risk.Severity))
				sb.WriteString(fmt.Sprintf("   - Impact: %s\n", risk.Impact))
				sb.WriteString(fmt.Sprintf("   - Likelihood: %s\n", risk.Likelihood))
			}
		}
		
		// Add top actions if available
		if len(insights.TopActions) > 0 {
			sb.WriteString("\n#### 🎯 Top Action Items\n\n")
			for i, action := range insights.TopActions {
				if i >= 3 {
					break
				}
				sb.WriteString(fmt.Sprintf("**%d. %s**\n", i+1, action.Description))
				sb.WriteString(fmt.Sprintf("   - Priority: %s\n", action.Priority))
				sb.WriteString(fmt.Sprintf("   - Estimated Effort: %s\n", action.EstimatedEffort))
				if len(action.Categories) > 0 {
					sb.WriteString(fmt.Sprintf("   - Categories: %s\n", strings.Join(action.Categories, ", ")))
				}
			}
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
			sb.WriteString("### Callgraph Analysis\n\n")
			sb.WriteString("```\n")
			for _, node := range finding.Callgraph {
				renderCallgraphNode(&node, 0, &sb)
			}
			sb.WriteString("```\n\n")
			
			// Add callgraph statistics
			totalNodes := countCallgraphNodes(finding.Callgraph)
			maxDepth := findMaxCallgraphDepth(finding.Callgraph)
			sb.WriteString(fmt.Sprintf("**Callgraph Statistics**: %d nodes, max depth %d\n\n", totalNodes, maxDepth))
		}

		// Add allocation analysis if available
		if finding.AllocationAnalysis != nil {
			sb.WriteString("### Allocation Analysis\n\n")
			sb.WriteString(fmt.Sprintf("- **Total Allocations**: %.0f\n", finding.AllocationAnalysis.TotalAllocations))
			sb.WriteString(fmt.Sprintf("- **Top 10%% Concentration**: %.1f%%\n", finding.AllocationAnalysis.TopConcentration*100))
			sb.WriteString(fmt.Sprintf("- **Allocation Severity**: %s\n", strings.Title(finding.AllocationAnalysis.Severity)))
			sb.WriteString(fmt.Sprintf("- **Allocation Score**: %d/100\n", finding.AllocationAnalysis.Score))
			sb.WriteString("\n")

			if finding.AllocationAnalysis.TopConcentration > 0.5 {
				sb.WriteString("⚠️ **High Allocation Concentration Detected**\n")
				sb.WriteString(fmt.Sprintf("Top functions account for %.1f%% of all allocations.\n", finding.AllocationAnalysis.TopConcentration*100))
				sb.WriteString("This indicates potential memory allocation hotspots that may benefit from optimization.\n")
			} else {
				sb.WriteString("✅ **Balanced Allocation Pattern**\n")
				sb.WriteString("Allocations are reasonably distributed across functions.\n")
			}
			sb.WriteString("\n")

			// Add allocation hotspots table
			if len(finding.AllocationAnalysis.Hotspots) > 0 {
				sb.WriteString("#### Top Allocation Hotspots\n\n")
				sb.WriteString("| Function | File | Line | Count | Percentage |\n")
				sb.WriteString("|----------|------|------|-------|------------|\n")
				for _, hotspot := range finding.AllocationAnalysis.Hotspots {
					sb.WriteString(fmt.Sprintf("| %s | %s | %d | %.0f | %.1f%% |\n",
						hotspot.Function, hotspot.File, hotspot.Line, hotspot.Count, hotspot.Percent))
				}
				sb.WriteString("\n")
			}
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
					sb.WriteString("### 🤖 LLM Insights\n\n")
					sb.WriteString(fmt.Sprintf("**Narrative**: %s\n\n", insight.Narrative))
					
					// Add root causes with emojis
					if len(insight.LikelyRootCauses) > 0 {
						sb.WriteString("**🔍 Likely Root Causes**:\n")
						for i, cause := range insight.LikelyRootCauses {
							sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, cause))
						}
						sb.WriteString("\n")
					}
					
					// Add suggestions with emojis
					if len(insight.Suggestions) > 0 {
						sb.WriteString("**💡 Suggestions**:\n")
						for i, suggestion := range insight.Suggestions {
							sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, suggestion))
						}
						sb.WriteString("\n")
					}
					
					// Add next measurements
					if len(insight.NextMeasurements) > 0 {
						sb.WriteString("**📊 Next Measurements**:\n")
						for i, measurement := range insight.NextMeasurements {
							sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, measurement))
						}
						sb.WriteString("\n")
					}
					
					// Add caveats with warning emoji
					if len(insight.Caveats) > 0 {
						sb.WriteString("**⚠️  Caveats**:\n")
						for i, caveat := range insight.Caveats {
							sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, caveat))
						}
						sb.WriteString("\n")
					}
					
					// Add confidence with appropriate emoji
					confidenceEmoji := "🟡"
					if insight.Confidence >= 80 {
						confidenceEmoji = "🟢"
					} else if insight.Confidence <= 50 {
						confidenceEmoji = "🔴"
					}
					sb.WriteString(fmt.Sprintf("**Confidence**: %s %d%%\n\n", confidenceEmoji, insight.Confidence))
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
			ID:               fmt.Sprintf("finding-%d", i),
			Category:         finding.Category,
			Title:            finding.Title,
			Severity:         finding.Severity,
			Score:            finding.Score,
			TopHotspots:      finding.Top,
			Callgraph:        finding.Callgraph,
			Regression:       finding.Regression,
			AllocationAnalysis: finding.AllocationAnalysis,
			Evidence:         finding.Evidence,
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
	
	// Use tree characters for better visualization
	treePrefix := "├── "
	if indent == 0 {
		treePrefix = ""
	}
	
	sb.WriteString(fmt.Sprintf("%s%s%s (cum: %.1f, flat: %.1f, depth: %d)\n",
		indentStr, treePrefix, node.Function, node.Cum, node.Flat, node.Depth))

	for _, child := range node.Children {
		renderCallgraphNode(&child, indent+1, sb)
	}
}

// countCallgraphNodes counts total nodes in callgraph
func countCallgraphNodes(nodes []model.CallgraphNode) int {
	count := 0
	for _, node := range nodes {
		count += countCallgraphNode(&node)
	}
	return count
}

// countCallgraphNode recursively counts nodes
func countCallgraphNode(node *model.CallgraphNode) int {
	count := 1
	for _, child := range node.Children {
		count += countCallgraphNode(&child)
	}
	return count
}

// findMaxCallgraphDepth finds maximum depth in callgraph
func findMaxCallgraphDepth(nodes []model.CallgraphNode) int {
	maxDepth := 0
	for _, node := range nodes {
		depth := findMaxCallgraphNodeDepth(&node)
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	return maxDepth
}

// findMaxCallgraphNodeDepth recursively finds max depth
func findMaxCallgraphNodeDepth(node *model.CallgraphNode) int {
	maxDepth := node.Depth
	for _, child := range node.Children {
		childDepth := findMaxCallgraphNodeDepth(&child)
		if childDepth > maxDepth {
			maxDepth = childDepth
		}
	}
	return maxDepth
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
