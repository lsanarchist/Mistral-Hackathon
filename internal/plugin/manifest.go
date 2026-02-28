package plugin

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Manifest represents a plugin manifest file
type Manifest struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	SDKVersion  string      `json:"sdkVersion"`
	Capabilities Capabilities `json:"capabilities"`
	Description string      `json:"description,omitempty"`
	Author      string      `json:"author,omitempty"`
}

// Capabilities defines what a plugin can handle
type Capabilities struct {
	Targets  []string `json:"targets"`
	Profiles []string `json:"profiles"`
}

// SDKVersionCompatibility defines the core's supported SDK version
const SDKVersionCompatibility = "1.0"

// LoadManifest loads and parses a plugin manifest file with strict validation
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest %s: %w", path, err)
	}

	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()

	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest %s: %w", path, err)
	}

	// Basic validation
	if manifest.Name == "" {
		return nil, fmt.Errorf("manifest %s: name is required", path)
	}
	if manifest.Version == "" {
		return nil, fmt.Errorf("manifest %s: version is required", path)
	}
	if manifest.SDKVersion == "" {
		return nil, fmt.Errorf("manifest %s: sdkVersion is required", path)
	}

	return &manifest, nil
}

// DiscoverManifests finds all valid plugin manifests in a directory
func DiscoverManifests(manifestsDir string) ([]*Manifest, error) {
	var manifests []*Manifest

	err := filepath.WalkDir(manifestsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Only process .json files
		if filepath.Ext(path) != ".json" {
			return nil
		}

		manifest, err := LoadManifest(path)
		if err != nil {
			// Log warning but continue with other manifests
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid manifest %s: %v\n", path, err)
			return nil
		}

		manifests = append(manifests, manifest)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover manifests: %w", err)
	}

	// Sort by name for consistent ordering
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].Name < manifests[j].Name
	})

	return manifests, nil
}

// ResolvePlugin finds a plugin manifest and validates its binary exists
func ResolvePlugin(manifestsDir, binDir, name string) (*Manifest, string, error) {
	manifests, err := DiscoverManifests(manifestsDir)
	if err != nil {
		return nil, "", fmt.Errorf("failed to discover plugins: %w", err)
	}

	// Find the requested plugin
	var manifest *Manifest
	for _, m := range manifests {
		if m.Name == name {
			manifest = m
			break
		}
	}

	if manifest == nil {
		available := make([]string, 0, len(manifests))
		for _, m := range manifests {
			available = append(available, m.Name)
		}
		return nil, "", fmt.Errorf("plugin %q not found. Available plugins: %s", name, strings.Join(available, ", "))
	}

	// Check SDK compatibility
	if manifest.SDKVersion != SDKVersionCompatibility {
		return nil, "", fmt.Errorf("plugin %s requires sdkVersion %s, but core supports %s", 
			manifest.Name, manifest.SDKVersion, SDKVersionCompatibility)
	}

	// Check binary exists
	binaryPath := filepath.Join(binDir, manifest.Name)
	if _, err := os.Stat(binaryPath); err != nil {
		if os.IsNotExist(err) {
			return nil, "", fmt.Errorf("plugin %s manifest found but binary missing at %s", manifest.Name, binaryPath)
		}
		return nil, "", fmt.Errorf("failed to access plugin binary %s: %w", binaryPath, err)
	}

	return manifest, binaryPath, nil
}

// ValidateTarget checks if a target type is supported by the plugin
func (m *Manifest) ValidateTarget(targetType string) error {
	for _, t := range m.Capabilities.Targets {
		if t == targetType {
			return nil
		}
	}
	return fmt.Errorf("target type %q not supported by plugin %s. Supported targets: %s", 
		targetType, m.Name, strings.Join(m.Capabilities.Targets, ", "))
}

// ValidateProfiles checks if requested profiles are supported by the plugin
func (m *Manifest) ValidateProfiles(requested []string) error {
	supported := make(map[string]bool, len(m.Capabilities.Profiles))
	for _, p := range m.Capabilities.Profiles {
		supported[p] = true
	}

	var unsupported []string
	for _, profile := range requested {
		if !supported[profile] {
			unsupported = append(unsupported, profile)
		}
	}

	if len(unsupported) > 0 {
		return fmt.Errorf("profiles %s not supported by plugin %s. Supported profiles: %s", 
			strings.Join(unsupported, ", "), m.Name, strings.Join(m.Capabilities.Profiles, ", "))
	}

	return nil
}