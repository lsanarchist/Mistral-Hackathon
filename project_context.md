# TriageProf Project Context - Comprehensive System Documentation

This document provides a complete overview of the TriageProf codebase, architecture, and all significant components. It serves as a single source of truth for AI agents and developers to understand the system's structure, purpose, and implementation details.

## 🎯 Project Overview

**TriageProf** is a production-grade, modular profiling and bottleneck-analysis tool for performance optimization. It integrates with existing profilers via a well-defined plugin SDK, produces evidence-backed bottleneck findings, and generates structured, machine-readable reports with optional AI-powered insights.

### Core Philosophy
- **Go deep, not wide**: Make each layer excellent before adding the next
- **Language-agnostic core**: Language-specific logic lives in plugins
- **AI/LLM as first-class feature**: Not optional glue, but a proper pipeline stage
- **Deterministic first**: Always collect deterministic profiling data as source of truth
- **Extensible architecture**: Plugin-based design for future growth

## 🗺️ Architecture Overview

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                                TriageProf Architecture                            │
├─────────────────┬─────────────────┬─────────────────┬─────────────────┤
│     CLI Layer    │   Core Pipeline  │   Plugin System  │   Web Server    │
├─────────────────┼─────────────────┼─────────────────┼─────────────────┤
│ cmd/triageprof/ │ internal/core/   │ internal/plugin/ │ internal/webserver/│
│ main.go         │ pipeline.go      │ jsonrpc.go       │ server.go       │
│                 │ demo.go          │ manifest.go      │ auth.go         │
│                 │ audit.go         │                 │                 │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
│
│
└─┬───────────────────────────────────────────────────────────────────────────┬─┘
  │                                                                               │
  ▼                                                                               ▼
┌───────────────────────────────────────────────────────────────────────────────┐
│                                Analysis Components                               │
├─────────────────┬─────────────────┬─────────────────┬─────────────────┤
│  Deterministic  │     LLM/AI       │     Reporting    │   Enterprise    │
│    Analysis      │    Insights      │    System        │    Features     │
├─────────────────┼─────────────────┼─────────────────┼─────────────────┤
│ internal/analyzer/│ internal/llm/    │ internal/report/ │ internal/auth/  │
│ analyzer.go     │ insights.go     │ report.go       │ rbac.go         │
│ deterministic.go│ client.go       │                 │                 │
│                 │ prompt.go       │                 │                 │
│                 │ provider.go     │                 │                 │
└─────────────────┴─────────────────┴─────────────────┴─────────────────┘
│
│
└─┬───────────────────────────────────────────────────────────────────────────┬─┘
  │                                                                               │
  ▼                                                                               ▼
┌───────────────────────────────────────────────────────────────────────────────┐
│                                Data Models & Types                               │
├───────────────────────────────────────────────────────────────────────────────┤
│ internal/model/                                                                  │
│ types.go (core types, error handling)                                           │
│ insights.go (LLM insights structures)                                           │
│ remediation.go (automated remediation)                                           │
│ report.go (report structures)                                                    │
└───────────────────────────────────────────────────────────────────────────────┘
```

## 📁 Core Directory Structure

```
.
├── bin/                  # Compiled binaries
├── cmd/triageprof/       # Main CLI application
├── internal/              # Core application logic
│   ├── analyzer/          # Performance analysis engine
│   ├── auth/              # Enterprise authentication & RBAC
│   ├── core/              # Core pipeline and orchestration
│   ├── llm/               # LLM integration and insights
│   ├── model/             # Data models and types
│   ├── plugin/            # Plugin system and SDK
│   ├── report/            # Reporting system
│   └── webserver/         # Web server and WebSocket
├── plugins/               # Plugin manifests and source
│   ├── manifests/         # Plugin manifest JSON files
│   └── src/               # Plugin implementations
├── web/                   # Web interface assets
├── docs/                  # Documentation
├── examples/              # Example applications and demos
└── test-llm-output/       # Test data and outputs
```

## 🔧 Key Components Deep Dive

### 1. CLI Layer (`cmd/triageprof/main.go`)

**Purpose**: Command-line interface and entry point for all TriageProf operations

**Key Functions**:
- `main()`: Entry point with command routing
- `runDemoCommand()`: Execute demo analysis workflows
- `runPluginsCommand()`: Plugin management commands
- `runCollectCommand()`: Profile collection commands
- `runAnalyzeCommand()`: Analysis pipeline execution
- `runLLMCommand()`: LLM insights generation
- `runWebCommand()`: Web report generation
- `runWebSocketCommand()`: Real-time WebSocket server

**Supported Commands**:
- `plugins list`: List available plugins
- `collect`: Collect profiles using plugins
- `analyze`: Analyze collected profiles
- `report`: Generate reports from findings
- `llm`: Generate AI insights
- `run`: Complete collection+analysis pipeline
- `web`: Generate web reports
- `websocket`: Start WebSocket server
- `demo`: Run built-in demo workflows

### 2. Core Pipeline (`internal/core/pipeline.go`)

**Purpose**: Orchestrates the entire profiling and analysis workflow

**Key Components**:
- `Pipeline` struct: Main orchestrator with plugin manager, analyzer, reporter, LLM generator
- `NewPipeline()`: Constructor with default configuration
- `WithLLM()`: Configure LLM insights generation
- `WithPerformanceGates()`: Configure CI/CD performance gates
- `WithEnterpriseConfig()`: Configure enterprise features
- `Run()`: Execute complete profiling pipeline
- `Collect()`: Collect profiles using plugins
- `Analyze()`: Analyze collected profiles
- `GenerateReport()`: Generate reports
- `GenerateInsights()`: Generate LLM insights

**Enterprise Integration**:
- Audit logging via `AuditLogger`
- RBAC via `RBACManager`
- Performance gates for CI/CD
- Team/user management

### 3. Plugin System (`internal/plugin/`)

**Purpose**: Extensible plugin architecture for language-specific profilers

**Key Files**:
- `jsonrpc.go`: JSON-RPC 2.0 communication protocol
- `manifest.go`: Plugin manifest parsing and validation
- `plugin_test.go`: Plugin testing utilities

**Plugin Communication**:
- JSON-RPC 2.0 over stdio
- Methods: `initialize`, `collect`, `analyze`, `shutdown`
- Plugin discovery via manifest files

**Built-in Plugins**:
- `go-pprof-http`: Go HTTP pprof endpoint profiler
- `node-inspector`: Node.js profiling (archived)
- `python-cprofile`: Python cProfile support (archived)
- `ruby-stackprof`: Ruby stackprof integration (archived)

### 4. Analysis Engine (`internal/analyzer/`)

**Purpose**: Performance analysis with deterministic rules and AI insights

**Key Files**:
- `analyzer.go`: Main analysis orchestrator
- `deterministic.go`: 8+ deterministic analysis rules

**Deterministic Rules**:
1. **CPU Hotpath Dominance**: Functions consuming >70% CPU time
2. **Allocation Churn**: High mallocgc/memmove patterns
3. **JSON Hotspots**: encoding/json bottlenecks
4. **String Churn**: strings.Builder/bytes.Buffer usage
5. **GC Pressure**: runtime.gcBgMarkWorker impact
6. **Mutex Contention**: sync.(*Mutex).Lock issues
7. **Heap Allocation**: Memory allocation hotspots
8. **Block Contention**: runtime.chan/select patterns

**Analysis Process**:
1. Parse pprof profiles
2. Extract top functions
3. Build callgraphs (optional)
4. Apply deterministic rules
5. Calculate severity scores
6. Generate structured findings

### 5. LLM Integration (`internal/llm/`)

**Purpose**: AI-powered insights and remediation suggestions

**Key Files**:
- `insights.go`: Insights generation orchestrator
- `client.go`: LLM API client
- `prompt.go`: Prompt construction and guardrails
- `provider.go`: Multi-provider support (Mistral, OpenAI)
- `cache.go`: Insights caching system
- `remediation.go`: Automated code remediation

**Providers Supported**:
- Mistral AI (default)
- OpenAI

**Key Features**:
- Structured prompt building with guardrails
- Insights caching for performance
- Dry-run mode for prompt inspection
- Configurable timeouts and character limits
- Multi-model support

### 6. Reporting System (`internal/report/`)

**Purpose**: Generate professional reports from analysis findings

**Key Files**:
- `report.go`: Report generation engine

**Report Formats**:
- HTML: Interactive web interface
- Markdown: Structured text reports
- JSON: Machine-readable findings
- Raw profiles: Original pprof files

**Report Structure**:
```
output-directory/
├── findings.json          # Structured performance findings
├── insights.json          # LLM insights (if enabled)
├── report.md              # Markdown report
├── report.html            # Interactive HTML report
├── bundle.json            # Complete data bundle
├── run.json               # Run metadata
└── profiles/              # Raw profile files
```

### 7. Web Server (`internal/webserver/`)

**Purpose**: Real-time data streaming and interactive dashboards

**Key Files**:
- `server.go`: WebSocket server implementation
- `auth.go`: Authentication middleware
- `server_new.go`: Enhanced server features

**Key Features**:
- WebSocket-based real-time updates
- JWT authentication
- Connection quality monitoring
- Performance history tracking
- Anomaly detection
- ML-based pattern recognition
- Compression and batching support
- Multi-phase feature sets (Phases 4-6)

### 8. Enterprise Features (`internal/auth/`, `internal/core/audit.go`)

**Purpose**: Team collaboration and enterprise-grade features

**Key Components**:
- **RBAC System** (`internal/auth/rbac.go`):
  - Roles: admin, analyst, viewer
  - Permissions: run_analysis, view_reports, manage_users, manage_teams, configure_system
  - User and team management
  - Permission inheritance

- **Audit Logging** (`internal/core/audit.go`):
  - Action tracking with timestamps
  - User attribution
  - Resource-level logging
  - Persistent storage
  - Log rotation

- **Performance Gates** (`internal/model/types.go`):
  - CI/CD integration
  - Configurable thresholds (critical/high/medium)
  - Build failure on threshold violations
  - Regression analysis

### 9. Data Models (`internal/model/`)

**Purpose**: Core data structures and types

**Key Files**:
- `types.go`: Core types, error handling, configurations
- `insights.go`: LLM insights structures
- `remediation.go`: Automated remediation structures
- `report.go`: Report structures

**Key Data Structures**:
- `ProfileBundle`: Collection of profiles and metadata
- `FindingsBundle`: Analysis findings with evidence
- `InsightsBundle`: LLM-generated insights
- `PerformanceGateConfig`: CI/CD gate configuration
- `EnterpriseConfig`: Enterprise feature configuration
- `ErrorContext`: Structured error handling
- `AuditLogEntry`: Audit logging structure

## 🔄 Pipeline Workflow

```
┌───────────────────────────────────────────────────────────────────────────────┐
│                            TriageProf Pipeline Workflow                          │
├───────────────────────────────────────────────────────────────────────────────┤
│                                                                                   │
│  1. CLI Command Parsing                                                           │
│     └─> cmd/triageprof/main.go                                                   │
│                                                                                   │
│  2. Pipeline Initialization                                                       │
│     └─> internal/core/pipeline.go NewPipeline()                                  │
│                                                                                   │
│  3. Plugin Discovery & Initialization                                            │
│     └─> internal/plugin/manifest.go                                              │
│                                                                                   │
│  4. Profile Collection                                                            │
│     └─> internal/plugin/jsonrpc.go Call()                                        │
│                                                                                   │
│  5. Deterministic Analysis                                                        │
│     └─> internal/analyzer/analyzer.go Analyze()                                 │
│         └─> internal/analyzer/deterministic.go ApplyRules()                     │
│                                                                                   │
│  6. LLM Insights Generation (optional)                                            │
│     └─> internal/llm/insights.go GenerateInsights()                              │
│         └─> internal/llm/prompt.go Build()                                      │
│             └─> internal/llm/client.go GenerateInsights()                        │
│                                                                                   │
│  7. Report Generation                                                             │
│     └─> internal/report/report.go GenerateReport()                               │
│                                                                                   │
│  8. Web Server / WebSocket (optional)                                            │
│     └─> internal/webserver/server.go Start()                                     │
│                                                                                   │
│  9. Audit Logging (enterprise)                                                    │
│     └─> internal/core/audit.go LogAction()                                      │
│                                                                                   │
│  10. Performance Gate Checking (CI/CD)                                            │
│      └─> internal/core/pipeline.go CheckPerformanceGates()                        │
│                                                                                   │
└───────────────────────────────────────────────────────────────────────────────┘
```

## 🛠️ Configuration System

### Environment Variables
- `TRIAGEPROF_PLUGINS`: Custom plugin directory
- `MISTRAL_API_KEY`: Mistral AI API key
- `OPENAI_API_KEY`: OpenAI API key
- `LLM_PROVIDER`: Default LLM provider (mistral, openai)
- `LLM_MODEL`: Default LLM model
- `LLM_TIMEOUT`: LLM request timeout (seconds)
- `LLM_MAX_CHARS`: Maximum prompt characters

### CLI Flags
- `--llm`: Enable LLM insights
- `--llm-provider`: Specify LLM provider
- `--llm-model`: Specify LLM model
- `--llm-timeout`: LLM timeout
- `--llm-max-chars`: Maximum prompt characters
- `--llm-dry-run`: Print prompt without API call
- `--concurrent`: Enable concurrent benchmarks
- `--max-workers`: Maximum concurrent workers
- `--sampling-rate`: Profile sampling rate
- `--memory-optimization`: Enable memory optimization
- `--large-codebase`: Optimize for large codebases
- `--performance-gates`: Enable CI/CD performance gates
- `--enterprise`: Enable enterprise features

## 🔌 Plugin Development

### Plugin Manifest Structure
```json
{
  "name": "go-pprof-http",
  "version": "1.0.0",
  "sdkVersion": "1.0",
  "capabilities": {
    "targets": ["go"],
    "profiles": ["cpu", "heap", "allocs", "block", "mutex"]
  }
}
```

### Plugin Communication Protocol
- **Transport**: JSON-RPC 2.0 over stdio
- **Methods**:
  - `initialize`: Plugin initialization
  - `collect`: Profile collection
  - `analyze`: Profile analysis
  - `shutdown`: Plugin cleanup

### Built-in Plugins

#### go-pprof-http
- **Language**: Go
- **Target**: HTTP pprof endpoints
- **Profiles**: CPU, heap, allocs, block, mutex
- **Location**: `plugins/src/go-pprof-http/main.go`

#### node-inspector (archived)
- **Language**: Node.js
- **Target**: Node.js applications
- **Profiles**: CPU, heap
- **Location**: `plugins/src/node-inspector/main.go`

#### python-cprofile (archived)
- **Language**: Python
- **Target**: Python applications
- **Profiles**: CPU, memory
- **Location**: `plugins/src/python-cprofile/main.py`

#### ruby-stackprof (archived)
- **Language**: Ruby
- **Target**: Ruby applications
- **Profiles**: CPU, memory
- **Location**: `plugins/src/ruby-stackprof/main.rb`

## 🤖 LLM Integration Details

### Prompt Construction
- **Guardrails**: Prevent prompt injection and ensure safety
- **Structure**: Organized sections with clear delimiters
- **Context**: Includes profile data, findings, and analysis context
- **Limitations**: Configurable maximum character limits

### Insights Generation
- **Providers**: Mistral (default), OpenAI
- **Caching**: Insights cached for performance
- **Dry-run**: Generate prompts without API calls for testing
- **Timeout**: Configurable request timeouts
- **Error Handling**: Graceful degradation when LLM fails

### Remediation Features
- **Code Suggestions**: Concrete code examples for fixes
- **Confidence Scoring**: Rank suggestions by confidence
- **Impact Analysis**: Estimate performance impact
- **Safety Checks**: Validate suggestions before applying

## 🏢 Enterprise Features

### RBAC System
- **Roles**: admin, analyst, viewer
- **Permissions**: Fine-grained access control
- **Teams**: User grouping and management
- **Audit Trail**: Complete action logging

### Audit Logging
- **Storage**: JSON-based log files
- **Rotation**: Configurable maximum entries
- **Format**: Structured entries with timestamps
- **Persistence**: Survives application restarts

### Performance Gates
- **CI/CD Integration**: Fail builds on performance regressions
- **Thresholds**: Configurable severity levels
- **Regression Analysis**: Compare against baselines
- **Reporting**: Detailed gate results

## 📊 Analysis Rules

### 1. CPU Hotpath Dominance
- **Trigger**: Function consumes >70% CPU time
- **Severity**: Critical
- **Evidence**: Stack traces, CPU samples
- **Remediation**: Optimization suggestions

### 2. Allocation Churn
- **Trigger**: High mallocgc/memmove patterns
- **Severity**: High
- **Evidence**: Allocation profiles
- **Remediation**: Memory pool suggestions

### 3. JSON Hotspots
- **Trigger**: encoding/json bottlenecks
- **Severity**: Medium
- **Evidence**: CPU profiles with JSON functions
- **Remediation**: Alternative serialization

### 4. String Churn
- **Trigger**: Inefficient string operations
- **Severity**: Medium
- **Evidence**: strings.Builder/bytes.Buffer usage
- **Remediation**: Builder pattern suggestions

### 5. GC Pressure
- **Trigger**: High runtime.gcBgMarkWorker impact
- **Severity**: High
- **Evidence**: GC-related CPU usage
- **Remediation**: Allocation reduction

### 6. Mutex Contention
- **Trigger**: sync.(*Mutex).Lock contention
- **Severity**: High
- **Evidence**: Mutex profiles
- **Remediation**: Lock-free alternatives

### 7. Heap Allocation
- **Trigger**: Memory allocation hotspots
- **Severity**: Medium
- **Evidence**: Heap profiles
- **Remediation**: Object pool suggestions

### 8. Block Contention
- **Trigger**: runtime.chan/select patterns
- **Severity**: Medium
- **Evidence**: Block profiles
- **Remediation**: Channel optimization

## 🌐 Web Server Features

### WebSocket Server
- **Real-time Updates**: Streaming analysis results
- **Authentication**: JWT-based security
- **Compression**: Reduce bandwidth usage
- **Batching**: Optimize message delivery
- **Connection Quality**: Monitor client connections

### Advanced Features
- **Phase 4**: Anomaly detection, ML patterns
- **Phase 5**: Advanced ML, root cause analysis
- **Phase 6**: Enhanced ML, time series analysis
- **Performance History**: Track metrics over time
- **Alerting**: Configurable performance alerts

## 🔧 Build System

### Makefile Targets
- `make build`: Build main binary
- `make plugins`: Build all plugins
- `make demo`: Run built-in demo
- `make test`: Run tests
- `make clean`: Clean build artifacts

### Build Process
1. Go module initialization
2. Dependency resolution
3. Main binary compilation
4. Plugin compilation
5. Asset bundling

## 🧪 Testing Strategy

### Test Coverage
- **Unit Tests**: Core components (`*_test.go` files)
- **Integration Tests**: Pipeline workflows
- **Plugin Tests**: Plugin functionality
- **Web Server Tests**: WebSocket functionality
- **Performance Tests**: Benchmark analysis

### Test Files
- `internal/analyzer/analyzer_test.go`: Analysis engine tests
- `internal/core/pipeline_test.go`: Pipeline tests
- `internal/llm/*_test.go`: LLM integration tests
- `internal/webserver/*_test.go`: Web server tests
- `internal/auth/rbac_test.go`: RBAC tests
- `internal/core/audit_test.go`: Audit logging tests

## 📈 Performance Optimization

### Features
- **Concurrent Execution**: Parallel benchmark runs
- **Profile Sampling**: Reduce overhead
- **Memory Optimization**: Efficient data structures
- **Caching**: LLM insights caching
- **Batching**: WebSocket message batching
- **Compression**: Data compression

### Configuration
- `--concurrent`: Enable concurrent execution
- `--max-workers`: Set worker count
- `--sampling-rate`: Set sampling rate
- `--memory-optimization`: Enable optimizations
- `--large-codebase`: Optimize for large projects

## 🚀 Deployment Options

### Standalone Binary
- Single executable with embedded assets
- Portable across platforms
- No external dependencies

### Docker Container
- Containerized deployment
- Pre-configured environments
- Easy CI/CD integration

### CI/CD Integration
- GitHub Actions workflows
- Performance gate checking
- Automated reporting
- Build failure on regressions

## 🔒 Security Features

### Authentication
- JWT-based WebSocket authentication
- Role-based access control
- Permission management

### Data Protection
- Secure API key handling
- Input validation
- Prompt injection prevention
- Error handling with context

### Audit Trail
- Complete action logging
- User attribution
- Resource tracking
- Persistent storage

## 📁 Significant Files Reference

### Core Application
- `cmd/triageprof/main.go`: CLI entry point
- `internal/core/pipeline.go`: Core pipeline orchestrator
- `internal/core/demo.go`: Demo workflow implementation
- `internal/core/audit.go`: Audit logging system

### Analysis Engine
- `internal/analyzer/analyzer.go`: Analysis orchestrator
- `internal/analyzer/deterministic.go`: Deterministic rules
- `internal/analyzer/analyzer_test.go`: Analysis tests

### LLM Integration
- `internal/llm/insights.go`: Insights generation
- `internal/llm/client.go`: LLM client
- `internal/llm/prompt.go`: Prompt construction
- `internal/llm/provider.go`: Provider management
- `internal/llm/cache.go`: Insights caching
- `internal/llm/remediation.go`: Automated remediation

### Plugin System
- `internal/plugin/jsonrpc.go`: JSON-RPC communication
- `internal/plugin/manifest.go`: Manifest parsing
- `internal/plugin/plugin_test.go`: Plugin tests

### Web Server
- `internal/webserver/server.go`: WebSocket server
- `internal/webserver/auth.go`: Authentication
- `internal/webserver/server_new.go`: Enhanced features
- `internal/webserver/server_test.go`: Web server tests

### Enterprise Features
- `internal/auth/rbac.go`: RBAC system
- `internal/auth/rbac_test.go`: RBAC tests
- `internal/core/audit.go`: Audit logging
- `internal/core/audit_test.go`: Audit logging tests

### Data Models
- `internal/model/types.go`: Core types and error handling
- `internal/model/insights.go`: LLM insights structures
- `internal/model/remediation.go`: Remediation structures
- `internal/model/report.go`: Report structures

### Reporting
- `internal/report/report.go`: Report generation
- `internal/report/report_test.go`: Report tests

### Plugins
- `plugins/manifests/go-pprof-http.json`: Go plugin manifest
- `plugins/src/go-pprof-http/main.go`: Go plugin implementation
- `plugins/manifests/node-inspector.json`: Node.js plugin manifest
- `plugins/src/node-inspector/main.go`: Node.js plugin implementation

### Web Interface
- `web/app.js`: Main web application
- `web/report.js`: Report visualization
- `web/visualization.js`: Data visualization
- `web/style.css`: Styling

### Documentation
- `README.md`: Main documentation
- `docs/USER_GUIDE.md`: User guide
- `docs/CLI_REFERENCE.md`: CLI reference
- `docs/API_DOCUMENTATION.md`: Plugin API documentation
- `docs/CONTRIBUTING.md`: Contribution guidelines
- `COMPASS.md`: Project direction and philosophy
- `AGENTS.md`: Development guidelines
- `change.log`: Recent changes and implementation details

## 🎯 Key Design Decisions

### 1. Plugin Architecture
- **Rationale**: Language-agnostic core with extensible plugins
- **Benefit**: Support multiple languages without core changes
- **Tradeoff**: Plugin protocol must remain stable

### 2. Deterministic First
- **Rationale**: Always collect objective profiling data first
- **Benefit**: LLM insights are grounded in real data
- **Tradeoff**: Requires proper profiling setup

### 3. LLM as Pipeline Stage
- **Rationale**: AI insights are first-class feature, not bolt-on
- **Benefit**: Clean integration with analysis workflow
- **Tradeoff**: Requires careful prompt engineering

### 4. Structured Error Handling
- **Rationale**: Better debugging and user experience
- **Benefit**: Clear error messages with context
- **Tradeoff**: More complex error handling code

### 5. Enterprise Features
- **Rationale**: Support team collaboration and compliance
- **Benefit**: Audit trails, RBAC, performance gates
- **Tradeoff**: Increased complexity for simple use cases

## 🔮 Future Roadmap

### Short-term Priorities
- Enhance plugin SDK stability
- Improve LLM prompt engineering
- Expand deterministic analysis rules
- Optimize WebSocket performance
- Enhance enterprise features

### Long-term Vision
- Advanced ML-based anomaly detection
- Automated remediation workflows
- Cloud-based analysis service
- IDE integration
- Continuous profiling support

## 📚 Glossary

- **pprof**: Google's profiling tool for Go programs
- **JSON-RPC**: Remote procedure call protocol over JSON
- **RBAC**: Role-Based Access Control
- **JWT**: JSON Web Token for authentication
- **CI/CD**: Continuous Integration/Continuous Deployment
- **LLM**: Large Language Model
- **ML**: Machine Learning
- **SDK**: Software Development Kit
- **API**: Application Programming Interface

## 🎓 Understanding the System

### For New Developers
1. Start with `README.md` for overview
2. Review `COMPASS.md` for philosophy
3. Examine `cmd/triageprof/main.go` for CLI structure
4. Study `internal/core/pipeline.go` for core workflow
5. Look at `internal/analyzer/` for analysis logic
6. Explore `internal/llm/` for AI integration
7. Review `internal/plugin/` for extensibility

### For AI Agents
- This document provides complete context for decision-making
- All significant files and their purposes are documented
- Architecture diagrams show component relationships
- Key design decisions explain rationale
- Future roadmap indicates direction

### For Enterprise Users
- RBAC system provides fine-grained access control
- Audit logging ensures compliance
- Performance gates integrate with CI/CD
- Team management supports collaboration

## 🔍 Decision-Making Guide

### When to Modify Core
- Adding new analysis capabilities
- Enhancing pipeline orchestration
- Improving error handling
- Adding enterprise features

### When to Create Plugins
- Supporting new languages
- Adding new profiler integrations
- Extending collection capabilities

### When to Enhance LLM
- Improving prompt quality
- Adding new insight types
- Enhancing remediation suggestions

### When to Update Web Server
- Adding real-time features
- Enhancing visualization
- Improving performance monitoring

This comprehensive documentation should enable AI agents and developers to make informed decisions about the TriageProf system architecture, implementation, and future enhancements.