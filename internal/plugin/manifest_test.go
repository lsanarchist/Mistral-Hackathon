package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadManifest_StrictParsing(t *testing.T) {
	t.Run("valid manifest", func(t *testing.T) {
		manifest := `{
			"name": "test-plugin",
			"version": "1.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu", "heap"]
			},
			"description": "Test plugin",
			"author": "Test Author"
		}`

		tmpDir := t.TempDir()
		manifestPath := filepath.Join(tmpDir, "test.json")
		if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
			t.Fatalf("Failed to create test manifest: %v", err)
		}

		m, err := LoadManifest(manifestPath)
		if err != nil {
			t.Fatalf("LoadManifest failed: %v", err)
		}

		if m.Name != "test-plugin" {
			t.Errorf("Expected name 'test-plugin', got %q", m.Name)
		}
		if m.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got %q", m.Version)
		}
		if m.SDKVersion != "1.0" {
			t.Errorf("Expected sdkVersion '1.0', got %q", m.SDKVersion)
		}
		if len(m.Capabilities.Targets) != 1 || m.Capabilities.Targets[0] != "url" {
			t.Errorf("Expected targets ['url'], got %v", m.Capabilities.Targets)
		}
		if len(m.Capabilities.Profiles) != 2 {
			t.Errorf("Expected 2 profiles, got %d", len(m.Capabilities.Profiles))
		}
	})

	t.Run("unknown field fails", func(t *testing.T) {
		manifest := `{
			"name": "test-plugin",
			"version": "1.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu"]
			},
			"unknownField": "should fail"
		}`

		tmpDir := t.TempDir()
		manifestPath := filepath.Join(tmpDir, "test.json")
		if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
			t.Fatalf("Failed to create test manifest: %v", err)
		}

		_, err := LoadManifest(manifestPath)
		if err == nil {
			t.Fatal("Expected error for unknown field, got nil")
		}
		if !strings.Contains(err.Error(), "unknown field") {
			t.Errorf("Expected error about unknown field, got: %v", err)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		manifest := `{
			"version": "1.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu"]
			}
		}` // Missing "name"

		tmpDir := t.TempDir()
		manifestPath := filepath.Join(tmpDir, "test.json")
		if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
			t.Fatalf("Failed to create test manifest: %v", err)
		}

		_, err := LoadManifest(manifestPath)
		if err == nil {
			t.Fatal("Expected error for missing name, got nil")
		}
		if !strings.Contains(err.Error(), "name is required") {
			t.Errorf("Expected error about missing name, got: %v", err)
		}
	})
}

func TestDiscoverManifests(t *testing.T) {
	t.Run("discovers all json files", func(t *testing.T) {
		tmpDir := t.TempDir()
		manifestsDir := filepath.Join(tmpDir, "manifests")
		if err := os.MkdirAll(manifestsDir, 0755); err != nil {
			t.Fatalf("Failed to create manifests dir: %v", err)
		}

		// Create valid manifest
		validManifest := `{
			"name": "plugin-a",
			"version": "1.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu"]
			}
		}`
		if err := os.WriteFile(filepath.Join(manifestsDir, "plugin-a.json"), []byte(validManifest), 0644); err != nil {
			t.Fatalf("Failed to create valid manifest: %v", err)
		}

		// Create another valid manifest
		validManifest2 := `{
			"name": "plugin-b",
			"version": "2.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["heap"]
			}
		}`
		if err := os.WriteFile(filepath.Join(manifestsDir, "plugin-b.json"), []byte(validManifest2), 0644); err != nil {
			t.Fatalf("Failed to create valid manifest: %v", err)
		}

		// Create invalid manifest (should be skipped with warning)
		invalidManifest := `{
			"name": "plugin-invalid",
			"version": "1.0.0",
			// Missing sdkVersion
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu"]
			}
		}`
		if err := os.WriteFile(filepath.Join(manifestsDir, "plugin-invalid.json"), []byte(invalidManifest), 0644); err != nil {
			t.Fatalf("Failed to create invalid manifest: %v", err)
		}

		// Create non-json file (should be ignored)
		if err := os.WriteFile(filepath.Join(manifestsDir, "readme.txt"), []byte("not a manifest"), 0644); err != nil {
			t.Fatalf("Failed to create non-json file: %v", err)
		}

		manifests, err := DiscoverManifests(manifestsDir)
		if err != nil {
			t.Fatalf("DiscoverManifests failed: %v", err)
		}

		// Should find 2 valid manifests
		if len(manifests) != 2 {
			t.Errorf("Expected 2 manifests, got %d", len(manifests))
		}

		// Check sorting
		if manifests[0].Name != "plugin-a" {
			t.Errorf("Expected first manifest to be 'plugin-a', got %q", manifests[0].Name)
		}
		if manifests[1].Name != "plugin-b" {
			t.Errorf("Expected second manifest to be 'plugin-b', got %q", manifests[1].Name)
		}
	})
}

func TestResolvePlugin(t *testing.T) {
	t.Run("fails when binary missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		manifestsDir := filepath.Join(tmpDir, "manifests")
		binDir := filepath.Join(tmpDir, "bin")
		
		if err := os.MkdirAll(manifestsDir, 0755); err != nil {
			t.Fatalf("Failed to create manifests dir: %v", err)
		}
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatalf("Failed to create bin dir: %v", err)
		}

		// Create manifest but no binary
		manifest := `{
			"name": "missing-binary",
			"version": "1.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu"]
			}
		}`
		if err := os.WriteFile(filepath.Join(manifestsDir, "missing-binary.json"), []byte(manifest), 0644); err != nil {
			t.Fatalf("Failed to create manifest: %v", err)
		}

		_, _, err := ResolvePlugin(manifestsDir, binDir, "missing-binary")
		if err == nil {
			t.Fatal("Expected error for missing binary, got nil")
		}
		if !strings.Contains(err.Error(), "binary missing") {
			t.Errorf("Expected error about missing binary, got: %v", err)
		}
	})

	t.Run("succeeds when manifest and binary present", func(t *testing.T) {
		tmpDir := t.TempDir()
		manifestsDir := filepath.Join(tmpDir, "manifests")
		binDir := filepath.Join(tmpDir, "bin")
		
		if err := os.MkdirAll(manifestsDir, 0755); err != nil {
			t.Fatalf("Failed to create manifests dir: %v", err)
		}
		if err := os.MkdirAll(binDir, 0755); err != nil {
			t.Fatalf("Failed to create bin dir: %v", err)
		}

		// Create manifest
		manifest := `{
			"name": "valid-plugin",
			"version": "1.0.0",
			"sdkVersion": "1.0",
			"capabilities": {
				"targets": ["url"],
				"profiles": ["cpu", "heap"]
			},
			"description": "Valid plugin for testing"
		}`
		if err := os.WriteFile(filepath.Join(manifestsDir, "valid-plugin.json"), []byte(manifest), 0644); err != nil {
			t.Fatalf("Failed to create manifest: %v", err)
		}

		// Create dummy executable
		dummyBinary := filepath.Join(binDir, "valid-plugin")
		if err := os.WriteFile(dummyBinary, []byte("#!/bin/sh\necho 'dummy plugin'"), 0755); err != nil {
			t.Fatalf("Failed to create dummy binary: %v", err)
		}

		m, binaryPath, err := ResolvePlugin(manifestsDir, binDir, "valid-plugin")
		if err != nil {
			t.Fatalf("ResolvePlugin failed: %v", err)
		}

		if m.Name != "valid-plugin" {
			t.Errorf("Expected manifest name 'valid-plugin', got %q", m.Name)
		}
		if binaryPath != dummyBinary {
			t.Errorf("Expected binary path %q, got %q", dummyBinary, binaryPath)
		}
	})

	t.Run("fails for non-existent plugin", func(t *testing.T) {
		tmpDir := t.TempDir()
		manifestsDir := filepath.Join(tmpDir, "manifests")
		binDir := filepath.Join(tmpDir, "bin")
		
		if err := os.MkdirAll(manifestsDir, 0755); err != nil {
			t.Fatalf("Failed to create manifests dir: %v", err)
		}

		_, _, err := ResolvePlugin(manifestsDir, binDir, "non-existent")
		if err == nil {
			t.Fatal("Expected error for non-existent plugin, got nil")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected error about plugin not found, got: %v", err)
		}
	})
}

func TestManifestValidation(t *testing.T) {
	t.Run("target validation", func(t *testing.T) {
		manifest := &Manifest{
			Name:       "test-plugin",
			Version:    "1.0.0",
			SDKVersion: "1.0",
			Capabilities: Capabilities{
				Targets:  []string{"url", "file"},
				Profiles: []string{"cpu"},
			},
		}

		// Valid target
		if err := manifest.ValidateTarget("url"); err != nil {
			t.Errorf("Expected no error for valid target, got: %v", err)
		}

		// Invalid target
		err := manifest.ValidateTarget("database")
		if err == nil {
			t.Fatal("Expected error for invalid target, got nil")
		}
		if !strings.Contains(err.Error(), "not supported") {
			t.Errorf("Expected error about unsupported target, got: %v", err)
		}
	})

	t.Run("profile validation", func(t *testing.T) {
		manifest := &Manifest{
			Name:       "test-plugin",
			Version:    "1.0.0",
			SDKVersion: "1.0",
			Capabilities: Capabilities{
				Targets:  []string{"url"},
				Profiles: []string{"cpu", "heap"},
			},
		}

		// Valid profiles
		if err := manifest.ValidateProfiles([]string{"cpu"}); err != nil {
			t.Errorf("Expected no error for valid profiles, got: %v", err)
		}

		// Invalid profile
		err := manifest.ValidateProfiles([]string{"cpu", "mutex"})
		if err == nil {
			t.Fatal("Expected error for invalid profile, got nil")
		}
		if !strings.Contains(err.Error(), "not supported") {
			t.Errorf("Expected error about unsupported profile, got: %v", err)
		}
	})
}