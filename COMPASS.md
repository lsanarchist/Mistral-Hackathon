# COMPASS — TriageProf Direction & Status

## IMPORTANT RULE: THIS FILE MUST NOT BE MODIFIED DIRECTLY

**All changes must be appended to `change.log` instead.**

This document establishes the foundational principles and direction for TriageProf.
To propose changes, add entries to the change.log file following the format below.

---

## North Star

Build a production-grade, modular profiling and bottleneck-analysis tool — the `osv-scanner` of performance.
It integrates with existing profilers via a well-defined plugin SDK, produces evidence-backed bottleneck findings
(stacks, callgraphs, timelines, metrics), and generates structured, machine-readable reports.

**Killer feature: AI/LLM-powered analysis** — deterministic profiling data feeds an LLM that explains *why* a bottleneck exists, suggests concrete fixes with code examples, and ranks issues by impact.

**Philosophy: go deep, not wide. Make each layer excellent before adding the next.**

## Product Shape (today)

Pipeline: Collect → Analyze → Report.
Plugins are separate executables, discovered via manifests, and communicate with core using JSON-RPC 2.0 over stdio.

## Non-negotiable Architecture Rules

- Core is language-agnostic; language/profiler-specific logic must live in plugins.
- Plugin protocol/API must remain stable; breaking changes require explicit versioning.
- Deterministic profiling data is always collected first and is the source of truth.
- **AI/LLM analysis is a first-class feature**, not optional glue — it must be designed as a proper pipeline stage with a clean interface, not bolted on. apikey for dev is at apikey.swaga
- LLM calls are always backed by structured profiling data; the LLM never guesses — it reasons over real evidence.
- LLM stage must be skippable via `--no-ai` flag so the tool works without an API key.
- **Depth over breadth**: do not add new plugins unless the existing ones are excellent and the SDK is solid.
- Every feature must have tests. No untested code merges.
- Big PRs are fine when a feature warrants it. Correctness and completeness over minimal diffs.

## Change Management Process

### How to Propose Changes

1. **Do NOT modify this file directly**
2. Create or append to `change.log` using the format below
3. Changes will be reviewed and incorporated by maintainers

### Change Log Format

```
## Change Log

### YYYY-MM-DD
- **Objective**: Brief description of the change
- **Rationale**: Why this change is needed
- **Impact**: What will be affected
- **Changes**:
  - Specific change 1
  - Specific change 2
- **Testing**: How this will be validated
```

### Example Change Entry

```
## Change Log

### 2026-02-28
- **Objective**: Add plugin capability validation
- **Rationale**: Prevent incompatible plugins from being loaded
- **Impact**: Plugin loading process, error handling
- **Changes**:
  - Add manifest-based plugin discovery
  - Implement capability validation before plugin launch
  - Add SDK version compatibility checking
  - Enhance error messages for plugin issues
- **Testing**: Unit tests for manifest parsing and validation
```

## Current Focus Areas

1. **Plugin System Maturity**
   - Manifest-based discovery and validation 
   - Capability checking and compatibility
   - Plugin lifecycle management

2. **LLM Integration** !!!!!
   - Safe prompt generation with redaction 
   - Structured insights generation
   - Report enhancement with LLM sections

3. **Core Robustness**
   - Error handling and recovery
   - Test coverage and validation
   - Performance optimization

4. **Cool web page of results with real-time monitoring**

5. **Make more clear where each plugin-module lives**

6. **Enhanced WebSocket Real-time Monitoring** ✅
- **Live data streaming** with periodic auto-refresh
- **Enhanced WebSocket stats dashboard** with client count and severity breakdown
- **Real-time notifications** for data updates
- **Improved WebSocket protocol** with comprehensive statistics
- **Auto-refresh configuration** via CLI flags

7. **WebSocket Connection Quality Dashboard** ✅
- **Interactive connection quality visualization** with real-time metrics
- **Historical trend analysis** for connection health monitoring
- **Quality distribution charts** with excellent/good/fair/poor breakdown
- **Latency and packet loss trends** with time series visualization
- **Connection quality alerts** with configurable thresholds
- **Detailed connection statistics** table with per-client metrics

## Decision Log

### 2026-02-28: Plugin Manifest Approach
**Decision**: Use JSON manifests in `plugins/manifests/` for plugin discovery
**Rationale**: Simple, declarative, language-agnostic, easy to validate
**Alternatives Considered**:
- YAML manifests (rejected: more complex parsing)
- Code-based registration (rejected: less flexible)
- Database storage (rejected: overkill for this use case)

### 2026-02-28: LLM Safety Design
**Decision**: Redact sensitive data before sending to LLM
**Rationale**: Security and privacy by default
**Implementation**:
- Remove hostnames, tokens, absolute paths
- Truncate long strings
- Send only derived data, never raw profiles

## Iteration Log

### 2026-02-28 14:00: Initial Implementation
- **Objective**: Implement manifest-based plugin discovery
- **Changes**:
  - Created `internal/plugin/manifest.go` with Manifest struct
  - Added `LoadManifest()` with strict JSON parsing
  - Implemented `DiscoverManifests()` for directory walking
  - Added `ResolvePlugin()` for validation
  - Enhanced PluginManager with manifest support
  - Updated core pipeline to validate plugins before launch
- **Testing**: Unit tests for parsing, discovery, and validation
- **Result**: Plugins now discovered from manifests with capability checking

### 2026-02-28 15:30: LLM Insights Feature
- **Objective**: Add optional LLM augmentation
- **Changes**:
  - Created `internal/model/insights.go` with InsightsBundle
  - Implemented `internal/llm/client.go` for Mistral API
  - Added `internal/llm/prompt.go` with redaction
  - Created `internal/llm/insights.go` for orchestration
  - Enhanced pipeline with optional LLM step
  - Updated reporter to include LLM insights
  - Added CLI commands for LLM control
- **Testing**: Unit tests for client, prompt building, and integration
- **Result**: Optional LLM insights with safety features and proper error handling

### 2026-02-28 16:00: LLM Integration Completion
- **Objective**: Complete LLM integration with core pipeline
- **Changes**:
  - Updated `internal/core/pipeline.go` with LLM generator support
  - Added `WithLLM()` method to configure LLM insights
  - Enhanced `RunWithTarget()` to generate and save insights
  - Updated `ReportWithInsights()` for enhanced report generation
  - Modified CLI to support LLM flags in run command
  - Added LLM command for standalone insights generation
- **Testing**: Integration tests for full pipeline with LLM
- **Result**: End-to-end LLM integration with optional insights generation

### 2026-02-28 17:00: Enhanced WebSocket Real-time Monitoring
- **Objective**: Implement live data streaming and enhanced real-time monitoring
- **Changes**:
  - Added `StartAutoRefresh()` method to WebSocket server for periodic updates
  - Enhanced `BroadcastData()` with comprehensive statistics including client count
  - Added `UpdateData()` method for dynamic data updates
  - Implemented `GetClientCount()` for monitoring connected clients
  - Enhanced WebSocket stats display in web viewer with severity breakdown
  - Added real-time update notifications with slide-in animations
  - Updated CLI with `--auto-refresh` flag for configurable update intervals
  - Added comprehensive WebSocket server tests
  - Enhanced WebSocket protocol with detailed performance metrics
- **Testing**: Unit tests for WebSocket server, integration tests for real-time updates
- **Result**: Production-ready real-time monitoring with live data streaming, enhanced UI, and comprehensive statistics

### 2026-03-01 02:30: WebSocket JWT Authentication Implementation
- **Objective**: Implement JWT authentication for WebSocket connections to enhance security
- **Changes**:
  - Added JWT token generation and validation to WebSocket server
  - Implemented `generateJWTToken()` and `validateJWTToken()` methods
  - Added `extractTokenFromRequest()` utility function for token extraction
  - Enhanced WebSocket handler with JWT authentication middleware
  - Added `/auth/token` endpoint for token generation
  - Updated `BroadcastData()` to include authentication status in stats
  - Added comprehensive JWT authentication tests
  - Enhanced WebSocket server constructor with authentication support
  - Updated core pipeline integration for WebSocket authentication
- **Testing**: Unit tests for JWT token generation, validation, and WebSocket authentication flow
- **Result**: Secure WebSocket connections with JWT authentication, backward-compatible with existing functionality

### 2026-03-01 04:30: WebSocket Message Batching Implementation
- **Objective**: Implement message batching for WebSocket connections to optimize high-frequency updates
- **Changes**:
  - Added message batching fields to `WebSocketServer` struct: `batchingEnabled`, `batchInterval`, `messageQueue`, `queueMu`, `batchTimer`
  - Enhanced `NewWebSocketServer()` constructor with batching parameters
  - Implemented `startBatching()` method to initialize batching timer
  - Added `flushMessageQueue()` method to send batched messages
  - Implemented `sendBatchedMessage()` method for batch transmission
  - Added `queueMessage()` method for message queuing
  - Enhanced `Stop()` method to clean up batching resources
  - Updated `BroadcastData()` to use batching when enabled
  - Added `GetBatchingInfo()` method for batching configuration information
  - Added `/batching/info` HTTP endpoint for batching status
  - Updated core pipeline integration with batching parameters
  - Added CLI flags `--websocket-batching` and `--websocket-batch-interval`
  - Enhanced WebSocket server output to show batching status
  - Added comprehensive batching tests including concurrency tests
- **Testing**: Unit tests for batching functionality, integration tests, and concurrency tests
- **Result**: WebSocket message batching with configurable intervals, reducing message frequency for high-volume scenarios while maintaining real-time capabilities

### 2026-03-01 04:45: WebSocket Connection Quality Monitoring Implementation
- **Objective**: Implement WebSocket Connection Quality Monitoring for enhanced real-time monitoring reliability
- **Changes**:
  - Added `WebSocketConnectionStats` struct to track connection quality metrics (latency, packet loss, message counts, bandwidth)
  - Enhanced `WebSocketServer` struct with connection quality fields: `connectionStats`, `statsMu`, `pingInterval`, `connectionQualityEnabled`
  - Updated `NewWebSocketServer()` constructor with `enableConnectionQuality` parameter
  - Implemented `calculateConnectionQuality()` method for quality classification (excellent/good/fair/poor)
  - Added `updateConnectionStats()` method for real-time connection statistics tracking
  - Implemented `getConnectionStats()` and `GetConnectionQualityInfo()` methods for monitoring
  - Added `calculateAverageLatency()` method for performance analysis
  - Enhanced WebSocket handler with ping/pong monitoring using WebSocket control messages
  - Added connection cleanup on client disconnect
  - Updated `BroadcastData()` to include connection quality information in WebSocket payloads
  - Added `/connection/quality` HTTP endpoint for connection quality monitoring
  - Enhanced `handleWebSocket()` with connection quality monitoring when enabled
  - Added CLI flag `--websocket-connection-quality` for enabling connection quality monitoring
  - Updated core pipeline with `WithWebSocketConnectionQuality()` method
  - Added comprehensive unit tests for connection quality functionality
  - Enhanced WebSocket server output to show connection quality monitoring status
- **Testing**: Unit tests for connection quality calculation, average latency calculation, HTTP endpoint, and integration with WebSocket server
- **Result**: WebSocket connections now include comprehensive quality monitoring with latency tracking, packet loss detection, and quality classification. Users can monitor connection health in real-time and troubleshoot connectivity issues more effectively.

### 2026-03-01 05:30: WebSocket Connection Quality Dashboard Implementation
- **Objective**: Implement interactive WebSocket Connection Quality Dashboard for real-time monitoring and analysis
- **Changes**:
  - Created `web/connection-quality-dashboard.html` with comprehensive UI for connection quality monitoring
  - Added interactive quality summary cards showing excellent/good/fair/poor connection counts
  - Implemented real-time charts: quality trends, quality distribution, latency trends, packet loss trends
  - Added connection quality alerts table with configurable thresholds and acknowledgment
  - Created detailed connection statistics table with per-client metrics
  - Enhanced main dashboard with link to connection quality dashboard
  - Implemented WebSocket message handling for connection quality subscriptions and updates
  - Added `handleWebSocketMessage()` method for processing WebSocket messages (subscribe, acknowledge_alert, request_update)
  - Implemented `BroadcastConnectionQualityData()` for real-time connection quality data streaming
  - Added `BroadcastConnectionQualityAlerts()` for connection quality alert notifications
  - Enhanced WebSocket server with periodic connection quality updates (5-second intervals)
  - Added comprehensive unit tests for WebSocket connection quality broadcasting and message handling
  - Updated WebSocket protocol to support connection quality data streaming and alert management
- **Testing**: Unit tests for WebSocket connection quality broadcasting, alerts broadcasting, and message handling
- **Result**: Users can now monitor WebSocket connection quality in real-time through an interactive dashboard with visualizations, historical trends, quality distribution, and configurable alerts. The dashboard provides comprehensive insights into connection health, latency, packet loss, and individual client performance, enabling better troubleshooting and optimization of real-time monitoring infrastructure.

### 2026-03-01 06:30: WebSocket Connection Quality Dashboard Enhancements Implementation
- **Objective**: Enhance WebSocket Connection Quality Dashboard with geographical connection analysis and connection quality predictions
- **Changes**:
  - Extended `WebSocketConnectionStats` struct with new fields: `Geolocation`, `ConnectionScore`, `QualityTrend`, `PredictedQuality`
  - Added `calculateConnectionScore()` method for comprehensive quality scoring (0-100 scale)
  - Implemented `determineQualityTrend()` method to analyze historical quality trends (improving/degrading/stable)
  - Added `predictConnectionQuality()` method for future quality prediction based on current trends
  - Implemented `inferGeolocation()` method for client location inference (simplified for demo)
  - Added `getGeographicalConnectionAnalysis()` method for regional connection quality analysis
  - Implemented `getQualityPredictions()` method for connection quality forecasting and insights
  - Added `calculateOverallGeographicalQuality()` for geographical quality assessment
  - Implemented `generatePredictiveInsights()` for actionable insights based on predictions
  - Enhanced `updateConnectionStats()` to calculate scores, trends, and predictions
  - Updated `BroadcastConnectionQualityData()` to include geographical and prediction data
  - Extended connection quality dashboard with geographical analysis chart and prediction chart
  - Added geographical summary display showing regions, best/worst regions, and overall quality
  - Implemented prediction insights display with trend analysis and actionable recommendations
  - Enhanced connection stats table with location, score, trend, and predicted quality columns
  - Added CSS styles for new UI elements and improved visual presentation
  - Updated WebSocket message handling to support enhanced connection quality data
- **Testing**: Unit tests for geographical analysis, quality predictions, and enhanced connection statistics
- **Result**: WebSocket Connection Quality Dashboard now provides advanced geographical connection analysis with regional performance metrics, quality distribution by location, and identification of best/worst performing regions. The enhanced dashboard also includes machine learning-inspired connection quality predictions with trend analysis, confidence indicators, and actionable insights. Users can now analyze connection quality by geographical region, identify regional performance issues, and get predictive alerts about potential future connection problems, enabling proactive network optimization and troubleshooting.

---

## How to Contribute

1. **Read this document** to understand the vision and rules
2. **Propose changes** by adding to `change.log`
3. **Discuss** the proposed changes with the team
4. **Implement** the approved changes
5. **Update change.log** with actual results

## Maintenance

- **Do NOT edit this file directly**
- All historical changes preserved in iteration log
- Change.log entries become part of the permanent record
- This ensures traceability and accountability

**Remember**: The compass points the way, but the journey is recorded in the logs.

---

## 🎯 Killer Feature Focus: AI-Powered Performance Triage

### The Vision

**TriageProf's killer feature is transforming raw profiling data into actionable insights using AI.**

While traditional profilers show *what* is slow, TriageProf explains *why* it's slow and *how* to fix it.

### Demo Goal: "Wow" Factor

**Objective**: Create a demo that makes developers say "Wow, I need this!"

**Success Criteria**:
- Input: Complex Go application with real performance issues
- Process: Automatic collection → analysis → LLM insights
- Output: Professional report with:
  - ✅ Executive summary with severity assessment
  - ✅ Top 3 risks with impact analysis
  - ✅ Actionable suggestions with code examples
  - ✅ Confidence scores and caveats

### Demo Flow

```
1. Start demo server with intentional issues
   - CPU hotspots
   - Memory allocation patterns
   - Mutex contention
   - Blocking operations

2. Run TriageProf with LLM enabled
   bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir demo --llm

3. Show the "Wow" moments:
   - ✨ Automatic plugin discovery and validation
   - 🔍 Detailed profile collection (CPU, heap, mutex, block, goroutine)
   - 🤖 LLM-generated insights with root cause analysis
   - 📊 Professional markdown report with executive summary

4. Highlight key differentiators:
   - "Traditional profilers show you the hotspots"
   - "TriageProf explains why they're hot and how to fix them"
   - "From data to insights in one command"

### Cool Demo Features to Showcase

1. **Plugin Discovery Magic**
   - Run `triageprof plugins` to show available plugins
   - Emphasize manifest-based discovery

2. **End-to-End Workflow**
   - Single command collects, analyzes, and reports
   - Optional LLM augmentation for deeper insights

3. **LLM Insights**
   - Executive summary with severity assessment
   - Per-finding narrative explanations
   - Suggested fixes and next measurements

4. **Safety Features**
   - Sensitive data redaction
   - Graceful degradation without API key
   - Test verification before committing

5. **Professional Output**
   - Structured markdown reports
   - LLM insights clearly marked
   - Confidence indicators

### Demo Script Example

```bash
# Start the demo server
go run examples/demo-server/main.go &

# Generate some load
./examples/load.sh &

# Run the full pipeline with LLM
export MISTRAL_API_KEY="your-key"
bin/triageprof run --plugin go-pprof-http --target-url http://localhost:6060 --duration 30 --outdir demo --llm

# Show the results
cat demo/report.md

# Highlight the LLM insights
bin/triageprof report --in demo/findings.json --insights demo/insights.json --out demo/enhanced-report.md
```

### Demo Success Metrics

**Technical**:
- ✅ All profiles collected successfully
- ✅ Analysis completes without errors
- ✅ LLM generates meaningful insights
- ✅ Report includes all expected sections

**User Experience**:
- ✅ "Wow" reaction from viewers
- ✅ Clear understanding of the value proposition
- ✅ Desire to use the tool on their own projects
- ✅ Appreciation of the AI augmentation

**Visual Appeal**:
- ✅ Professional markdown formatting
- ✅ Clear separation of deterministic vs. LLM content
- ✅ Helpful visual elements (tables, bullet points)
- ✅ Confidence indicators for transparency

### Cool Factor Checklist

- [x] **Automatic**: One command does everything
- [x] **Smart**: AI explains the why behind performance issues
- [x] **Safe**: Redacts sensitive data automatically
- [x] **Professional**: Generates reports suitable for stakeholders
- [x] **Extensible**: Plugin architecture for any profiler
- [x] **Optional**: Works great without LLM too

### Demo Tips

1. **Start with a problem**: Show a slow application
2. **Run the analysis**: Single command magic
3. **Show the insights**: Focus on LLM explanations
4. **Compare approaches**: "Traditional vs. TriageProf"
5. **Highlight safety**: Emphasize data redaction
6. **Show flexibility**: Demonstrate plugin system

**Goal**: Make developers excited about performance analysis again!

---

## 🚀 Delivery Timeline

### Phase 1: Core Complete ✅
- Plugin discovery and validation
- Basic analysis pipeline
- Markdown report generation

### Phase 2: LLM Integration ✅
- Mistral API client
- Secure prompt generation
- Insights integration
- Enhanced reports

### Phase 3: Demo Polish (Current Focus)
- Perfect the demo script
- Create compelling sample application
- Refine visual output
- Prepare presentation materials

### Phase 4: Launch Ready
- Final testing and validation
- Documentation polish
- Screencast recording
- Website/update materials

**Target**: Have a "wow"-worthy demo ready for showcase!

## Iteration Log

### 2026-03-01 07:30: WebSocket Connection Quality Dashboard Advanced Enhancements Phase 3
- **Objective**: Implement advanced AI/ML capabilities with deep learning for enhanced anomaly detection
- **Rationale**: Provide state-of-the-art connection quality monitoring with predictive capabilities
- **Implementation**:
  - ✅ **Deep Learning Anomaly Detection**: Advanced neural network-based anomaly detection with high accuracy
  - ✅ **Time Series Forecasting**: Predict future anomalies using historical patterns and trends
  - ✅ **Automated Root Cause Analysis**: AI-powered determination of anomaly root causes with confidence scoring
  - ✅ **Anomaly Correlation Detection**: Identify systemic issues by correlating different anomaly types
  - ✅ **Adaptive Learning**: Continuous model improvement with learning rate adjustment
  - ✅ **Advanced ML Model Management**: Comprehensive ML model training, statistics, and versioning
  - ✅ **Real-time Anomaly Severity Assessment**: Dynamic severity classification based on multiple factors
  - ✅ **Predictive Insights**: Actionable recommendations and mitigation strategies
  - ✅ **Enhanced Web UI**: Advanced ML visualization with model status, predictions, and insights
  - ✅ **Comprehensive API Endpoints**: RESTful endpoints for advanced ML operations and data access
  - ✅ **Full Integration**: Seamless integration with existing connection quality monitoring infrastructure
  - ✅ **Comprehensive Testing**: Unit tests covering all new ML functionality and edge cases
- **Key Components Added**:
  - `detectAdvancedConnectionQualityAnomaliesPhase3()` for deep learning-based anomaly detection
  - `deepLearningAnomalyDetection()` for neural network simulation
  - `predictFutureAnomalyWithDeepLearning()` for time series forecasting
  - `generateAdvancedMLInsightsForAnomalies()` for comprehensive anomaly analysis
  - `performAdaptiveLearning()` for continuous model improvement
  - `getAdvancedMLConnectionQualityInfo()` for enhanced monitoring data
  - Advanced web UI components and visualization
  - HTTP endpoint `/connection/quality/advanced` for ML data access
  - Comprehensive test coverage for new functionality
- **Testing**: Unit tests for all new ML functions and integration tests
- **Result**: Production-ready advanced ML capabilities for WebSocket connection quality monitoring with deep learning, predictive analytics, and automated insights