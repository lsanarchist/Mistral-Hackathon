package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mistral-hackathon/triageprof/internal/model"
)

// RunManifest represents metadata about a demo run
type RunManifest struct {
	RepoURL      string    `json:"repoUrl"`
	RepoRef      string    `json:"repoRef,omitempty"`
	LocalPath    string    `json:"localPath"`
	Timestamp    time.Time `json:"timestamp"`
	GoVersion    string    `json:"goVersion"`
	Benchmarks   []string  `json:"benchmarks"`
	Profiles     []string  `json:"profiles"`
	DurationSec  int       `json:"durationSec"`
	Success      bool      `json:"success"`
	Error        string    `json:"error,omitempty"`
}

// cloneRepo clones a Git repository to the specified directory
func cloneRepo(ctx context.Context, repoURL, ref, destDir string) error {
	// Remove existing directory if it exists
	if _, err := os.Stat(destDir); err == nil {
		if err := os.RemoveAll(destDir); err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Build git clone command
	cmd := exec.CommandContext(ctx, "git", "clone", repoURL, destDir)
	if ref != "" {
		cmd.Args = append(cmd.Args, "--branch", ref)
	}

	// Set up output capture
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	return nil
}

// detectBenchmarks finds Go benchmark functions in the repository
func detectBenchmarks(ctx context.Context, repoPath string) ([]string, error) {
	var benchmarks []string

	// Find all Go files in the repository
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Look for benchmark functions
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "func Benchmark") {
				// Extract benchmark name
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					benchName := parts[1]
					// Remove the (b *testing.B) part
					benchName = strings.Split(benchName, "(")[0]
					benchmarks = append(benchmarks, benchName)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk repository: %w", err)
	}

	return benchmarks, nil
}

// runBenchmarks runs Go benchmarks and collects profiles
func runBenchmarks(ctx context.Context, repoPath string, durationSec int, outDir string) ([]string, error) {
	var profilePaths []string

	// Change to repository directory
	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		return nil, fmt.Errorf("failed to change to repository directory: %w", err)
	}

	// Create profiles directory in repo
	profilesDir := filepath.Join(repoPath, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profiles directory: %w", err)
	}
	
	// Create profiles directory in output
	outputProfilesDir := filepath.Join(outDir, "profiles")
	if err := os.MkdirAll(outputProfilesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output profiles directory: %w", err)
	}

	// Run benchmarks with profiling - use shorter duration for each profile type
	// to avoid timeout issues, and run a subset of benchmarks
	profileTypes := []string{"cpu", "heap", "allocs"} // Reduced set for demo
	benchmarkDuration := durationSec / len(profileTypes) // Split duration among profile types
	if benchmarkDuration < 1 {
		benchmarkDuration = 1 // Minimum 1 second per profile
	}

	for _, profileType := range profileTypes {
		profilePath := filepath.Join(profilesDir, fmt.Sprintf("%s.pprof", profileType))

		// Build benchmark command - run a subset of benchmarks for demo
		args := []string{"test", "-bench=BenchmarkProcessStrings|BenchmarkGenerateRandomData|BenchmarkProcessJSON", "-benchtime", fmt.Sprintf("%ds", benchmarkDuration)}
		if profileType == "cpu" {
			args = append(args, "-cpuprofile", profilePath)
		} else if profileType == "heap" {
			args = append(args, "-memprofile", profilePath)
		} else if profileType == "allocs" {
			args = append(args, "-memprofile", profilePath)
		} else if profileType == "block" {
			args = append(args, "-blockprofile", profilePath)
		} else if profileType == "mutex" {
			args = append(args, "-mutexprofile", profilePath)
		}

		cmd := exec.CommandContext(ctx, "go", args...)
		var stdout, stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			// Benchmarks might fail but still produce profiles
			// Only return error if the profile file doesn't exist
			if _, statErr := os.Stat(profilePath); statErr != nil {
				log.Printf("Benchmark command failed: %v\nArgs: %v\nStdout: %s\nStderr: %s", err, args, stdout.String(), stderr.String())
				return nil, fmt.Errorf("benchmark failed and no profile generated: %w\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
			}
		}

		// Check if profile file was created and copy it immediately
		if _, err := os.Stat(profilePath); err == nil {
			// Copy the profile to a safe location immediately
			destProfilePath := filepath.Join(outputProfilesDir, filepath.Base(profilePath))
			if err := copyFile(profilePath, destProfilePath); err != nil {
				log.Printf("Warning: Failed to copy profile %s to %s: %v", profilePath, destProfilePath, err)
			} else {
				profilePaths = append(profilePaths, destProfilePath)
			}
		}
	}

	return profilePaths, nil
}

// collectProfiles collects profiles from the repository
func (p *Pipeline) collectProfiles(ctx context.Context, repoPath string, durationSec int, outDir string) (*model.ProfileBundle, error) {
	// Detect benchmarks
	benchmarks, err := detectBenchmarks(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect benchmarks: %w", err)
	}

	if len(benchmarks) == 0 {
		return nil, fmt.Errorf("no benchmarks found in repository")
	}

	// Run benchmarks and collect profiles
	profilePaths, err := runBenchmarks(ctx, repoPath, durationSec, outDir)
	if err != nil {
		return nil, fmt.Errorf("failed to run benchmarks: %w", err)
	}

	// Create profile bundle
	bundle := &model.ProfileBundle{
		Metadata: model.Metadata{
			Timestamp:   time.Now(),
			DurationSec: durationSec,
			Service:     "demo",
			Scenario:    "benchmark",
		},
		Target: model.Target{
			Type:    "go",
			BaseURL: repoPath,
		},
		Plugin: model.PluginRef{
			Name: "go-pprof",
		},
	}

	// Add artifacts
	for _, profilePath := range profilePaths {
		profileType := strings.TrimSuffix(filepath.Base(profilePath), ".pprof")
		bundle.Artifacts = append(bundle.Artifacts, model.Artifact{
			Kind:        "profile",
			ProfileType: profileType,
			Path:        profilePath,
			ContentType: "application/octet-stream",
		})
	}

	return bundle, nil
}

// Demo runs the complete demo workflow: clone repo, detect benchmarks, collect profiles, and generate reports
func (p *Pipeline) Demo(ctx context.Context, repoURL, ref, outDir string, durationSec int) (*RunManifest, error) {
	manifest := &RunManifest{
		RepoURL:     repoURL,
		RepoRef:     ref,
		Timestamp:   time.Now(),
		DurationSec: durationSec,
	}

	var repoPath string
	var err error

	// Check if repoURL is a local path or a git URL
	if _, err := os.Stat(repoURL); err == nil {
		// It's a local path
		repoPath = repoURL
		fmt.Printf("📁 Using local repository: %s\n", repoPath)
	} else if strings.HasPrefix(repoURL, "http://") || strings.HasPrefix(repoURL, "https://") || strings.HasPrefix(repoURL, "git@") {
		// It's a git URL - clone it
		repoName := filepath.Base(strings.TrimSuffix(repoURL, ".git"))
		repoPath = filepath.Join(outDir, "repo", repoName)
		
		if err := cloneRepo(ctx, repoURL, ref, repoPath); err != nil {
			manifest.Success = false
			manifest.Error = fmt.Sprintf("clone failed: %v", err)
			return manifest, fmt.Errorf("demo failed: %w", err)
		}
	} else {
		// Assume it's a local path that doesn't exist yet
		repoPath = repoURL
	}

	manifest.LocalPath = repoPath

	// Get Go version
	goVersion, err := getGoVersion(ctx)
	if err != nil {
		log.Printf("Warning: failed to get Go version: %v", err)
		goVersion = "unknown"
	}
	manifest.GoVersion = goVersion

	// Detect benchmarks
	benchmarks, err := detectBenchmarks(ctx, repoPath)
	if err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("benchmark detection failed: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	if len(benchmarks) == 0 {
		manifest.Success = false
		manifest.Error = "no benchmarks found"
		return manifest, fmt.Errorf("no benchmarks found in repository")
	}

	manifest.Benchmarks = benchmarks

	// Collect profiles
	bundle, err := p.collectProfiles(ctx, repoPath, durationSec, outDir)
	if err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("profile collection failed: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Create output directory and profiles subdirectory if they don't exist
	if err := os.MkdirAll(outDir, 0755); err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to create output directory: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(outDir, "profiles"), 0755); err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to create profiles directory: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Save bundle
	bundlePath := filepath.Join(outDir, "bundle.json")
	bundleData, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to serialize bundle: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	if err := os.WriteFile(bundlePath, bundleData, 0644); err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to write bundle: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Analyze profiles with deterministic rules
	findingsPath := filepath.Join(outDir, "findings.json")
	_, err = p.AnalyzeWithDeterministicRules(ctx, bundlePath, 20, findingsPath)
	if err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("analysis failed: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Generate report
	reportPath := filepath.Join(outDir, "report.md")
	if err := p.Report(ctx, findingsPath, reportPath); err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("report generation failed: %v", err)
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Profiles were already copied during generation, just update the manifest
	for _, artifact := range bundle.Artifacts {
		if artifact.Kind == "profile" {
			// Use the output directory path instead of the repo path
			destPath := filepath.Join(outDir, "profiles", filepath.Base(artifact.Path))
			manifest.Profiles = append(manifest.Profiles, destPath)
		}
	}

	manifest.Success = true
	return manifest, nil
}

// getGoVersion gets the current Go version
func getGoVersion(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Go version: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}
