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
		return nil, fmt.Errorf("failed to read manifest file %s: %w", path, err)
	}

	var manifest Manifest
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest file %s: %w", path, err)
	}

	// Validate required fields
	if manifest.Name == "" {
		return nil, fmt.Errorf("manifest %s is missing required field 'name'", path)
	}
	if manifest.Version == "" {
		return nil, fmt.Errorf("manifest %s is missing required field 'version'", path)
	}
	if manifest.SDKVersion == "" {
		return nil, fmt.Errorf("manifest %s is missing required field 'sdkVersion'", path)
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

		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		manifest, err := LoadManifest(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: skipping invalid manifest %s: %v\n", path, err)
			return nil
		}

		manifests = append(manifests, manifest)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk manifests directory %s: %w", manifestsDir, err)
	}

	// Sort by plugin name for consistent ordering
	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].Name < manifests[j].Name
	})

	return manifests, nil
}

// ResolvePlugin finds a plugin manifest and validates its binary exists
func ResolvePlugin(manifestsDir, binDir, name string) (*Manifest, string, error) {
	manifests, err := DiscoverManifests(manifestsDir)
	if err != nil {
		return nil, "", fmt.Errorf("failed to discover manifests: %w", err)
	}

	var foundManifest *Manifest
	for _, m := range manifests {
		if m.Name == name {
			foundManifest = m
			break
		}
	}

	if foundManifest == nil {
		available := make([]string, len(manifests))
		for i, m := range manifests {
			available[i] = m.Name
		}
		return nil, "", fmt.Errorf("plugin %s not found. Available plugins: %s", name, strings.Join(available, ", "))
	}

	// Check SDK version compatibility
	if foundManifest.SDKVersion != SDKVersionCompatibility {
		return nil, "", fmt.Errorf("plugin %s requires sdkVersion %s, but core supports %s", 
			foundManifest.Name, foundManifest.SDKVersion, SDKVersionCompatibility)
	}

	// Validate binary exists
	binaryPath := filepath.Join(binDir, foundManifest.Name)
	if _, err := os.Stat(binaryPath); err != nil {
		if os.IsNotExist(err) {
			return nil, "", fmt.Errorf("plugin %s manifest found but binary missing at path %s", foundManifest.Name, binaryPath)
		}
		return nil, "", fmt.Errorf("failed to check plugin %s binary: %w", foundManifest.Name, err)
	}

	return foundManifest, binaryPath, nil
}

// ValidateTarget checks if a target type is supported by the plugin
func (m *Manifest) ValidateTarget(targetType string) error {
	for _, supported := range m.Capabilities.Targets {
		if supported == targetType {
			return nil
		}
	}
	return fmt.Errorf("target type '%s' not supported by plugin %s. Supported targets: %s", 
		targetType, m.Name, strings.Join(m.Capabilities.Targets, ", "))
}

// ValidateProfiles checks if requested profiles are supported by the plugin
func (m *Manifest) ValidateProfiles(requested []string) error {
	var unsupported []string
	for _, req := range requested {
		supported := false
		for _, supp := range m.Capabilities.Profiles {
			if supp == req {
				supported = true
				break
			}
		}
		if !supported {
			unsupported = append(unsupported, req)
		}
	}

	if len(unsupported) > 0 {
		return fmt.Errorf("profiles '%s' not supported by plugin %s. Supported profiles: %s", 
			strings.Join(unsupported, ", "), m.Name, strings.Join(m.Capabilities.Profiles, ", "))
	}
	return nil
}
