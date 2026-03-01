package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
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
	ErrorContext *model.ErrorContext `json:"errorContext,omitempty"`
	PerformanceConfig *model.PerformanceOptimizationConfig `json:"performanceConfig,omitempty"`
	RemediationConfig *model.RemediationConfig `json:"remediationConfig,omitempty"`
}

// cloneRepo clones a Git repository to the specified directory
func cloneRepo(ctx context.Context, repoURL, ref, destDir string) error {
	// Remove existing directory if it exists
	if _, err := os.Stat(destDir); err == nil {
		if err := os.RemoveAll(destDir); err != nil {
			errContext := model.NewErrorContext(
				model.ErrorTypeIO,
				model.ErrorCodeFileOperation,
				"Failed to remove existing directory",
				fmt.Sprintf("Directory: %s, Error: %v", destDir, err),
				"Try running with elevated permissions or check if files are locked",
				true,
			)
			return errContext
		}
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destDir), 0755); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeIO,
			model.ErrorCodeFileOperation,
			"Failed to create parent directory",
			fmt.Sprintf("Directory: %s, Error: %v", filepath.Dir(destDir), err),
			"Check directory permissions and disk space",
			true,
		)
		return errContext
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeDependency,
			model.ErrorCodeDependencyMissing,
			"Git is not installed or not in PATH",
			"Git executable not found",
			"Install Git and ensure it's in your system PATH",
			true,
		)
		return errContext
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
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeGitCloneFailed,
			"Git clone operation failed",
			fmt.Sprintf("Repository: %s, Reference: %s\nstdout: %s\nstderr: %s", repoURL, ref, stdout.String(), stderr.String()),
			"Check repository URL, network connectivity, and authentication credentials",
			true,
		)
		return errContext
	}

	return nil
}

// detectBenchmarks finds Go benchmark functions in the repository
func detectBenchmarks(ctx context.Context, repoPath string) ([]string, error) {
	var benchmarks []string

	// Check if repository path exists
	if _, err := os.Stat(repoPath); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeValidation,
			model.ErrorCodeFileOperation,
			"Repository path does not exist",
			fmt.Sprintf("Path: %s, Error: %v", repoPath, err),
			"Verify the repository path is correct and accessible",
			true,
		)
		return nil, errContext
	}

	// Find all Go files in the repository
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				errContext := model.NewErrorContext(
					model.ErrorTypeIO,
					model.ErrorCodeFileOperation,
					"Permission denied accessing file",
					fmt.Sprintf("Path: %s, Error: %v", path, err),
					"Check file permissions and try running with elevated privileges",
					true,
				)
				return errContext
			}
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			errContext := model.NewErrorContext(
				model.ErrorTypeIO,
				model.ErrorCodeFileOperation,
				"Failed to read Go file",
				fmt.Sprintf("File: %s, Error: %v", path, err),
				"Check file permissions and disk health",
				true,
			)
			return errContext
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
		var errContext model.ErrorContext
		if errors.As(err, &errContext) {
			return nil, errContext
		}
		errContext = model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeFileOperation,
			"Failed to walk repository",
			fmt.Sprintf("Path: %s, Error: %v", repoPath, err),
			"Check directory structure and file permissions",
			true,
		)
		return nil, errContext
	}

	if len(benchmarks) == 0 {
		errContext := model.NewErrorContext(
			model.ErrorTypeValidation,
			model.ErrorCodeNoBenchmarksFound,
			"No Go benchmarks found in repository",
			fmt.Sprintf("Searched path: %s", repoPath),
			"Ensure your Go files contain functions starting with 'func Benchmark' and follow Go benchmark naming conventions",
			true,
		)
		return nil, errContext
	}

	return benchmarks, nil
}

// runBenchmarks runs Go benchmarks and collects profiles
func runBenchmarks(ctx context.Context, repoPath string, durationSec int, outDir string, perfConfig *model.PerformanceOptimizationConfig) ([]string, error) {
	var profilePaths []string

	// Check if Go is available
	if _, err := exec.LookPath("go"); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeDependency,
			model.ErrorCodeDependencyMissing,
			"Go is not installed or not in PATH",
			"Go executable not found",
			"Install Go and ensure it's in your system PATH",
			true,
		)
		return nil, errContext
	}

	// Change to repository directory
	originalDir, err := os.Getwd()
	if err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeFileOperation,
			"Failed to get current directory",
			fmt.Sprintf("Error: %v", err),
			"Check working directory permissions",
			true,
		)
		return nil, errContext
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(repoPath); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeFileOperation,
			"Failed to change to repository directory",
			fmt.Sprintf("Directory: %s, Error: %v", repoPath, err),
			"Check directory exists and permissions are correct",
			true,
		)
		return nil, errContext
	}

	// Create profiles directory in repo
	profilesDir := filepath.Join(repoPath, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeIO,
			model.ErrorCodeFileOperation,
			"Failed to create profiles directory",
			fmt.Sprintf("Directory: %s, Error: %v", profilesDir, err),
			"Check directory permissions and disk space",
			true,
		)
		return nil, errContext
	}
	
	// Create profiles directory in output
	outputProfilesDir := filepath.Join(outDir, "profiles")
	if err := os.MkdirAll(outputProfilesDir, 0755); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeIO,
			model.ErrorCodeFileOperation,
			"Failed to create output profiles directory",
			fmt.Sprintf("Directory: %s, Error: %v", outputProfilesDir, err),
			"Check output directory permissions and disk space",
			true,
		)
		return nil, errContext
	}

	// Run benchmarks with profiling - use shorter duration for each profile type
	// to avoid timeout issues, and run a subset of benchmarks
	profileTypes := []string{"cpu", "heap", "allocs"} // Reduced set for demo
	benchmarkDuration := durationSec / len(profileTypes) // Split duration among profile types
	if benchmarkDuration < 1 {
		benchmarkDuration = 1 // Minimum 1 second per profile
	}

	// Determine if we should run concurrently
	shouldRunConcurrently := perfConfig != nil && perfConfig.EnableConcurrentBenchmarks && perfConfig.MaxConcurrentWorkers > 1
	maxWorkers := 1
	if shouldRunConcurrently {
		maxWorkers = perfConfig.MaxConcurrentWorkers
		if maxWorkers > len(profileTypes) {
			maxWorkers = len(profileTypes) // Don't exceed number of profile types
		}
	}

	if shouldRunConcurrently {
		// Run benchmarks concurrently
		profilePaths, err = runBenchmarksConcurrently(ctx, repoPath, profileTypes, benchmarkDuration, profilesDir, outputProfilesDir, maxWorkers)
		if err != nil {
			return nil, err
		}
	} else {
		// Run benchmarks sequentially (original behavior)
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
					errContext := model.NewErrorContext(
						model.ErrorTypeExecution,
						model.ErrorCodeBenchmarkExecution,
						"Benchmark execution failed and no profile generated",
						fmt.Sprintf("Profile type: %s\nError: %v\nstdout: %s\nstderr: %s", profileType, err, stdout.String(), stderr.String()),
						"Check benchmark code, dependencies, and test environment",
						true,
					)
					return nil, errContext
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
	}

	if len(profilePaths) == 0 {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeProfileCollection,
			"No profiles were generated",
			"Benchmark execution completed but no profile files were created",
			"Check benchmark configuration and ensure profiling flags are correctly specified",
			true,
		)
		return nil, errContext
	}

	return profilePaths, nil
}

// runBenchmarksConcurrently runs benchmarks in parallel using worker pool pattern
func runBenchmarksConcurrently(ctx context.Context, repoPath string, profileTypes []string, benchmarkDuration int, profilesDir, outputProfilesDir string, maxWorkers int) ([]string, error) {
	var profilePaths []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstError error
	var errorMutex sync.Mutex

	// Create worker pool
	workerChan := make(chan string, maxWorkers)
	resultChan := make(chan string, len(profileTypes))
	errorChan := make(chan error, len(profileTypes))

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			for profileType := range workerChan {
				profilePath := filepath.Join(profilesDir, fmt.Sprintf("%s.pprof", profileType))

				// Build benchmark command
				args := []string{"test", "-bench=BenchmarkProcessStrings|BenchmarkGenerateRandomData|BenchmarkProcessJSON", "-benchtime", fmt.Sprintf("%ds", benchmarkDuration)}
				if profileType == "cpu" {
					args = append(args, "-cpuprofile", profilePath)
				} else if profileType == "heap" || profileType == "allocs" {
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
					// Only report error if the profile file doesn't exist
					if _, statErr := os.Stat(profilePath); statErr != nil {
						log.Printf("Worker %d: Benchmark command failed for %s: %v\nArgs: %v\nStdout: %s\nStderr: %s", workerID, profileType, err, args, stdout.String(), stderr.String())
						errorChan <- model.NewErrorContext(
							model.ErrorTypeExecution,
							model.ErrorCodeBenchmarkExecution,
							"Benchmark execution failed and no profile generated",
							fmt.Sprintf("Profile type: %s\nError: %v\nstdout: %s\nstderr: %s", profileType, err, stdout.String(), stderr.String()),
							"Check benchmark code, dependencies, and test environment",
							true,
						)
						continue
					}
				}

				// Check if profile file was created and copy it immediately
				if _, err := os.Stat(profilePath); err == nil {
					// Copy the profile to a safe location immediately
					destProfilePath := filepath.Join(outputProfilesDir, filepath.Base(profilePath))
					if err := copyFile(profilePath, destProfilePath); err != nil {
						log.Printf("Worker %d: Warning: Failed to copy profile %s to %s: %v", workerID, profilePath, destProfilePath, err)
					} else {
						resultChan <- destProfilePath
					}
				}
			}
		}(i)
	}

	// Distribute work to workers
	for _, profileType := range profileTypes {
		workerChan <- profileType
	}
	close(workerChan)

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	// Collect results
	for result := range resultChan {
		mu.Lock()
		profilePaths = append(profilePaths, result)
		mu.Unlock()
	}

	// Check for errors
	for err := range errorChan {
		errorMutex.Lock()
		if firstError == nil {
			firstError = err
		}
		errorMutex.Unlock()
	}

	if firstError != nil {
		return nil, firstError
	}

	if len(profilePaths) == 0 {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeProfileCollection,
			"No profiles were generated in concurrent execution",
			"Benchmark execution completed but no profile files were created",
			"Check benchmark configuration and ensure profiling flags are correctly specified",
			true,
		)
		return nil, errContext
	}

	return profilePaths, nil
}

// collectProfiles collects profiles from the repository
func (p *Pipeline) collectProfiles(ctx context.Context, repoPath string, durationSec int, outDir string, perfConfig *model.PerformanceOptimizationConfig) (*model.ProfileBundle, error) {
	// Detect benchmarks
	benchmarks, err := detectBenchmarks(ctx, repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect benchmarks: %w", err)
	}

	if len(benchmarks) == 0 {
		return nil, fmt.Errorf("no benchmarks found in repository")
	}

	// Run benchmarks and collect profiles
	profilePaths, err := runBenchmarks(ctx, repoPath, durationSec, outDir, perfConfig)
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
	// Create default performance config (backward compatibility)
	defaultPerfConfig := &model.PerformanceOptimizationConfig{
		EnableConcurrentBenchmarks: false,
		MaxConcurrentWorkers:       1,
		EnableProfileSampling:      false,
		SamplingRate:               1.0,
		EnableMemoryOptimization:   false,
		LargeCodebaseMode:          false,
	}
	return p.DemoWithPerformance(ctx, repoURL, ref, outDir, durationSec, defaultPerfConfig)
}

// DemoWithPerformance runs the complete demo workflow with performance optimization
func (p *Pipeline) DemoWithPerformance(ctx context.Context, repoURL, ref, outDir string, durationSec int, perfConfig *model.PerformanceOptimizationConfig) (*RunManifest, error) {
	manifest := &RunManifest{
		RepoURL:     repoURL,
		RepoRef:     ref,
		Timestamp:   time.Now(),
		DurationSec: durationSec,
		PerformanceConfig: perfConfig,
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
			
			// Extract error context if available
			var errContext model.ErrorContext
			if errors.As(err, &errContext) {
				manifest.ErrorContext = &errContext
			}
			
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
		
		// Extract error context if available
		var errContext model.ErrorContext
		if errors.As(err, &errContext) {
			manifest.ErrorContext = &errContext
		}
		
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	if len(benchmarks) == 0 {
		errContext := model.NewErrorContext(
			model.ErrorTypeValidation,
			model.ErrorCodeNoBenchmarksFound,
			"No Go benchmarks found in repository",
			fmt.Sprintf("Searched path: %s", repoPath),
			"Ensure your Go files contain functions starting with 'func Benchmark' and follow Go benchmark naming conventions",
			true,
		)
		manifest.Success = false
		manifest.Error = "no benchmarks found"
		manifest.ErrorContext = &errContext
		return manifest, errContext
	}

	manifest.Benchmarks = benchmarks

	// Collect profiles
	bundle, err := p.collectProfiles(ctx, repoPath, durationSec, outDir, perfConfig)
	if err != nil {
		manifest.Success = false
		manifest.Error = fmt.Sprintf("profile collection failed: %v", err)
		
		// Extract error context if available
		var errContext model.ErrorContext
		if errors.As(err, &errContext) {
			manifest.ErrorContext = &errContext
		}
		
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Create output directory and profiles subdirectory if they don't exist
	if err := os.MkdirAll(outDir, 0755); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeIO,
			model.ErrorCodeFileOperation,
			"Failed to create output directory",
			fmt.Sprintf("Directory: %s, Error: %v", outDir, err),
			"Check directory permissions and disk space",
			true,
		)
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to create output directory: %v", err)
		manifest.ErrorContext = &errContext
		return manifest, fmt.Errorf("demo failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(outDir, "profiles"), 0755); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeIO,
			model.ErrorCodeFileOperation,
			"Failed to create profiles directory",
			fmt.Sprintf("Directory: %s, Error: %v", filepath.Join(outDir, "profiles"), err),
			"Check directory permissions and disk space",
			true,
		)
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to create profiles directory: %v", err)
		manifest.ErrorContext = &errContext
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Save bundle
	bundlePath := filepath.Join(outDir, "bundle.json")
	bundleData, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeFileOperation,
			"Failed to serialize bundle",
			fmt.Sprintf("Error: %v", err),
			"Check bundle data structure and available memory",
			true,
		)
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to serialize bundle: %v", err)
		manifest.ErrorContext = &errContext
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	if err := os.WriteFile(bundlePath, bundleData, 0644); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeIO,
			model.ErrorCodeFileOperation,
			"Failed to write bundle file",
			fmt.Sprintf("File: %s, Error: %v", bundlePath, err),
			"Check file permissions and disk space",
			true,
		)
		manifest.Success = false
		manifest.Error = fmt.Sprintf("failed to write bundle: %v", err)
		manifest.ErrorContext = &errContext
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Analyze profiles with deterministic rules and performance configuration
	findingsPath := filepath.Join(outDir, "findings.json")
	_, err = p.AnalyzeWithDeterministicRulesAndOptions(ctx, bundlePath, 20, findingsPath, perfConfig)
	if err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeFileOperation,
			"Analysis failed",
			fmt.Sprintf("Error: %v", err),
			"Check profile data and analysis configuration",
			true,
		)
		manifest.Success = false
		manifest.Error = fmt.Sprintf("analysis failed: %v", err)
		manifest.ErrorContext = &errContext
		return manifest, fmt.Errorf("demo failed: %w", err)
	}

	// Generate report
	reportPath := filepath.Join(outDir, "report.md")
	if err := p.Report(ctx, findingsPath, reportPath); err != nil {
		errContext := model.NewErrorContext(
			model.ErrorTypeExecution,
			model.ErrorCodeFileOperation,
			"Report generation failed",
			fmt.Sprintf("Error: %v", err),
			"Check findings data and report template",
			true,
		)
		manifest.Success = false
		manifest.Error = fmt.Sprintf("report generation failed: %v", err)
		manifest.ErrorContext = &errContext
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
