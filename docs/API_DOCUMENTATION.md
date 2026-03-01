# TriageProf API Documentation

## Table of Contents

- [Plugin SDK](#plugin-sdk)
- [JSON-RPC Protocol](#json-rpc-protocol)
- [Data Models](#data-models)
- [Plugin Manifest](#plugin-manifest)
- [Development Guide](#development-guide)

## Plugin SDK

TriageProf uses a plugin architecture where profilers are separate executables that communicate via JSON-RPC over stdio.

### Plugin Requirements

1. **Executable**: Must be a standalone binary
2. **JSON-RPC 2.0**: Must implement the JSON-RPC protocol
3. **Stdio Communication**: Must read/write JSON-RPC messages via stdin/stdout
4. **Manifest**: Must have a corresponding JSON manifest file

### Supported Methods

#### `initialize`

Initialize the plugin with configuration.

**Parameters:**
```json
{
  "config": {
    "target": "string",
    "timeout": "number",
    "outputDir": "string"
  }
}
```

**Returns:**
```json
{
  "capabilities": {
    "profileTypes": ["cpu", "heap", "allocs", "block", "mutex"],
    "features": ["concurrent", "sampling"]
  }
}
```

#### `collectProfile`

Collect a specific profile type.

**Parameters:**
```json
{
  "profileType": "cpu",
  "config": {
    "duration": "number",
    "samplingRate": "number"
  }
}
```

**Returns:**
```json
{
  "profilePath": "string",
  "metadata": {
    "type": "string",
    "duration": "number",
    "timestamp": "string"
  }
}
```

#### `shutdown`

Clean up resources.

**Returns:**
```json
{
  "status": "ok"
}
```

## JSON-RPC Protocol

### Message Format

All messages follow the JSON-RPC 2.0 specification:

```json
{
  "jsonrpc": "2.0",
  "id": "number",
  "method": "string",
  "params": {}
}
```

### Response Format

```json
{
  "jsonrpc": "2.0",
  "id": "number",
  "result": {},
  "error": {
    "code": "number",
    "message": "string",
    "data": {}
  }
}
```

### Error Codes

| Code | Meaning |
|------|---------|
| -32600 | Invalid Request |
| -32601 | Method not found |
| -32602 | Invalid params |
| -32603 | Internal error |
| -32000 | Plugin-specific error |

## Data Models

### Finding

```go
type Finding struct {
    ID               string
    Title            string
    Category         string // cpu, alloc, heap, gc, mutex, block
    Severity         string // low, medium, high, critical
    Confidence       float64 // 0.0-1.0
    ImpactSummary    string
    Evidence         []EvidenceItem
    DeterministicHints []string
    Tags             []string
}
```

### EvidenceItem

```go
type EvidenceItem struct {
    Type        string  // stack, metric, pattern, callgraph
    Description string
    Value       string  // JSON-encoded data
    Weight      float64 // 0.0-1.0 importance
}
```

### RunManifest

```go
type RunManifest struct {
    Version         string
    Timestamp       string
    ToolVersion     string
    GoVersion       string
    RepoURL         string
    RepoRef         string
    BenchmarksFound int
    ProfilesGenerated []ProfileInfo
    FindingsCount   int
    PerformanceConfig PerformanceOptimizationConfig
    ErrorContext     *ErrorContext
}
```

## Plugin Manifest

### Manifest Structure

```json
{
  "name": "go-pprof-http",
  "version": "1.0.0",
  "description": "Go pprof HTTP profiler",
  "author": "TriageProf Team",
  "license": "MIT",
  "homepage": "https://github.com/triageprof/triageprof",
  "binaries": {
    "linux/amd64": "bin/go-pprof-http",
    "darwin/amd64": "bin/go-pprof-http-mac",
    "windows/amd64": "bin/go-pprof-http.exe"
  },
  "capabilities": {
    "profileTypes": ["cpu", "heap", "allocs", "block", "mutex"],
    "features": ["concurrent", "sampling"],
    "configSchema": {
      "type": "object",
      "properties": {
        "target": {"type": "string"},
        "timeout": {"type": "number"}
      }
    }
  },
  "documentation": "https://github.com/triageprof/triageprof/blob/main/docs/plugins/go-pprof-http.md"
}
```

### Manifest Fields

| Field | Type | Description |
|-------|------|-------------|
| name | string | Plugin name (unique identifier) |
| version | string | Semantic version |
| description | string | Human-readable description |
| author | string | Author/team name |
| license | string | License identifier |
| homepage | string | Project URL |
| binaries | object | Platform-specific binary paths |
| capabilities | object | Supported features and profile types |
| documentation | string | Documentation URL |

## Development Guide

### Creating a New Plugin

1. **Initialize Plugin Structure**
   ```bash
   mkdir plugins/src/my-plugin
   cd plugins/src/my-plugin
   go mod init github.com/triageprof/my-plugin
   ```

2. **Implement JSON-RPC Handler**
   ```go
   package main

   import (
       "encoding/json"
       "fmt"
       "os"
   )

   type JSONRPCRequest struct {
       JSONRPC string          `json:"jsonrpc"`
       ID     int             `json:"id"`
       Method string          `json:"method"`
       Params json.RawMessage `json:"params"`
   }

   type JSONRPCResponse struct {
       JSONRPC string      `json:"jsonrpc"`
       ID    int         `json:"id"`
       Result interface{} `json:"result,omitempty"`
       Error  *ErrorObj   `json:"error,omitempty"`
   }

   func main() {
       // Read JSON-RPC requests from stdin
       // Process initialize, collectProfile, shutdown methods
       // Write JSON-RPC responses to stdout
   }
   ```

3. **Create Manifest File**
   ```bash
   # Create manifest.json in plugins/manifests/
   {
     "name": "my-plugin",
     "version": "1.0.0",
     "description": "My custom profiler",
     "capabilities": {
       "profileTypes": ["cpu", "memory"]
     }
   }
   ```

4. **Build and Test**
   ```bash
   # Build plugin
   go build -o ../../bin/my-plugin
   
   # Test with TriageProf
   triageprof plugin test my-plugin
   ```

### Plugin Best Practices

1. **Error Handling**: Return proper JSON-RPC error codes
2. **Timeout Handling**: Respect timeout parameters
3. **Resource Cleanup**: Implement proper shutdown
4. **Validation**: Validate all input parameters
5. **Logging**: Use stderr for logging (not stdout)
6. **Performance**: Optimize for minimal overhead

### Testing Your Plugin

```bash
# Test plugin initialization
triageprof plugin test my-plugin --method initialize

# Test profile collection
triageprof plugin test my-plugin --method collectProfile --params '{"profileType": "cpu"}'

# Test with real analysis
triageprof demo --repo ./test-app --plugin my-plugin --out analysis/
```

## Integration with Core

### How TriageProf Uses Plugins

1. **Discovery**: Scans `plugins/manifests/` for available plugins
2. **Initialization**: Calls `initialize` with configuration
3. **Profile Collection**: Calls `collectProfile` for each profile type
4. **Shutdown**: Calls `shutdown` for cleanup
5. **Analysis**: Processes collected profiles with analyzers

### Plugin Configuration

```json
{
  "plugins": [
    {
      "name": "go-pprof-http",
      "config": {
        "target": "http://localhost:6060",
        "timeout": 30
      }
    }
  ]
}
```

## Advanced Topics

### Concurrent Profile Collection

Plugins can support concurrent collection:

```json
{
  "capabilities": {
    "features": ["concurrent"]
  }
}
```

### Profile Sampling

```json
{
  "capabilities": {
    "features": ["sampling"]
  }
}
```

### Custom Profile Types

Define custom profile types in capabilities:

```json
{
  "capabilities": {
    "profileTypes": ["cpu", "heap", "custom-metric"]
  }
}
```

## Troubleshooting

### Common Plugin Issues

#### Plugin not discovered
- **Cause**: Missing or invalid manifest file
- **Solution**: Verify manifest JSON structure and location

#### JSON-RPC errors
- **Cause**: Malformed JSON or protocol violations
- **Solution**: Validate JSON structure and error codes

#### Timeout errors
- **Cause**: Plugin takes too long to respond
- **Solution**: Implement proper timeout handling

#### Profile collection failures
- **Cause**: Target application not running or misconfigured
- **Solution**: Verify target configuration and connectivity

## Example: Go pprof HTTP Plugin

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type InitializeParams struct {
    Config map[string]interface{} `json:"config"`
}

type CollectProfileParams struct {
    ProfileType string                 `json:"profileType"`
    Config      map[string]interface{} `json:"config"`
}

func main() {
    decoder := json.NewDecoder(os.Stdin)
    encoder := json.NewEncoder(os.Stdout)

    for {
        var req JSONRPCRequest
        if err := decoder.Decode(&req); err != nil {
            if err == io.EOF {
                break
            }
            sendError(encoder, req.ID, -32600, "Invalid request")
            continue
        }

        switch req.Method {
        case "initialize":
            handleInitialize(encoder, req)
        case "collectProfile":
            handleCollectProfile(encoder, req)
        case "shutdown":
            handleShutdown(encoder, req)
        default:
            sendError(encoder, req.ID, -32601, "Method not found")
        }
    }
}

func handleCollectProfile(encoder *json.Encoder, req JSONRPCRequest) {
    var params CollectProfileParams
    if err := json.Unmarshal(req.Params, &params); err != nil {
        sendError(encoder, req.ID, -32602, "Invalid params")
        return
    }

    // Collect profile based on profileType
    // Return profile path and metadata
}
```

## Resources

- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Go JSON-RPC Examples](https://pkg.go.dev/encoding/json#example-Unmarshal)
- [TriageProf Plugin Examples](https://github.com/triageprof/triageprof/tree/main/plugins/src)
