package model

import "time"

// InsightsBundle contains LLM-generated insights about performance findings
type InsightsBundle struct {
    SchemaVersion    string            `json:"schema_version"`
    GeneratedAt      time.Time         `json:"generated_at"`
    Model            string            `json:"model,omitempty"`
    RequestID        string            `json:"request_id,omitempty"`
    DisabledReason   string            `json:"disabled_reason,omitempty"`
    ExecutiveSummary ExecutiveSummary `json:"executive_summary"`
    TopRisks        []RiskItem        `json:"top_risks,omitempty"`
    TopActions      []ActionItem      `json:"top_actions,omitempty"`
    PerformanceCategories map[string]int `json:"performance_categories,omitempty"`
    PerFinding      []FindingInsight   `json:"per_finding,omitempty"`
    ROIAnalysis      []ROIItem         `json:"roi_analysis,omitempty"`
    TechnicalDeepDive TechnicalAnalysis `json:"technical_deep_dive,omitempty"`
}

// ExecutiveSummary provides high-level overview
type ExecutiveSummary struct {
    Overview        string   `json:"overview"`
    OverallSeverity string   `json:"overall_severity"`
    KeyThemes       []string `json:"key_themes,omitempty"`
    Confidence      int      `json:"confidence"` // 0-100
    PerformanceScore int     `json:"performance_score,omitempty"` // 0-100
    ImprovementPotential int `json:"improvement_potential,omitempty"` // percentage
}

// RiskItem represents a significant risk identified by LLM
type RiskItem struct {
    Description string `json:"description"`
    Severity    string `json:"severity"`
    Impact      string `json:"impact"`
    Likelihood  string `json:"likelihood"`
    AffectedComponents []string `json:"affected_components,omitempty"`
    PotentialImpact string `json:"potential_impact,omitempty"` // quantitative estimate
}

// ActionItem represents a recommended action
type ActionItem struct {
    Description    string   `json:"description"`
    Priority       string   `json:"priority"`
    EstimatedEffort string   `json:"estimated_effort"`
    Categories     []string `json:"categories,omitempty"`
    ImplementationComplexity string `json:"implementation_complexity,omitempty"` // Low/Medium/High
    ExpectedImpact string   `json:"expected_impact,omitempty"` // quantitative estimate
    CodeExamples    []string `json:"code_examples,omitempty"`
    ValidationMetrics []string `json:"validation_metrics,omitempty"`
}

// FindingInsight provides per-finding analysis
type FindingInsight struct {
    FindingID        string   `json:"finding_id"`
    Narrative        string   `json:"narrative"`
    LikelyRootCauses []string `json:"likely_root_causes,omitempty"`
    Suggestions      []string `json:"suggestions,omitempty"`
    NextMeasurements  []string `json:"next_measurements,omitempty"`
    Caveats          []string `json:"caveats,omitempty"`
    Confidence       int      `json:"confidence"` // 0-100
    PerformanceImpact string   `json:"performance_impact,omitempty"` // quantitative estimate
    ImplementationComplexity string `json:"implementation_complexity,omitempty"` // Low/Medium/High
    CodeExamples    []string `json:"code_examples,omitempty"`
    BeforeAfterMetrics []string `json:"before_after_metrics,omitempty"`
}

// ROIItem represents return on investment analysis
type ROIItem struct {
    ActionID        string  `json:"action_id"`
    Description     string  `json:"description"`
    EstimatedEffort string  `json:"estimated_effort"`
    ExpectedImpact  string  `json:"expected_impact"`
    CostBenefitRatio float64 `json:"cost_benefit_ratio"`
    PriorityScore   int     `json:"priority_score"` // 0-100
}

// TechnicalAnalysis provides advanced technical insights
type TechnicalAnalysis struct {
    MemoryPatterns        []MemoryPattern        `json:"memory_patterns,omitempty"`
    CPUUtilization        []CPUUtilization       `json:"cpu_utilization,omitempty"`
    BlockingOperations    []BlockingOperation   `json:"blocking_operations,omitempty"`
    ConcurrencyPatterns   []ConcurrencyPattern  `json:"concurrency_patterns,omitempty"`
    AlgorithmAnalysis     []AlgorithmAnalysis    `json:"algorithm_analysis,omitempty"`
    CacheEfficiency       []CacheEfficiency     `json:"cache_efficiency,omitempty"`
    IOPatterns            []IOPattern           `json:"io_patterns,omitempty"`
    GarbageCollection     []GCAnalysis          `json:"garbage_collection,omitempty"`
}

// MemoryPattern represents memory allocation and usage patterns
type MemoryPattern struct {
    PatternType     string  `json:"pattern_type"`
    Description     string  `json:"description"`
    CurrentUsage    string  `json:"current_usage"`
    Optimization    string  `json:"optimization"`
    ExpectedSavings string  `json:"expected_savings"`
    Implementation  string  `json:"implementation"`
}

// CPUUtilization represents CPU usage analysis
type CPUUtilization struct {
    Component      string  `json:"component"`
    CurrentUsage   float64 `json:"current_usage"`
    HotspotAnalysis string  `json:"hotspot_analysis"`
    Optimization   string  `json:"optimization"`
    ExpectedGain   float64 `json:"expected_gain"`
}

// BlockingOperation represents synchronization bottleneck analysis
type BlockingOperation struct {
    OperationType  string  `json:"operation_type"`
    Location       string  `json:"location"`
    CurrentLatency string  `json:"current_latency"`
    RootCause      string  `json:"root_cause"`
    Solution       string  `json:"solution"`
    ExpectedGain   string  `json:"expected_gain"`
}

// ConcurrencyPattern represents concurrency and parallelism analysis
type ConcurrencyPattern struct {
    PatternType    string  `json:"pattern_type"`
    CurrentUsage    string  `json:"current_usage"`
    Bottleneck     string  `json:"bottleneck"`
    Optimization   string  `json:"optimization"`
    ExpectedGain   string  `json:"expected_gain"`
}

// AlgorithmAnalysis represents algorithm complexity analysis
type AlgorithmAnalysis struct {
    Algorithm      string  `json:"algorithm"`
    CurrentComplexity string `json:"current_complexity"`
    Analysis        string  `json:"analysis"`
    Optimization    string  `json:"optimization"`
    ExpectedGain    string  `json:"expected_gain"`
}

// CacheEfficiency represents cache usage analysis
type CacheEfficiency struct {
    CacheType      string  `json:"cache_type"`
    CurrentHitRate float64 `json:"current_hit_rate"`
    Analysis        string  `json:"analysis"`
    Optimization   string  `json:"optimization"`
    ExpectedGain   float64 `json:"expected_gain"`
}

// IOPattern represents I/O operation analysis
type IOPattern struct {
    OperationType  string  `json:"operation_type"`
    CurrentPattern string  `json:"current_pattern"`
    Analysis       string  `json:"analysis"`
    Optimization   string  `json:"optimization"`
    ExpectedGain   string  `json:"expected_gain"`
}

// GCAnalysis represents garbage collection analysis
type GCAnalysis struct {
    GCType         string  `json:"gc_type"`
    CurrentPressure string `json:"current_pressure"`
    Analysis        string  `json:"analysis"`
    Optimization   string  `json:"optimization"`
    ExpectedGain   string  `json:"expected_gain"`
}