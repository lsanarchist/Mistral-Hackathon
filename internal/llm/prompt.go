package llm

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// MaxPromptSize is the maximum size of the prompt in characters
const MaxPromptSize = 12000

// Redaction patterns for sensitive data
var (
	hostnamePattern = regexp.MustCompile(`\b(?:localhost|[a-zA-Z0-9\-]+\.[a-zA-Z0-9\-]+\.[a-zA-Z]{2,})\b`)
	tokenPattern    = regexp.MustCompile(`\b(?:[A-Za-z0-9]{32,}|[A-Za-z0-9_]{20,})\b`)
	pathPattern     = regexp.MustCompile(`[A-Za-z]:\\|\/[^\/]+\/[^\/]+`)
)

// BuildPrompt creates a structured prompt from bundle and findings with redaction
type PromptBuilder struct {
	Bundle   *model.ProfileBundle
	Findings *model.FindingsBundle
	MaxSize  int
}

func NewPromptBuilder(bundle *model.ProfileBundle, findings *model.FindingsBundle) *PromptBuilder {
	return &PromptBuilder{
		Bundle:   bundle,
		Findings: findings,
		MaxSize:  MaxPromptSize,
	}
}

// Build creates the final prompt with redaction and size limiting
func (p *PromptBuilder) Build() (string, error) {
	// Build structured data
	metadata := p.buildMetadata()
	findingsSummary := p.buildFindingsSummary()

	// Create prompt template
	prompt := fmt.Sprintf(`Analyze the following performance findings and provide insights in JSON format matching the schema:

METADATA:
%s

FINDINGS SUMMARY:
%s

INSTRUCTIONS:
1. Generate an executive summary with overall severity assessment
2. Identify top 3 risks with impact/likelihood
3. Suggest top 3 actions with priority/effort estimates
4. For each finding, provide narrative, likely root causes, suggestions, and next measurements
5. Include caveats about limitations and confidence levels
6. Output JSON ONLY, no markdown or explanations`,
		metadata, findingsSummary)

	// Check size limit
	if len(prompt) > p.MaxSize {
		return "", fmt.Errorf("prompt size %d exceeds maximum %d characters", len(prompt), p.MaxSize)
	}

	return prompt, nil
}

// buildMetadata creates a redacted metadata section
func (p *PromptBuilder) buildMetadata() string {
	if p.Bundle == nil {
		return "No metadata available"
	}

	// Redact sensitive information from target
	targetURL := p.redactSensitiveInfo(p.Bundle.Target.BaseURL)

	// Build artifact summary
	var artifacts []string
	for _, artifact := range p.Bundle.Artifacts {
		artifacts = append(artifacts, fmt.Sprintf("- %s (%s)",
			p.redactPath(artifact.ProfileType),
			p.redactPath(artifact.Kind)))
	}

	return fmt.Sprintf(`Service: %s
Scenario: %s
Duration: %d seconds
Target: %s
Profiles Collected: %s
Artifacts: %d total
Timestamp: %s`,
		p.redactSensitiveInfo(p.Bundle.Metadata.Service),
		p.redactSensitiveInfo(p.Bundle.Metadata.Scenario),
		p.Bundle.Metadata.DurationSec,
		targetURL,
		strings.Join(artifacts, ", "),
		len(p.Bundle.Artifacts),
		p.Bundle.Metadata.Timestamp.Format("2006-01-02 15:04:05"))
}

// buildFindingsSummary creates a redacted findings summary
func (p *PromptBuilder) buildFindingsSummary() string {
	if p.Findings == nil || len(p.Findings.Findings) == 0 {
		return "No findings available"
	}

	var findingsSummary strings.Builder
	findingsSummary.WriteString(fmt.Sprintf("Overall Score: %d/100\n", p.Findings.Summary.OverallScore))
	findingsSummary.WriteString(fmt.Sprintf("Top Issue Tags: %s\n\n", strings.Join(p.Findings.Summary.TopIssueTags, ", ")))

	// Limit to top 5 findings to control prompt size
	findings := p.Findings.Findings
	if len(findings) > 5 {
		findings = findings[:5]
	}

	for i, finding := range findings {
		findingsSummary.WriteString(fmt.Sprintf("\nFINDING %d: %s\n", i+1, finding.Title))
		findingsSummary.WriteString(fmt.Sprintf("Category: %s\n", finding.Category))
		findingsSummary.WriteString(fmt.Sprintf("Severity: %s\n", finding.Severity))
		findingsSummary.WriteString(fmt.Sprintf("Score: %d\n", finding.Score))
		findingsSummary.WriteString("Top Stack Frames:\n")

		// Limit to top 10 frames and redact sensitive info
		frames := finding.Top
		if len(frames) > 10 {
			frames = frames[:10]
		}

		for _, frame := range frames {
			findingsSummary.WriteString(fmt.Sprintf("  - %s (%.1fs cum, %.1fs flat)\n",
				p.redactStackFrame(frame),
				frame.Cum,
				frame.Flat))
		}

		findingsSummary.WriteString(fmt.Sprintf("\nEvidence: %s profile from %s\n",
			p.redactPath(finding.Evidence.ProfileType),
			p.redactPath(filepath.Base(finding.Evidence.ArtifactPath))))
	}

	return findingsSummary.String()
}

// redactSensitiveInfo removes sensitive data from strings
func (p *PromptBuilder) redactSensitiveInfo(text string) string {
	// Redact hostnames
	text = hostnamePattern.ReplaceAllString(text, "[REDACTED_HOSTNAME]")
	// Redact long tokens
	text = tokenPattern.ReplaceAllString(text, "[REDACTED_TOKEN]")
	// Redact paths
	text = pathPattern.ReplaceAllString(text, "[REDACTED_PATH]")
	// Truncate long strings
	if len(text) > 200 {
		text = text[:200] + "..."
	}
	return text
}

// redactPath removes sensitive path information
func (p *PromptBuilder) redactPath(path string) string {
	// Remove absolute paths, keep only filename/extension
	base := filepath.Base(path)
	if base == "." || base == "/" {
		return "[REDACTED_PATH]"
	}
	return base
}

// redactStackFrame redacts sensitive information from stack frames
func (p *PromptBuilder) redactStackFrame(frame model.StackFrame) string {
	function := p.redactSensitiveInfo(frame.Function)
	_ = p.redactSensitiveInfo(frame.File)

	// Limit line number to reasonable range
	line := frame.Line
	if line < 0 {
		line = 0
	} else if line > 100000 {
		line = 100000
	}

	return fmt.Sprintf("%s (line %d)", function, line)
}
