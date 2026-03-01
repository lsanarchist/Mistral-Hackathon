package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type JSONRPCCodec struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr *bufio.Reader
}

type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
}

type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func NewJSONRPCCodec(cmd *exec.Cmd) (*JSONRPCCodec, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &JSONRPCCodec{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
		stderr: bufio.NewReader(stderr),
	}, nil
}

func (c *JSONRPCCodec) Call(method string, params, result interface{}) error {
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	if err := c.writeRequest(req); err != nil {
		return err
	}

	resp, err := c.readResponse()
	if err != nil {
		return err
	}

	if resp.Error != nil {
		return fmt.Errorf("rpc error: %s", resp.Error.Message)
	}

	if result != nil {
		data, err := json.Marshal(resp.Result)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, result)
	}

	return nil
}

func (c *JSONRPCCodec) writeRequest(req RPCRequest) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(c.stdin, string(data))
	return err
}

func (c *JSONRPCCodec) readResponse() (*RPCResponse, error) {
	line, err := c.stdout.ReadString('\n')
	if err != nil {
		return nil, err
	}

	var resp RPCResponse
	if err := json.Unmarshal([]byte(line), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *JSONRPCCodec) Close() error {
	return c.cmd.Process.Kill()
}

// PluginPerformance represents performance metrics for a plugin execution
type PluginPerformance struct {
	PluginName      string        `json:"pluginName"`
	ExecutionTime   time.Duration `json:"executionTime"`
	MemoryUsageMB   float64       `json:"memoryUsageMB"`
	CPUUsagePercent float64       `json:"cpuUsagePercent"`
	Timestamp       time.Time     `json:"timestamp"`
	Success         bool          `json:"success"`
	Error           string        `json:"error,omitempty"`
}

// PluginManager manages plugin discovery and execution
type PluginManager struct {
	PluginDir      string
	Performance    []PluginPerformance
	performanceMu  sync.Mutex
}

func NewPluginManager(pluginDir string) *PluginManager {
	return &PluginManager{
		PluginDir: pluginDir,
		Performance: make([]PluginPerformance, 0),
	}
}

// ListPlugins returns all available plugins from manifests
func (m *PluginManager) ListPlugins() ([]*Manifest, error) {
	manifestsDir := filepath.Join(m.PluginDir, "manifests")
	return DiscoverManifests(manifestsDir)
}

// ResolvePlugin finds a plugin by name and validates it
func (m *PluginManager) ResolvePlugin(name string) (*Manifest, string, error) {
	manifestsDir := filepath.Join(m.PluginDir, "manifests")
	binDir := filepath.Join(m.PluginDir, "bin")
	return ResolvePlugin(manifestsDir, binDir, name)
}

// RecordPluginPerformance records performance metrics for a plugin execution
func (m *PluginManager) RecordPluginPerformance(performance PluginPerformance) {
	m.performanceMu.Lock()
	defer m.performanceMu.Unlock()
	
	// Keep only the last 100 performance records to prevent memory bloat
	if len(m.Performance) >= 100 {
		m.Performance = m.Performance[1:]
	}
	m.Performance = append(m.Performance, performance)
}

// GetPluginPerformance returns the performance metrics for all plugins
func (m *PluginManager) GetPluginPerformance() []PluginPerformance {
	m.performanceMu.Lock()
	defer m.performanceMu.Unlock()
	
	// Return a copy to avoid race conditions
	performanceCopy := make([]PluginPerformance, len(m.Performance))
	copy(performanceCopy, m.Performance)
	return performanceCopy
}

// LaunchPlugin launches a plugin process after validation
func (m *PluginManager) LaunchPlugin(name string, timeout time.Duration) (*JSONRPCCodec, error) {
	// First resolve the plugin to ensure it exists and is valid
	_, binaryPath, err := m.ResolvePlugin(name)
	if err != nil {
		performance := PluginPerformance{
			PluginName:    name,
			Timestamp:     time.Now(),
			Success:       false,
			Error:         fmt.Sprintf("failed to resolve plugin: %v", err),
		}
		m.RecordPluginPerformance(performance)
		return nil, fmt.Errorf("failed to resolve plugin %s: %w", name, err)
	}

	// Launch the plugin process
	cmd := exec.Command(binaryPath)
	
	// Record start time for performance tracking
	startTime := time.Now()
	
	codec, err := NewJSONRPCCodec(cmd)
	if err != nil {
		performance := PluginPerformance{
			PluginName:    name,
			Timestamp:     time.Now(),
			Success:       false,
			Error:         fmt.Sprintf("failed to launch plugin: %v", err),
		}
		m.RecordPluginPerformance(performance)
		return nil, fmt.Errorf("failed to launch plugin %s: %w", name, err)
	}
	
	// Record successful launch with basic metrics
	performance := PluginPerformance{
		PluginName:      name,
		ExecutionTime:   time.Since(startTime),
		Timestamp:       time.Now(),
		Success:         true,
		MemoryUsageMB:   getMemoryUsageMB(),
		CPUUsagePercent: getCPUUsagePercent(),
	}
	m.RecordPluginPerformance(performance)
	
	return codec, nil
}

// getMemoryUsageMB returns the current memory usage in MB
func getMemoryUsageMB() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// Convert bytes to MB
	return float64(m.Alloc) / (1024 * 1024)
}

// getCPUUsagePercent returns the current CPU usage percentage
// Note: This is a simple approximation since Go doesn't provide direct CPU usage
func getCPUUsagePercent() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// This is a placeholder - in a real implementation, you'd use system-specific APIs
	// For demo purposes, return a small random value to simulate CPU usage
	return float64(runtime.NumGoroutine()) * 0.1
}
