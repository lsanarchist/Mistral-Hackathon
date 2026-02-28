package plugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
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

type PluginManager struct {
	PluginDir string
}

func NewPluginManager(pluginDir string) *PluginManager {
	return &PluginManager{PluginDir: pluginDir}
}

func (m *PluginManager) ListPlugins() ([]model.PluginInfo, error) {
	// TODO: implement plugin discovery
	return nil, nil
}

func (m *PluginManager) LaunchPlugin(name string, timeout time.Duration) (*JSONRPCCodec, error) {
	pluginPath := fmt.Sprintf("%s/bin/%s", m.PluginDir, name)
	cmd := exec.Command(pluginPath)
	return NewJSONRPCCodec(cmd)
}