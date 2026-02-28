package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

type Plugin struct {
	info model.PluginInfo
}

func main() {
	plugin := &Plugin{
		info: model.PluginInfo{
			Name:       "go-pprof-http",
			Version:    "0.1.0",
			SDKVersion: "1.0",
			Capabilities: model.Capabilities{
				Targets:  []string{"url"},
				Profiles: []string{"cpu", "heap", "mutex", "block", "goroutine", "allocs"},
			},
		},
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			continue
		}

		var req RPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
			continue
		}

		var result interface{}
		var methodErr error

		switch req.Method {
		case "rpc.info":
			result = plugin.info
		case "rpc.validateTarget":
			var target model.Target
			paramsData, err := json.Marshal(req.Params)
			if err != nil {
				result = nil
				break
			}
			if err := json.Unmarshal(paramsData, &target); err != nil {
				result = nil
				break
			}
			methodErr = plugin.validateTarget(target)
		case "rpc.collect":
			var collectReq model.CollectRequest
			paramsData, err := json.Marshal(req.Params)
			if err != nil {
				result = nil
				break
			}
			if err := json.Unmarshal(paramsData, &collectReq); err != nil {
				result = nil
				break
			}
			result, methodErr = plugin.collect(collectReq)
		default:
			methodErr = fmt.Errorf("unknown method: %s", req.Method)
		}

		resp := RPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
		}

		if methodErr != nil {
			resp.Error = &RPCError{
				Code:    -32603,
				Message: methodErr.Error(),
			}
		} else {
			resp.Result = result
		}

		data, _ := json.Marshal(resp)
		fmt.Fprintln(writer, string(data))
		writer.Flush()
	}
}

type RPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
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

func (p *Plugin) validateTarget(target model.Target) error {
	if target.Type != "url" {
		return fmt.Errorf("unsupported target type: %s", target.Type)
	}
	if !strings.HasPrefix(target.BaseURL, "http://") && !strings.HasPrefix(target.BaseURL, "https://") {
		return fmt.Errorf("invalid URL scheme")
	}
	return nil
}

func (p *Plugin) collect(req model.CollectRequest) (model.ArtifactBundle, error) {
	// Create output directory
	if err := os.MkdirAll(req.OutDir, 0755); err != nil {
		return model.ArtifactBundle{}, err
	}

	artifacts := []model.Artifact{}
	client := &http.Client{Timeout: time.Duration(req.DurationSec+10) * time.Second}

	// Collect CPU profile
	if contains(req.Profiles, "cpu") {
		cpuPath := filepath.Join(req.OutDir, "cpu.pb.gz")
		cpuURL := fmt.Sprintf("%s/debug/pprof/profile?seconds=%d", req.Target.BaseURL, req.DurationSec)
		if err := downloadFile(client, cpuURL, cpuPath); err == nil {
			artifacts = append(artifacts, model.Artifact{
				Kind:        "pprof",
				ProfileType: "cpu",
				Path:        cpuPath,
				ContentType: "application/octet-stream",
			})
		}
	}

	// Collect Heap profile
	if contains(req.Profiles, "heap") {
		heapPath := filepath.Join(req.OutDir, "heap.pb.gz")
		heapURL := fmt.Sprintf("%s/debug/pprof/heap", req.Target.BaseURL)
		if err := downloadFile(client, heapURL, heapPath); err == nil {
			artifacts = append(artifacts, model.Artifact{
				Kind:        "pprof",
				ProfileType: "heap",
				Path:        heapPath,
				ContentType: "application/octet-stream",
			})
		}
	}

	// Collect Mutex profile
	if contains(req.Profiles, "mutex") {
		mutexPath := filepath.Join(req.OutDir, "mutex.pb.gz")
		mutexURL := fmt.Sprintf("%s/debug/pprof/mutex", req.Target.BaseURL)
		if err := downloadFile(client, mutexURL, mutexPath); err == nil {
			artifacts = append(artifacts, model.Artifact{
				Kind:        "pprof",
				ProfileType: "mutex",
				Path:        mutexPath,
				ContentType: "application/octet-stream",
			})
		}
	}

	// Collect Block profile
	if contains(req.Profiles, "block") {
		blockPath := filepath.Join(req.OutDir, "block.pb.gz")
		blockURL := fmt.Sprintf("%s/debug/pprof/block", req.Target.BaseURL)
		if err := downloadFile(client, blockURL, blockPath); err == nil {
			artifacts = append(artifacts, model.Artifact{
				Kind:        "pprof",
				ProfileType: "block",
				Path:        blockPath,
				ContentType: "application/octet-stream",
			})
		}
	}

	// Collect Goroutine profile
	if contains(req.Profiles, "goroutine") {
		goroutinePath := filepath.Join(req.OutDir, "goroutine.txt")
		goroutineURL := fmt.Sprintf("%s/debug/pprof/goroutine?debug=2", req.Target.BaseURL)
		if err := downloadFile(client, goroutineURL, goroutinePath); err == nil {
			artifacts = append(artifacts, model.Artifact{
				Kind:        "text",
				ProfileType: "goroutine",
				Path:        goroutinePath,
				ContentType: "text/plain",
			})
		}
	}

	// Collect Allocs profile
	if contains(req.Profiles, "allocs") {
		allocsPath := filepath.Join(req.OutDir, "allocs.pb.gz")
		allocsURL := fmt.Sprintf("%s/debug/pprof/allocs", req.Target.BaseURL)
		if err := downloadFile(client, allocsURL, allocsPath); err == nil {
			artifacts = append(artifacts, model.Artifact{
				Kind:        "pprof",
				ProfileType: "allocs",
				Path:        allocsPath,
				ContentType: "application/octet-stream",
			})
		}
	}

	return model.ArtifactBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: req.DurationSec,
			Service:     req.Metadata["service"],
			Scenario:    req.Metadata["scenario"],
			GitSha:      req.Metadata["gitSha"],
		},
		Target:    req.Target,
		Artifacts: artifacts,
	}, nil
}

func downloadFile(client *http.Client, url, path string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
