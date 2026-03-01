package model

import (
	"errors"
	"fmt"
	"time"
)

// PerformanceOptimizationConfig contains configuration for performance optimization features
type PerformanceOptimizationConfig struct {
	EnableConcurrentBenchmarks bool    `json:"enableConcurrentBenchmarks,omitempty"`
	MaxConcurrentWorkers       int     `json:"maxConcurrentWorkers,omitempty"`
	EnableProfileSampling      bool    `json:"enableProfileSampling,omitempty"`
	SamplingRate               float64 `json:"samplingRate,omitempty"`
	EnableMemoryOptimization   bool    `json:"enableMemoryOptimization,omitempty"`
	EnableEnhancedCaching     bool    `json:"enableEnhancedCaching,omitempty"`
	CacheMaxSizeMB            int     `json:"cacheMaxSizeMB,omitempty"`
	LargeCodebaseMode          bool    `json:"largeCodebaseMode,omitempty"`
}

// RemediationConfig contains configuration for automated remediation
type RemediationConfig struct {
	Enabled           bool    `json:"enabled"`
	MinConfidence     float64 `json:"minConfidence"`
	MaxCodeChanges    int     `json:"maxCodeChanges"`
	CodeChangeLimit   int     `json:"codeChangeLimit"`
	Provider          string  `json:"provider"`
	Model             string  `json:"model"`
	Temperature       float64 `json:"temperature"`
}

// ErrorContext provides structured error information for better error handling
type ErrorContext struct {
	ErrorType    string `json:"errorType,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
	Message      string `json:"message"`
	Details      string `json:"details,omitempty"`
	Suggestion   string `json:"suggestion,omitempty"`
	IsRecoverable bool   `json:"isRecoverable,omitempty"`
}

// Error types
const (
	ErrorTypeValidation   = "validation"
	ErrorTypeIO           = "io"
	ErrorTypeExecution    = "execution"
	ErrorTypeNetwork      = "network"
	ErrorTypeConfiguration = "configuration"
	ErrorTypeDependency   = "dependency"
)

// Error codes
const (
	ErrorCodeGitCloneFailed       = "git_clone_failed"
	ErrorCodeNoBenchmarksFound     = "no_benchmarks_found"
	ErrorCodeBenchmarkExecution    = "benchmark_execution_failed"
	ErrorCodeProfileCollection     = "profile_collection_failed"
	ErrorCodeFileOperation         = "file_operation_failed"
	ErrorCodeJSONParse             = "json_parse_failed"
	ErrorCodeDependencyMissing     = "dependency_missing"
	ErrorCodeInvalidConfiguration  = "invalid_configuration"
	ErrorCodeNetworkRequestFailed  = "network_request_failed"
	ErrorCodeTimeout               = "operation_timeout"
)

// NewErrorContext creates a new ErrorContext with the given parameters
func NewErrorContext(errorType, errorCode, message, details, suggestion string, isRecoverable bool) ErrorContext {
	return ErrorContext{
		ErrorType:    errorType,
		ErrorCode:    errorCode,
		Message:      message,
		Details:      details,
		Suggestion:   suggestion,
		IsRecoverable: isRecoverable,
	}
}

// Error implements the error interface
func (e ErrorContext) Error() string {
	if e.Details == "" && e.Suggestion == "" {
		return fmt.Sprintf("[%s:%s] %s", e.ErrorType, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("[%s:%s] %s\nDetails: %s\nSuggestion: %s", e.ErrorType, e.ErrorCode, e.Message, e.Details, e.Suggestion)
}

// Unwrap returns the underlying error if this wraps another error
func (e ErrorContext) Unwrap() error {
	if e.Details != "" {
		return errors.New(e.Details)
	}
	return nil
}

type Target struct {
	Type    string   `json:"type"`
	BaseURL string   `json:"baseUrl"`
	Command []string `json:"command,omitempty"`
}

type PluginInfo struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	SDKVersion   string       `json:"sdkVersion"`
	Capabilities Capabilities `json:"capabilities"`
}

type Capabilities struct {
	Targets  []string `json:"targets"`
	Profiles []string `json:"profiles"`
}

type Artifact struct {
	Kind        string `json:"kind"`
	ProfileType string `json:"profileType"`
	Path        string `json:"path"`
	ContentType string `json:"contentType"`
}

type ArtifactBundle struct {
	Metadata  Metadata   `json:"metadata"`
	Target    Target     `json:"target"`
	Artifacts []Artifact `json:"artifacts"`
}

type Metadata struct {
	Timestamp   time.Time `json:"timestamp"`
	DurationSec int       `json:"durationSec"`
	Service     string    `json:"service"`
	Scenario    string    `json:"scenario"`
	GitSha      string    `json:"gitSha"`
}

type ProfileBundle struct {
	Metadata  Metadata   `json:"metadata"`
	Target    Target     `json:"target"`
	Plugin    PluginRef  `json:"plugin"`
	Artifacts []Artifact `json:"artifacts"`
}

type PluginRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Finding represents a performance finding with deterministic analysis
type Finding struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	Category         string      `json:"category"` // cpu, alloc, heap, gc, mutex, block
	Severity         string      `json:"severity"` // low, medium, high, critical
	Confidence       float64     `json:"confidence"` // 0.0-1.0
	ImpactSummary    string      `json:"impactSummary"`
	Evidence         []EvidenceItem `json:"evidence"`
	DeterministicHints []string    `json:"deterministicHints"`
	Tags             []string    `json:"tags"`
	// Legacy fields for backward compatibility
	Score           int                     `json:"score,omitempty"`
	Top             []StackFrame            `json:"top,omitempty"`
	Callgraph       []CallgraphNode         `json:"callgraph,omitempty"`
	Regression      *RegressionAnalysis    `json:"regression,omitempty"`
	AllocationAnalysis *AllocationAnalysis `json:"allocationAnalysis,omitempty"`
	EvidenceLegacy  Evidence                `json:"evidenceLegacy,omitempty"`
}

// EvidenceItem represents a piece of evidence for a finding
type EvidenceItem struct {
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Value       string  `json:"value"`
	Weight      float64 `json:"weight"`
}

// AllocationAnalysis represents allocation-specific analysis results
type AllocationAnalysis struct {
	TotalAllocations float64            `json:"totalAllocations"`
	TopConcentration float64            `json:"topConcentration"`
	Severity         string             `json:"severity"`
	Score            int                `json:"score"`
	Hotspots         []AllocationHotspot `json:"hotspots"`
}

// AllocationHotspot represents a memory allocation hotspot
type AllocationHotspot struct {
	Function string  `json:"function"`
	File     string  `json:"file"`
	Line     int     `json:"line"`
	Count    float64 `json:"count"`
	Percent  float64 `json:"percent"`
}

type CallgraphNode struct {
	Function string          `json:"function"`
	File     string          `json:"file"`
	Line     int             `json:"line"`
	Depth    int             `json:"depth"`
	Cum      float64         `json:"cum"`
	Flat     float64         `json:"flat"`
	Children []CallgraphNode `json:"children,omitempty"`
}

type RegressionAnalysis struct {
	BaselineScore int     `json:"baseline_score"`
	CurrentScore  int     `json:"current_score"`
	Delta         int     `json:"delta"`
	Percentage    float64 `json:"percentage"`
	Severity      string  `json:"severity"`
	Confidence    int     `json:"confidence"`
	BaselineRef   string  `json:"baseline_ref,omitempty"`
	CurrentRef    string  `json:"current_ref,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
}

// BaselineComparison represents a performance baseline for comparative analysis
type BaselineComparison struct {
	BaselinePath string `json:"baselinePath"`
	CurrentPath  string `json:"currentPath"`
	BaselineRef   string `json:"baselineRef,omitempty"`
	CurrentRef    string `json:"currentRef,omitempty"`
	Threshold     float64 `json:"threshold,omitempty"` // Percentage threshold for regression detection (default: 10%)
}

// PerformanceTrend represents performance changes over time
type PerformanceTrend struct {
	Metric       string    `json:"metric"`
	Baseline     float64   `json:"baseline"`
	Current      float64   `json:"current"`
	Delta        float64   `json:"delta"`
	Percentage   float64   `json:"percentage"`
	Severity     string    `json:"severity"`
	Confidence   int       `json:"confidence"`
	Timestamps   []time.Time `json:"timestamps,omitempty"`
}

type StackFrame struct {
	Function string  `json:"function"`
	File     string  `json:"file"`
	Line     int     `json:"line"`
	Cum      float64 `json:"cum"`
	Flat     float64 `json:"flat"`
}

type Evidence struct {
	ArtifactPath string    `json:"artifactPath"`
	ProfileType  string    `json:"profileType"`
	ExtractedAt  time.Time `json:"extractedAt"`
}

type FindingsBundle struct {
	Summary  Summary   `json:"summary"`
	Findings []Finding `json:"findings"`
}

type Summary struct {
	TopIssueTags []string `json:"topIssueTags"`
	OverallScore int      `json:"overallScore"`
	Notes        []string `json:"notes"`
}

type CollectRequest struct {
	Target      Target            `json:"target"`
	DurationSec int               `json:"durationSec"`
	Profiles    []string          `json:"profiles"`
	OutDir      string            `json:"outDir"`
	Metadata    map[string]string `json:"metadata"`
}
